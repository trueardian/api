// Package api carries caller identity between HTTP headers and context.Context.
package api

import (
	"context"
	"net/http"
)

// Propagator sets a header on an outbound request from a value in ctx.
type Propagator func(ctx context.Context, h http.Header)

// Transport runs Propagators on a clone of each request, then sends it with Base.
// A nil Base means http.DefaultTransport.
type Transport struct {
	Base        http.RoundTripper
	Propagators []Propagator
}

// RoundTrip implements http.RoundTripper.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	clone := req.Clone(ctx)

	for _, p := range t.Propagators {
		p(ctx, clone.Header)
	}

	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(clone)
}

// CloseIdleConnections closes idle connections on Base.
func (t *Transport) CloseIdleConnections() {
	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}
	if c, ok := base.(interface{ CloseIdleConnections() }); ok {
		c.CloseIdleConnections()
	}
}
