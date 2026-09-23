package notifyjev

import (
	"context"
	"testing"

	jev "github.com/havlan/jev-go"
)

func TestShadowResolve_match_legacy_security_login(t *testing.T) {
	client := &mockJevClient{fn: func(ctx context.Context, state any, questions jev.Questions) (*jev.Response, error) {
		return loadJevFixture("security_new_login.json"), nil
	}}
	r := testRouter(client)
	rc := sampleSecurityContext(t)
	legacy := []Destination{
		{Channel: ChannelEmail, Target: TargetPrimary},
		{Channel: ChannelEmail, Target: TargetRecovery, Required: true},
		{Channel: ChannelPush, Target: TargetAllDevices},
	}
	obs, err := r.ShadowResolve(context.Background(), ShadowInput{Context: rc, Legacy: legacy})
	if err != nil {
		t.Fatal(err)
	}
	if !obs.Diff.Match {
		t.Fatalf("expected match, diff=%+v", obs.Diff)
	}
}

func TestShadowMetrics_category_rollup(t *testing.T) {
	m := &ShadowMetrics{}
	obsMatch := ShadowObservation{Diff: PlanDiff{Match: true}}
	obsMiss := ShadowObservation{Diff: PlanDiff{Match: false, Added: []Destination{{Channel: ChannelPush, Target: TargetAllDevices}}}}
	m.RecordShadowObservation(CategorySecurity, obsMatch)
	m.RecordShadowObservation(CategorySecurity, obsMiss)
	if m.Total != 2 || m.Matches != 1 || m.Mismatches != 1 {
		t.Fatalf("metrics=%+v", m)
	}
	st := m.ByCategory[CategorySecurity]
	if st.Total != 2 || st.Matches != 1 {
		t.Fatalf("category stats=%+v", st)
	}
}
