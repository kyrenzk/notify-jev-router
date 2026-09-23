package notifyjev

import (
	"strings"
	"testing"
)

func TestMergeDestinations_RequiredORAndDedup(t *testing.T) {
	a := []Destination{
		{Channel: ChannelEmail, Target: TargetPrimary, Required: false, Priority: 1},
		{Channel: ChannelEmail, Target: TargetRecovery, Required: true, Priority: 2},
	}
	b := []Destination{
		{Channel: ChannelEmail, Target: TargetPrimary, Required: true, Priority: 3},
		{Channel: ChannelPush, Target: TargetAllDevices, Required: false, Priority: 1},
	}
	merged := MergeDestinations(a, b)
	if len(merged) != 3 {
		t.Fatalf("len=%d want 3", len(merged))
	}
	var primary Destination
	for _, d := range merged {
		if d.Channel == ChannelEmail && d.Target == TargetPrimary {
			primary = d
		}
	}
	if !primary.Required || primary.Priority != 3 {
		t.Fatalf("primary merge: %+v", primary)
	}
}

func TestValidateRoutingContextJSON_RejectsPII(t *testing.T) {
	_, err := ValidateRoutingContextJSON([]byte(`{"email":"a@b.c","event_id":"e"}`))
	if err == nil || !strings.Contains(err.Error(), "forbidden") {
		t.Fatalf("expected forbidden field error, got %v", err)
	}
}

func TestValidateRoutingContextJSON_RejectsUnknownField(t *testing.T) {
	raw := validContextJSON()
	raw = strings.Replace(raw, `"locale": "ja-JP"`, `"locale": "ja-JP", "extra": true`, 1)
	_, err := ValidateRoutingContextJSON([]byte(raw))
	if err == nil {
		t.Fatal("expected unknown field error")
	}
}

func TestMergePreservesLenAtLeastOneWithSeed(t *testing.T) {
	seed := []Destination{
		{Channel: ChannelEmail, Target: TargetRecovery, Required: true, Priority: 10},
	}
	merged := MergeDestinations(seed)
	if len(merged) < 1 {
		t.Fatal("merge after policy seed must be >= 1")
	}
}

func validContextJSON() string {
	return `{
		"event_id": "evt-1",
		"template_id": "security.new_login",
		"category": "security",
		"severity": "critical",
		"locale": "ja-JP",
		"user_prefs": {"email_marketing": false, "push_enabled": true, "quiet_hours_active": false},
		"capabilities": {"email_verified": true, "push_device_count": 2, "has_recovery_email": true, "sms_available": false, "in_app_available": false}
	}`
}

func sampleSecurityContext(t *testing.T) RoutingContext {
	rc, err := ValidateRoutingContextJSON([]byte(validContextJSON()))
	if err != nil {
		t.Fatalf("sample context: %v", err)
	}
	return rc
}
