package notifyjev

import (
	"context"

	jev "github.com/havlan/jev-go"
)

// JevGoAdapter wraps github.com/havlan/jev-go with a pinned model from Config.
type JevGoAdapter struct {
	client *jev.Client
	model  string
}

// NewJevGoAdapter creates an adapter from an existing jev-go client and model name.
func NewJevGoAdapter(client *jev.Client, model string) *JevGoAdapter {
	return &JevGoAdapter{client: client, model: model}
}

// SystemOne calls Jev with the configured model.
func (a *JevGoAdapter) SystemOne(ctx context.Context, state any, questions jev.Questions) (*jev.Response, error) {
	req := jev.Request{
		State:     state,
		Model:     a.model,
		Questions: questions,
	}
	return a.client.Evaluate(ctx, req)
}
