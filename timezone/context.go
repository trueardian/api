// Package timezone carries the caller's IANA time zone from the time-zone header into context.Context.
package timezone

import (
	"context"
	"errors"
	"time"
)

type key struct{}

// ContextWithZone returns a copy of ctx carrying the IANA zone called name.
func ContextWithZone(ctx context.Context, name string) (context.Context, error) {
	if name == "" {
		return nil, errors.New("timezone cannot be empty")
	}
	if name == "Local" {
		return nil, errors.New("timezone must be an IANA name, not Local")
	}
	loc, err := time.LoadLocation(name)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, key{}, loc), nil
}

// ZoneFromContext returns the zone stored in ctx, if any.
func ZoneFromContext(ctx context.Context) (*time.Location, bool) {
	loc, ok := ctx.Value(key{}).(*time.Location)
	return loc, ok
}
