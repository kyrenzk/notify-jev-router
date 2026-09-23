package notifyjev

import "errors"

// ErrInvalidContext is returned when input validation fails, no deliverable channel exists,
// or the category is outside the v0.1.0 scope.
var ErrInvalidContext = errors.New("notifyjev: invalid routing context")
