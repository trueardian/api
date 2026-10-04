package time

import (
	nethttp "net/http"

	apihttp "go.trueardian.com/api/http"
)

func HTTPWithZone(next nethttp.Handler) nethttp.Handler {
	return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		tz := r.Header.Get("time-zone")
		ctx, err := ContextWithZone(r.Context(), tz)
		if err != nil {
			apihttp.WriteProblem(w, nethttp.StatusBadRequest, map[string]any{
				"detail": err.Error(),
			})
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
