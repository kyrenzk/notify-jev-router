package notifyjev

// FallbackTable maps category and severity to safe default destinations.
type FallbackTable map[Category]map[Severity][]Destination

// DefaultFallbackTable is embedded v0.1.0 defaults; each cell has len >= 1.
var DefaultFallbackTable = FallbackTable{
	CategorySecurity: {
		SeverityCritical: []Destination{
			{Channel: ChannelEmail, Target: TargetPrimary, Priority: 50},
			{Channel: ChannelPush, Target: TargetAllDevices, Priority: 40},
		},
		SeverityNormal: []Destination{
			{Channel: ChannelEmail, Target: TargetPrimary, Priority: 50},
			{Channel: ChannelPush, Target: TargetAllDevices, Priority: 40},
		},
		SeverityLow: []Destination{
			{Channel: ChannelEmail, Target: TargetPrimary, Priority: 50},
			{Channel: ChannelPush, Target: TargetAllDevices, Priority: 40},
		},
	},
	CategoryBilling: {
		SeverityCritical: []Destination{
			{Channel: ChannelEmail, Target: TargetPrimary, Priority: 50},
		},
		SeverityNormal: []Destination{
			{Channel: ChannelEmail, Target: TargetPrimary, Priority: 50},
		},
		SeverityLow: []Destination{
			{Channel: ChannelEmail, Target: TargetPrimary, Priority: 50},
		},
	},
	CategoryTransactional: {
		SeverityCritical: []Destination{
			{Channel: ChannelEmail, Target: TargetPrimary, Priority: 50},
			{Channel: ChannelPush, Target: TargetAllDevices, Priority: 40},
		},
		SeverityNormal: []Destination{
			{Channel: ChannelEmail, Target: TargetPrimary, Priority: 50},
			{Channel: ChannelPush, Target: TargetAllDevices, Priority: 40},
		},
		SeverityLow: []Destination{
			{Channel: ChannelEmail, Target: TargetPrimary, Priority: 50},
			{Channel: ChannelPush, Target: TargetAllDevices, Priority: 40},
		},
	},
}

// Lookup returns defaults for a category and severity.
func (ft FallbackTable) Lookup(category Category, severity Severity) []Destination {
	if ft == nil {
		ft = DefaultFallbackTable
	}
	bySev, ok := ft[category]
	if !ok {
		return []Destination{{Channel: ChannelEmail, Target: TargetPrimary, Priority: 1}}
	}
	if dests, ok := bySev[severity]; ok && len(dests) > 0 {
		return cloneDests(dests)
	}
	for _, dests := range bySev {
		if len(dests) > 0 {
			return cloneDests(dests)
		}
	}
	return []Destination{{Channel: ChannelEmail, Target: TargetPrimary, Priority: 1}}
}

func cloneDests(in []Destination) []Destination {
	out := make([]Destination, len(in))
	copy(out, in)
	return out
}

// QuestionFallbackDefaults supplies per-question defaults when confidence gates fail.
var QuestionFallbackDefaults = map[string]func(RoutingContext) []Destination{
	"primary_strategy": func(rc RoutingContext) []Destination {
		return DefaultFallbackTable.Lookup(rc.Category, rc.Severity)
	},
	"push_fanout": func(rc RoutingContext) []Destination {
		if rc.Category == CategorySecurity || rc.Category == CategoryTransactional {
			return []Destination{{Channel: ChannelPush, Target: TargetAllDevices, Priority: 30}}
		}
		return nil
	},
	"recovery_copy": func(rc RoutingContext) []Destination {
		if rc.Category == CategorySecurity && rc.Capabilities.HasRecoveryEmail {
			return []Destination{{Channel: ChannelEmail, Target: TargetRecovery, Priority: 20}}
		}
		return nil
	},
	"email_mandatory": func(rc RoutingContext) []Destination {
		return []Destination{{Channel: ChannelEmail, Target: TargetPrimary, Required: true, Priority: 60}}
	},
}

// FullFallback returns table defaults for API failure paths.
func FullFallback(rc RoutingContext, table FallbackTable) []Destination {
	return table.Lookup(rc.Category, rc.Severity)
}
