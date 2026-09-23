package notifyjev

import (
	"sort"
)

// Destination is one routed delivery target.
type Destination struct {
	Channel  Channel    `json:"channel"`
	Target   TargetKind `json:"target"`
	Required bool       `json:"required"`
	Priority int        `json:"priority"`
}

// AnswerSummary holds non-PII probability metadata for one Jev answer.
type AnswerSummary struct {
	Type          string             `json:"type"`
	Choice        string             `json:"choice,omitempty"`
	Noul          float64            `json:"noul,omitempty"`
	Confidence    float64            `json:"confidence,omitempty"`
	Probabilities map[string]float64 `json:"probabilities,omitempty"`
}

// PlanMetadata accompanies a RoutingPlan for audit and debugging.
type PlanMetadata struct {
	QuestionSetVersion string                    `json:"question_set_version"`
	JevModel           string                    `json:"jev_model"`
	AnswersSnapshot    map[string]AnswerSummary  `json:"answers_snapshot,omitempty"`
	FallbackQuestions  []string                  `json:"fallback_questions,omitempty"`
	PolicyOverrides    []string                  `json:"policy_overrides,omitempty"`
	FallbackMode       string                    `json:"fallback_mode,omitempty"` // "" | "partial" | "full"
}

// RoutingPlan is the resolved delivery plan; Destinations is never empty after a successful Resolve.
type RoutingPlan struct {
	PlanID       string        `json:"plan_id"`
	Destinations []Destination `json:"destinations"`
	Metadata     PlanMetadata  `json:"metadata"`
}

type destKey struct {
	channel Channel
	target  TargetKind
}

func destKeyOf(d Destination) destKey {
	return destKey{channel: d.Channel, target: d.Target}
}

// MergeDestinations merges by (channel, target). Required is OR-ed; Priority keeps the maximum.
func MergeDestinations(parts ...[]Destination) []Destination {
	byKey := make(map[destKey]Destination)
	for _, list := range parts {
		for _, d := range list {
			k := destKeyOf(d)
			if existing, ok := byKey[k]; ok {
				if d.Priority > existing.Priority {
					existing.Priority = d.Priority
				}
				existing.Required = existing.Required || d.Required
				byKey[k] = existing
			} else {
				byKey[k] = d
			}
		}
	}
	out := make([]Destination, 0, len(byKey))
	for _, d := range byKey {
		out = append(out, d)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Channel != out[j].Channel {
			return out[i].Channel < out[j].Channel
		}
		return out[i].Target < out[j].Target
	})
	return out
}

// EnsureNonEmpty unions fallback destinations when the merged list would be empty.
func EnsureNonEmpty(merged []Destination, fallback []Destination) []Destination {
	if len(merged) >= 1 {
		return merged
	}
	return MergeDestinations(fallback)
}
