package api

import (
	"context"
	"net/http"
	"testing"
)

type stubRoundTripper struct {
	got       *http.Request
	resp      *http.Response
	closeIdle bool
}

func (s *stubRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	s.got = r
	return s.resp, nil
}

func (s *stubRoundTripper) CloseIdleConnections() {
	s.closeIdle = true
}

func TestTransport_RoundTrip(t *testing.T) {
	stub := &stubRoundTripper{resp: &http.Response{StatusCode: http.StatusOK}}
	tr := &Transport{
		Base: stub,
		Propagators: []Propagator{
			func(ctx context.Context, h http.Header) { h.Set("user-id", "u_123") },
		},
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.test/", nil)
	if err != nil {
		t.Fatalf("NewRequestWithContext: %v", err)
	}

	if _, err := tr.RoundTrip(req); err != nil {
		t.Fatalf("RoundTrip: %v", err)
	}

	if got := req.Header.Get("user-id"); got != "" {
		t.Errorf("original request was mutated: user-id = %q, want unset", got)
	}
	if stub.got == nil {
		t.Fatal("base round tripper was not invoked")
	}
	if stub.got == req {
		t.Error("base received the original request, want a clone")
	}
	if got := stub.got.Header.Get("user-id"); got != "u_123" {
		t.Errorf("propagated header = %q, want u_123", got)
	}
}

func TestTransport_CloseIdleConnections(t *testing.T) {
	stub := &stubRoundTripper{}
	tr := &Transport{Base: stub}

	(&http.Client{Transport: tr}).CloseIdleConnections()

	if !stub.closeIdle {
		t.Error("CloseIdleConnections was not forwarded to Base")
	}
}
