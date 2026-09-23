package notifyjev

import (
	"context"
	"encoding/json"
	"log/slog"
	"sort"
)

// ShadowInput compares a new RoutingPlan against a host's legacy route (channels/targets only).
type ShadowInput struct {
	Context RoutingContext
	Legacy  []Destination
}

// ShadowObservation is the outcome of a shadow Resolve: plan, diff vs legacy, and metadata for metrics.
type ShadowObservation struct {
	Plan RoutingPlan
	Diff PlanDiff
}

// PlanDiff summarizes how the new plan differs from legacy delivery rules.
type PlanDiff struct {
	Match              bool
	Added              []Destination
	Removed            []Destination
	RequiredMismatches []DestinationPair
}

// DestinationPair names a channel/target pair where Required differed.
type DestinationPair struct {
	Channel  Channel
	Target   TargetKind
	Legacy   bool
	Resolved bool
}

// ShadowResolve runs Resolve and compares the result to legacy destinations without changing delivery.
func (r *Router) ShadowResolve(ctx context.Context, in ShadowInput) (ShadowObservation, error) {
	plan, err := r.Resolve(ctx, ResolveInput{Context: in.Context})
	if err != nil {
		return ShadowObservation{}, err
	}
	diff := DiffPlans(in.Legacy, plan.Destinations)
	return ShadowObservation{Plan: plan, Diff: diff}, nil
}

// DiffPlans compares legacy and resolved destination sets (channel + target + required).
func DiffPlans(legacy, resolved []Destination) PlanDiff {
	leg := indexDestinations(legacy)
	res := indexDestinations(resolved)

	var added, removed []Destination
	var reqMismatch []DestinationPair

	for k, d := range res {
		old, ok := leg[k]
		if !ok {
			added = append(added, d)
			continue
		}
		if old.Required != d.Required {
			reqMismatch = append(reqMismatch, DestinationPair{
				Channel: k.channel, Target: k.target, Legacy: old.Required, Resolved: d.Required,
			})
		}
	}
	for k, d := range leg {
		if _, ok := res[k]; !ok {
			removed = append(removed, d)
		}
	}

	sortDestinations(added)
	sortDestinations(removed)
	sort.Slice(reqMismatch, func(i, j int) bool {
		if reqMismatch[i].Channel != reqMismatch[j].Channel {
			return reqMismatch[i].Channel < reqMismatch[j].Channel
		}
		return reqMismatch[i].Target < reqMismatch[j].Target
	})

	match := len(added) == 0 && len(removed) == 0 && len(reqMismatch) == 0
	return PlanDiff{
		Match:              match,
		Added:              added,
		Removed:            removed,
		RequiredMismatches: reqMismatch,
	}
}

// LogShadowObservation emits structured logs for shadow mode (plan + metadata only, no PII state).
func LogShadowObservation(logger *slog.Logger, category Category, obs ShadowObservation) {
	if logger == nil {
		logger = slog.Default()
	}
	attrs := []any{
		"plan_id", obs.Plan.PlanID,
		"category", category,
		"destinations", PlanDestinationsJSON(obs.Plan),
		"diff_match", obs.Diff.Match,
		"diff_added", len(obs.Diff.Added),
		"diff_removed", len(obs.Diff.Removed),
		"diff_required_mismatch", len(obs.Diff.RequiredMismatches),
		"question_set", obs.Plan.Metadata.QuestionSetVersion,
		"jev_model", obs.Plan.Metadata.JevModel,
		"fallback_mode", obs.Plan.Metadata.FallbackMode,
		"fallback_questions", obs.Plan.Metadata.FallbackQuestions,
		"policy_overrides", obs.Plan.Metadata.PolicyOverrides,
	}
	logger.Info("notifyjev shadow observation", attrs...)
}

// ShadowMetrics aggregates shadow observations (e.g. per category).
type ShadowMetrics struct {
	Total            int
	Matches          int
	Mismatches       int
	ByCategory       map[Category]ShadowCategoryStats
	RequiredMismatch int
}

// ShadowCategoryStats holds per-category shadow counters.
type ShadowCategoryStats struct {
	Total      int
	Matches    int
	Mismatches int
}

// RecordShadowObservation updates rolling metrics from one observation.
func (m *ShadowMetrics) RecordShadowObservation(category Category, obs ShadowObservation) {
	if m.ByCategory == nil {
		m.ByCategory = make(map[Category]ShadowCategoryStats)
	}
	m.Total++
	st := m.ByCategory[category]
	st.Total++
	if obs.Diff.Match {
		m.Matches++
		st.Matches++
	} else {
		m.Mismatches++
		st.Mismatches++
	}
	m.ByCategory[category] = st
	m.RequiredMismatch += len(obs.Diff.RequiredMismatches)
}

// PlanDestinationsJSON returns a compact JSON snapshot of destinations for audit logs.
func PlanDestinationsJSON(plan RoutingPlan) string {
	b, err := json.Marshal(plan.Destinations)
	if err != nil {
		return "[]"
	}
	return string(b)
}

func indexDestinations(list []Destination) map[destKey]Destination {
	out := make(map[destKey]Destination, len(list))
	for _, d := range list {
		out[destKeyOf(d)] = d
	}
	return out
}

func sortDestinations(list []Destination) {
	sort.Slice(list, func(i, j int) bool {
		if list[i].Channel != list[j].Channel {
			return list[i].Channel < list[j].Channel
		}
		return list[i].Target < list[j].Target
	})
}
