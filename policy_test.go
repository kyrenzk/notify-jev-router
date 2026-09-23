package notifyjev

import (
	"testing"
)

func TestPolicy_SecurityRecoveryRequired(t *testing.T) {
	rc := sampleSecurityContext(t)
	dests, overrides := DefaultPolicies.Apply(rc)
	if len(dests) != 1 {
		t.Fatalf("dests=%v", dests)
	}
	if dests[0].Channel != ChannelEmail || dests[0].Target != TargetRecovery || !dests[0].Required {
		t.Fatalf("unexpected dest: %+v", dests[0])
	}
	if len(overrides) != 1 || overrides[0] != "security_recovery_email" {
		t.Fatalf("overrides=%v", overrides)
	}
}
