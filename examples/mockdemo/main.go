// Mock Jev demo: prints a RoutingPlan JSON without API keys.
// Run from repo root: go run ./examples/mockdemo
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	jev "github.com/havlan/jev-go"
	notifyjev "github.com/kyrenzk/notify-jev-router"
)

func main() {
	fixture, err := os.ReadFile("testdata/jev/security_new_login.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "read fixture: %v (run from repository root)\n", err)
		os.Exit(1)
	}
	var resp jev.Response
	if err := json.Unmarshal(fixture, &resp); err != nil {
		fmt.Fprintf(os.Stderr, "parse fixture: %v\n", err)
		os.Exit(1)
	}

	client := &fixtureJevClient{resp: &resp}
	router, err := notifyjev.New(notifyjev.Config{
		JevClient: client,
		Model:     "jev-fixture-1.0.0",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "router: %v\n", err)
		os.Exit(1)
	}

	rc, err := notifyjev.ValidateRoutingContextJSON([]byte(`{
		"event_id": "evt-demo",
		"template_id": "security.new_login",
		"category": "security",
		"severity": "critical",
		"locale": "ja-JP",
		"user_prefs": {"email_marketing": false, "push_enabled": true, "quiet_hours_active": false},
		"capabilities": {"email_verified": true, "push_device_count": 2, "has_recovery_email": true, "sms_available": false, "in_app_available": false}
	}`))
	if err != nil {
		fmt.Fprintf(os.Stderr, "context: %v\n", err)
		os.Exit(1)
	}

	plan, err := router.Resolve(context.Background(), notifyjev.ResolveInput{Context: rc})
	if err != nil {
		fmt.Fprintf(os.Stderr, "resolve: %v\n", err)
		os.Exit(1)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(plan); err != nil {
		fmt.Fprintf(os.Stderr, "encode: %v\n", err)
		os.Exit(1)
	}
}

type fixtureJevClient struct {
	resp *jev.Response
}

func (c *fixtureJevClient) SystemOne(ctx context.Context, state any, questions jev.Questions) (*jev.Response, error) {
	return c.resp, nil
}
