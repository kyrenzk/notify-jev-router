package notifyjev

// PolicySet applies hard must-rules before Jev interpretation.
type PolicySet struct {
	Rules []PolicyRule
}

// PolicyRule adds destinations and records override identifiers for metadata.
type PolicyRule struct {
	ID      string
	Applies func(RoutingContext) bool
	Add     func(RoutingContext) []Destination
}

// DefaultPolicies is the v0.1.0 built-in must-rule set.
var DefaultPolicies = PolicySet{
	Rules: []PolicyRule{
		{
			ID: "security_recovery_email",
			Applies: func(rc RoutingContext) bool {
				return rc.Category == CategorySecurity && rc.Capabilities.HasRecoveryEmail
			},
			Add: func(rc RoutingContext) []Destination {
				return []Destination{{
					Channel:  ChannelEmail,
					Target:   TargetRecovery,
					Required: true,
					Priority: 100,
				}}
			},
		},
		{
			ID: "billing_critical_primary_email",
			Applies: func(rc RoutingContext) bool {
				return rc.Category == CategoryBilling && rc.Severity == SeverityCritical
			},
			Add: func(rc RoutingContext) []Destination {
				return []Destination{{
					Channel:  ChannelEmail,
					Target:   TargetPrimary,
					Required: true,
					Priority: 100,
				}}
			},
		},
	},
}

// Apply evaluates all matching rules.
func (ps PolicySet) Apply(rc RoutingContext) ([]Destination, []string) {
	var dests []Destination
	var overrides []string
	for _, rule := range ps.Rules {
		if rule.Applies != nil && rule.Applies(rc) {
			if rule.Add != nil {
				dests = append(dests, rule.Add(rc)...)
			}
			if rule.ID != "" {
				overrides = append(overrides, rule.ID)
			}
		}
	}
	return dests, overrides
}
