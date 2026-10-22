package user

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPWithID(t *testing.T) {
	tests := []struct {
		name       string
		header     string
		wantStatus int
		wantNext   bool
	}{
		{name: "present", header: "u_123", wantStatus: http.StatusOK, wantNext: true},
		{name: "absent", header: "", wantStatus: http.StatusUnauthorized, wantNext: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				nextCalled bool
				gotID      string
			)
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				gotID, _ = IDFromContext(r.Context())
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.header != "" {
				req.Header.Set("user-id", tt.header)
			}
			rec := httptest.NewRecorder()

			HTTPWithID(next).ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if nextCalled != tt.wantNext {
				t.Errorf("next called = %v, want %v", nextCalled, tt.wantNext)
			}
			if tt.wantNext && gotID != tt.header {
				t.Errorf("context id = %q, want %q", gotID, tt.header)
			}
			if !tt.wantNext && rec.Body.Len() == 0 {
				t.Error("expected a problem body on rejection, got empty")
			}
		})
	}
}

func TestPropagateID(t *testing.T) {
	ctx, err := ContextWithID(context.Background(), "u_123")
	if err != nil {
		t.Fatalf("ContextWithID: %v", err)
	}

	h := http.Header{}
	PropagateID(ctx, h)
	if got := h.Get("user-id"); got != "u_123" {
		t.Errorf("header = %q, want u_123", got)
	}

	h = http.Header{}
	PropagateID(context.Background(), h)
	if len(h) != 0 {
		t.Errorf("header set on empty context: %v", h)
	}
}
