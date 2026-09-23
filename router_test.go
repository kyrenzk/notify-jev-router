package notifyjev

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	jev "github.com/havlan/jev-go"
)

type mockJevClient struct {
	fn func(ctx context.Context, state any, questions jev.Questions) (*jev.Response, error)
}

func (m *mockJevClient) SystemOne(ctx context.Context, state any, questions jev.Questions) (*jev.Response, error) {
	return m.fn(ctx, state, questions)
}

func loadJevFixture(name string) *jev.Response {
	path := filepath.Join("testdata", "jev", name)
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var resp jev.Response
	if err := json.Unmarshal(data, &resp); err != nil {
		panic(err)
	}
	return &resp
}

func testRouter(client JevClient) *Router {
	r, err := New(Config{
		JevClient: client,
		Model:     "jev-test-1.0.0",
	})
	if err != nil {
		panic(err)
	}
	return r
}

func TestAcceptance_security_new_login(t *testing.T) {
	client := &mockJevClient{fn: func(ctx context.Context, state any, questions jev.Questions) (*jev.Response, error) {
		return loadJevFixture("security_new_login.json"), nil
	}}
	r := testRouter(client)
	rc := sampleSecurityContext(t)
	plan, err := r.Resolve(context.Background(), ResolveInput{Context: rc})
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Destinations) < 1 {
		t.Fatal("expected destinations")
	}
	assertHas(t, plan.Destinations, ChannelEmail, TargetPrimary, false)
	assertHas(t, plan.Destinations, ChannelEmail, TargetRecovery, true)
	assertHas(t, plan.Destinations, ChannelPush, TargetAllDevices, false)
}

func TestAcceptance_JevAPIFailure_billing_critical_fallback(t *testing.T) {
	client := &mockJevClient{fn: func(ctx context.Context, state any, questions jev.Questions) (*jev.Response, error) {
		return nil, errors.New("jev unavailable")
	}}
	r := testRouter(client)
	rc := RoutingContext{
		EventID:    "evt-bill",
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
	plan, err := r.Resolve(context.Background(), ResolveInput{Context: rc})
	if err != nil {
		t.Fatal(err)
	}
	if plan.Metadata.FallbackMode != "full" {
		t.Fatalf("fallback mode=%q", plan.Metadata.FallbackMode)
	}
	assertHas(t, plan.Destinations, ChannelEmail, TargetPrimary, true)
}

func TestAcceptance_PolicyMust_security_recovery_survives_jev(t *testing.T) {
	client := &mockJevClient{fn: func(ctx context.Context, state any, questions jev.Questions) (*jev.Response, error) {
		return loadJevFixture("security_push_only.json"), nil
	}}
	r := testRouter(client)
	rc := sampleSecurityContext(t)
	plan, err := r.Resolve(context.Background(), ResolveInput{Context: rc})
	if err != nil {
		t.Fatal(err)
	}
	assertHas(t, plan.Destinations, ChannelEmail, TargetRecovery, true)
}

func assertHas(t *testing.T, dests []Destination, ch Channel, target TargetKind, required bool) {
	for _, d := range dests {
		if d.Channel == ch && d.Target == target {
			if d.Required != required {
				t.Fatalf("required=%v want %v for %s/%s", d.Required, required, ch, target)
			}
			return
		}
	}
	t.Fatalf("missing %s/%s in %v", ch, target, dests)
}
