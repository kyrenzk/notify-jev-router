package notifyjev

// Category identifies the v0.1.0 routing policy bucket (security, billing, transactional).
type Category string

const (
	CategorySecurity      Category = "security"
	CategoryBilling       Category = "billing"
	CategoryTransactional Category = "transactional"
)

// Severity classifies urgency for policy and fallback defaults.
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityNormal   Severity = "normal"
	SeverityLow      Severity = "low"
)

// Channel is a delivery channel identifier.
type Channel string

const (
	ChannelEmail Channel = "email"
	ChannelPush  Channel = "push"
	ChannelSMS   Channel = "sms"
	ChannelInApp Channel = "in_app"
)

// TargetKind describes which logical endpoint within a channel to use.
type TargetKind string

const (
	TargetPrimary          TargetKind = "primary"
	TargetRecovery         TargetKind = "recovery"
	TargetAllDevices       TargetKind = "all_devices"
	TargetLastActiveDevice TargetKind = "last_active_device"
)

func (c Category) valid() bool {
	switch c {
	case CategorySecurity, CategoryBilling, CategoryTransactional:
		return true
	default:
		return false
	}
}

func (s Severity) valid() bool {
	switch s {
	case SeverityCritical, SeverityNormal, SeverityLow:
		return true
	default:
		return false
	}
}
