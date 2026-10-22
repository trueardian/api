package user

import (
	"context"
	"net/http"

	"go.trueardian.com/api"
)

const header = "user-id"

// HTTPWithID stores the user-id header in the request context, or responds 401 when it is empty.
func HTTPWithID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, err := ContextWithID(r.Context(), r.Header.Get(header))
		if err != nil {
			api.WriteProblem(w, http.StatusUnauthorized, map[string]any{
				"detail": err.Error(),
			})
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// PropagateID sets the user-id header from ctx, if present.
func PropagateID(ctx context.Context, h http.Header) {
	if id, ok := IDFromContext(ctx); ok {
		h.Set(header, id)
	}
}
