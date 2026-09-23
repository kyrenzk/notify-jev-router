package notifyjev

import (
	"context"

	jev "github.com/havlan/jev-go"
)

// JevClient abstracts the TypeSafe Jev SystemOne API for tests and adapters.
type JevClient interface {
	SystemOne(ctx context.Context, state any, questions jev.Questions) (*jev.Response, error)
}
