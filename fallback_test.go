package notifyjev

import (
	"testing"
)

func TestFullFallback_BillingCritical(t *testing.T) {
	rc := RoutingContext{
		EventID:    "e1",
		TemplateID: "billing.payment_failed",
		Category:   CategoryBilling,
		Severity:   SeverityCritical,
		Locale:     "ja-JP",
		UserPrefs:  UserPrefs{PushEnabled: true},
		Capabilities: Capabilities{
			EmailVerified:   true,
			PushDeviceCount: 1,
		},
	}
	dests := FullFallback(rc, DefaultFallbackTable)
	if len(dests) < 1 {
		t.Fatal("expected len>=1")
	}
	if dests[0].Channel != ChannelEmail || dests[0].Target != TargetPrimary {
		t.Fatalf("dests=%v", dests)
	}
}
