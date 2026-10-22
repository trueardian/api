package timezone

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPWithZone(t *testing.T) {
	tests := []struct {
		name       string
		header     string
		setHeader  bool
		wantStatus int
		wantNext   bool
	}{
		{name: "valid", header: "Asia/Jakarta", setHeader: true, wantStatus: http.StatusOK, wantNext: true},
		{name: "absent", setHeader: false, wantStatus: http.StatusBadRequest, wantNext: false},
		{name: "local", header: "Local", setHeader: true, wantStatus: http.StatusBadRequest, wantNext: false},
		{name: "invalid", header: "Mars/Olympus", setHeader: true, wantStatus: http.StatusBadRequest, wantNext: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				nextCalled bool
				gotZone    string
			)
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				if loc, ok := ZoneFromContext(r.Context()); ok {
					gotZone = loc.String()
				}
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.setHeader {
				req.Header.Set("time-zone", tt.header)
			}
			rec := httptest.NewRecorder()

			HTTPWithZone(next).ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if nextCalled != tt.wantNext {
				t.Errorf("next called = %v, want %v", nextCalled, tt.wantNext)
			}
			if tt.wantNext && gotZone != tt.header {
				t.Errorf("context zone = %q, want %q", gotZone, tt.header)
			}
			if !tt.wantNext && rec.Body.Len() == 0 {
				t.Error("expected a problem body on rejection, got empty")
			}
		})
	}
}

func TestPropagateZone(t *testing.T) {
	ctx, err := ContextWithZone(context.Background(), "Asia/Jakarta")
	if err != nil {
		t.Fatalf("ContextWithZone: %v", err)
	}

	h := http.Header{}
	PropagateZone(ctx, h)
	if got := h.Get("time-zone"); got != "Asia/Jakarta" {
		t.Errorf("header = %q, want Asia/Jakarta", got)
	}

	h = http.Header{}
	PropagateZone(context.Background(), h)
	if len(h) != 0 {
		t.Errorf("header set on empty context: %v", h)
	}
}
