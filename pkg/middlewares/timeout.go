package middlewares

import (
	"app/product-api/pkg/utils"
	"net/http"
)

type TimeoutMiddleware struct {
	TimeoutInSeconds float64
}

func NewTimeoutMiddleware(seconds float64) *TimeoutMiddleware {
	return &TimeoutMiddleware{TimeoutInSeconds: seconds}
}

func (t *TimeoutMiddleware) AddTimeout(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx, cancel := utils.AddRequestTimeout(req.Context(), t.TimeoutInSeconds)
		defer cancel()

		next.ServeHTTP(w, req.WithContext(ctx))
	})
}
