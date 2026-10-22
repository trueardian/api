// Package session carries the caller's session ID from the session-id header into context.Context.
package session

import (
	"context"
	"errors"
)

type key struct{}

// ContextWithID returns a copy of ctx carrying id, which must not be empty.
func ContextWithID(ctx context.Context, id string) (context.Context, error) {
	if id == "" {
		return nil, errors.New("session ID cannot be empty")
	}
	return context.WithValue(ctx, key{}, id), nil
}

// IDFromContext returns the session ID stored in ctx, if any.
func IDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(key{}).(string)
	return id, ok
}
