package notifyjev

// ConfidenceConfig tunes per-question confidence gates.
type ConfidenceConfig struct {
	ChoiceMinConfidence float64
	ChoiceMinMargin     float64
	NoulThresholds      map[string]float64
}

// DefaultConfidenceConfig is the v0.1.0 default gate set.
func DefaultConfidenceConfig() ConfidenceConfig {
	return ConfidenceConfig{
		ChoiceMinConfidence: 0.55,
		ChoiceMinMargin:     0.15,
		NoulThresholds: map[string]float64{
			"recovery_copy":     0.65,
			"email_mandatory":   0.60,
		},
	}
}

func (cc ConfidenceConfig) noulThreshold(question string) float64 {
	if cc.NoulThresholds != nil {
		if v, ok := cc.NoulThresholds[question]; ok {
			return v
		}
	}
	return 0.60
}
