package middleware

import (
	"net/http"
	"time"
	"weather-api/internal/response"

	"github.com/go-chi/httprate"
)

func RateLimit(requestsPerWindow int, window time.Duration) func(http.Handler) http.Handler {
	return httprate.Limit(
		requestsPerWindow,
		window,
		httprate.WithKeyFuncs(httprate.KeyByIP),
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			response.Error(w, http.StatusTooManyRequests, "rate limit exceeded")
		}),
	)
}
