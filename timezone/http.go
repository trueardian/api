package timezone

import (
	"context"
	"net/http"

	"go.trueardian.com/api"
)

const header = "time-zone"

// HTTPWithZone stores the time-zone header in the request context, or responds 400 when it is not an IANA zone.
func HTTPWithZone(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, err := ContextWithZone(r.Context(), r.Header.Get(header))
		if err != nil {
			api.WriteProblem(w, http.StatusBadRequest, map[string]any{
				"detail": err.Error(),
			})
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// PropagateZone sets the time-zone header from ctx, if present.
func PropagateZone(ctx context.Context, h http.Header) {
	if loc, ok := ZoneFromContext(ctx); ok {
		h.Set(header, loc.String())
	}
}
