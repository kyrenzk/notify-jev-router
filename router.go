package notifyjev

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	questionv0 "github.com/kyrenzk/notify-jev-router/questions"
)

// ResolveInput is the host-provided routing request.
type ResolveInput struct {
	Context RoutingContext
}

// Config configures Router dependencies for v0.1.0.
type Config struct {
	JevClient     JevClient
	Model         string
	Policies      PolicySet
	Confidence    ConfidenceConfig
	FallbackTable FallbackTable
}

// Router resolves notification routing plans via policy, Jev, and fallbacks.
type Router struct {
	cfg Config
}

// New validates required configuration and returns a Router.
func New(cfg Config) (*Router, error) {
	if cfg.JevClient == nil {
		return nil, errors.New("notifyjev: JevClient is required")
	}
	if strings.TrimSpace(cfg.Model) == "" {
		return nil, errors.New("notifyjev: Model is required")
	}
	if len(cfg.Policies.Rules) == 0 {
		cfg.Policies = DefaultPolicies
	}
	if cfg.FallbackTable == nil {
		cfg.FallbackTable = DefaultFallbackTable
	}
	if cfg.Confidence.ChoiceMinConfidence == 0 && cfg.Confidence.ChoiceMinMargin == 0 && len(cfg.Confidence.NoulThresholds) == 0 {
		cfg.Confidence = DefaultConfidenceConfig()
	}
	return &Router{cfg: cfg}, nil
}

// Resolve runs the full routing pipeline and always returns len(Destinations) >= 1 on success.
func (r *Router) Resolve(ctx context.Context, in ResolveInput) (RoutingPlan, error) {
	rc := in.Context
	if err := rc.Validate(); err != nil {
		return RoutingPlan{}, err
	}

	policyDests, policyOverrides := r.cfg.Policies.Apply(rc)
	jevQuestions := questionv0.Build(string(rc.Category))
	asked := questionIDs(jevQuestions)

	meta := PlanMetadata{
		QuestionSetVersion: questionv0.Version,
		JevModel:           r.cfg.Model,
		PolicyOverrides:    policyOverrides,
	}

	var interpreted []Destination
	var jevModel string

	if rc.Capabilities.AnyDeliverable(rc.UserPrefs) {
		resp, err := r.cfg.JevClient.SystemOne(ctx, rc.JevState(), jevQuestions)
		if err != nil {
			interpreted = FullFallback(rc, r.cfg.FallbackTable)
			meta.FallbackMode = "full"
		} else {
			if resp.Model != "" {
				jevModel = resp.Model
			}
			interp := InterpretAnswers(rc, resp, asked, r.cfg.Confidence)
			interpreted = interp.Destinations
			meta.AnswersSnapshot = interp.AnswersSnapshot
			meta.FallbackQuestions = interp.FallbackQuestions
			if len(interp.FallbackQuestions) > 0 {
				meta.FallbackMode = "partial"
			}
		}
	}

	merged := MergeDestinations(policyDests, interpreted)
	merged = EnsureNonEmpty(merged, r.cfg.FallbackTable.Lookup(rc.Category, rc.Severity))
	if len(merged) < 1 {
		return RoutingPlan{}, fmt.Errorf("%w: empty routing plan", ErrInvalidContext)
	}

	if jevModel != "" {
		meta.JevModel = jevModel
	}

	return RoutingPlan{
		PlanID:       buildPlanID(rc, questionv0.Version),
		Destinations: merged,
		Metadata:     meta,
	}, nil
}

func buildPlanID(rc RoutingContext, qsVersion string) string {
	sum := sha256.Sum256([]byte(rc.EventID + "|" + qsVersion + "|" + rc.TemplateID))
	return fmt.Sprintf("%s-%s", rc.EventID, hex.EncodeToString(sum[:8]))
}
