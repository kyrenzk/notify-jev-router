//go:build integration

package notifyjev

import (
	"context"
	"os"
	"strings"
	"testing"

	jev "github.com/havlan/jev-go"
)

func TestIntegration_security_new_login_resolve(t *testing.T) {
	apiKey := strings.TrimSpace(os.Getenv("TYPESAFE_API_KEY"))
	if apiKey == "" {
		t.Skip("TYPESAFE_API_KEY not set")
	}
	model := strings.TrimSpace(os.Getenv("JEV_MODEL"))
	if model == "" {
		t.Fatal("JEV_MODEL must be set to a pinned model (not jev-latest) for integration tests")
	}
	if model == "jev-latest" {
		t.Fatal("JEV_MODEL must be pinned; jev-latest is not allowed in integration tests")
	}

	jc, err := jev.New(jev.Config{APIKey: apiKey, Model: model})
	if err != nil {
		t.Fatal(err)
	}
	router, err := New(Config{
		JevClient: NewJevGoAdapter(jc, model),
		Model:     model,
	})
	if err != nil {
		t.Fatal(err)
	}

	rc, err := ValidateRoutingContextJSON([]byte(validContextJSON()))
	if err != nil {
		t.Fatal(err)
	}
	plan, err := router.Resolve(context.Background(), ResolveInput{Context: rc})
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Destinations) < 1 {
		t.Fatal("expected destinations")
	}
	assertHas(t, plan.Destinations, ChannelEmail, TargetRecovery, true)
}
