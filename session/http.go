package session

import (
	nethttp "net/http"

	apihttp "go.trueardian.com/api/http"
)

func HTTPWithID(next nethttp.Handler) nethttp.Handler {
	return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		id := r.Header.Get("session-id")
		ctx, err := ContextWithID(r.Context(), id)
		if err != nil {
			apihttp.WriteProblem(w, nethttp.StatusUnauthorized, map[string]any{
				"detail": err.Error(),
			})
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
