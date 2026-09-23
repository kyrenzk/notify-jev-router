package notifyjev

import (
	"sort"

	jev "github.com/havlan/jev-go"
)

// InterpretResult holds destinations inferred from Jev answers and gate failures.
type InterpretResult struct {
	Destinations      []Destination
	AnswersSnapshot   map[string]AnswerSummary
	FallbackQuestions []string
}

// InterpretAnswers maps Jev answers to destinations with confidence gating.
func InterpretAnswers(rc RoutingContext, resp *jev.Response, asked []string, cc ConfidenceConfig) InterpretResult {
	result := InterpretResult{AnswersSnapshot: make(map[string]AnswerSummary)}
	if resp == nil {
		return result
	}
	for _, q := range asked {
		answer, ok := resp.Answers[q]
		if !ok {
			result.FallbackQuestions = append(result.FallbackQuestions, q)
			result.Destinations = append(result.Destinations, questionFallback(rc, q)...)
			continue
		}
		result.AnswersSnapshot[q] = summarizeAnswer(answer)
		switch q {
		case "primary_strategy":
			if !choicePassesGate(answer, cc) {
				result.FallbackQuestions = append(result.FallbackQuestions, q)
				result.Destinations = append(result.Destinations, questionFallback(rc, q)...)
				continue
			}
			result.Destinations = append(result.Destinations, interpretPrimaryStrategy(answer.Choice)...)
		case "push_fanout":
			if !choicePassesGate(answer, cc) {
				result.FallbackQuestions = append(result.FallbackQuestions, q)
				result.Destinations = append(result.Destinations, questionFallback(rc, q)...)
				continue
			}
			result.Destinations = append(result.Destinations, interpretPushFanout(answer.Choice)...)
		case "recovery_copy":
			threshold := cc.noulThreshold(q)
			if answer.Noul < threshold {
				result.FallbackQuestions = append(result.FallbackQuestions, q)
				result.Destinations = append(result.Destinations, questionFallback(rc, q)...)
				continue
			}
			result.Destinations = append(result.Destinations, Destination{
				Channel: ChannelEmail, Target: TargetRecovery, Priority: 30,
			})
		case "email_mandatory":
			threshold := cc.noulThreshold(q)
			if answer.Noul < threshold {
				result.FallbackQuestions = append(result.FallbackQuestions, q)
				result.Destinations = append(result.Destinations, questionFallback(rc, q)...)
				continue
			}
			result.Destinations = append(result.Destinations, Destination{
				Channel: ChannelEmail, Target: TargetPrimary, Required: true, Priority: 70,
			})
		}
	}
	return result
}

func questionFallback(rc RoutingContext, question string) []Destination {
	fn, ok := QuestionFallbackDefaults[question]
	if !ok || fn == nil {
		return nil
	}
	return fn(rc)
}

func summarizeAnswer(a jev.Answer) AnswerSummary {
	return AnswerSummary{
		Type:          string(a.Type),
		Choice:        a.Choice,
		Noul:          a.Noul,
		Confidence:    a.Confidence,
		Probabilities: a.Probabilities,
	}
}

func choicePassesGate(a jev.Answer, cc ConfidenceConfig) bool {
	if a.Confidence < cc.ChoiceMinConfidence {
		return false
	}
	if len(a.Probabilities) < 2 {
		return true
	}
	var probs []float64
	for _, p := range a.Probabilities {
		probs = append(probs, p)
	}
	sort.Float64s(probs)
	if len(probs) < 2 {
		return true
	}
	top := probs[len(probs)-1]
	second := probs[len(probs)-2]
	return top-second >= cc.ChoiceMinMargin
}

func interpretPrimaryStrategy(choice string) []Destination {
	switch choice {
	case "email_only":
		return []Destination{{Channel: ChannelEmail, Target: TargetPrimary, Priority: 40}}
	case "push_only":
		return []Destination{{Channel: ChannelPush, Target: TargetLastActiveDevice, Priority: 40}}
	case "email_and_push":
		return []Destination{
			{Channel: ChannelEmail, Target: TargetPrimary, Priority: 40},
			{Channel: ChannelPush, Target: TargetAllDevices, Priority: 35},
		}
	case "in_app_only":
		return []Destination{{Channel: ChannelInApp, Target: TargetPrimary, Priority: 40}}
	default:
		return nil
	}
}

func interpretPushFanout(choice string) []Destination {
	switch choice {
	case "all_devices":
		return []Destination{{Channel: ChannelPush, Target: TargetAllDevices, Priority: 35}}
	case "last_active_only":
		return []Destination{{Channel: ChannelPush, Target: TargetLastActiveDevice, Priority: 35}}
	default:
		return nil
	}
}

func questionIDs(questions jev.Questions) []string {
	ids := make([]string, 0, len(questions))
	for id := range questions {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}
