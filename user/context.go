// Package user carries the caller's user ID from the user-id header into context.Context.
package user

import (
	"context"
	"errors"
)

type key struct{}

// ContextWithID returns a copy of ctx carrying id, which must not be empty.
func ContextWithID(ctx context.Context, id string) (context.Context, error) {
	if id == "" {
		return nil, errors.New("user ID cannot be empty")
	}
	return context.WithValue(ctx, key{}, id), nil
}

// IDFromContext returns the user ID stored in ctx, if any.
func IDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(key{}).(string)
	return id, ok
}
