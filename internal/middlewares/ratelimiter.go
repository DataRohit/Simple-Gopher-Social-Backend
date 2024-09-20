package middlewares

import (
	"gopher-social-backend-server/pkg/ratelimiter"
	"gopher-social-backend-server/pkg/utils"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func RateLimiterMiddleware(rateLimiter *ratelimiter.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if allow, retryAfter := rateLimiter.Allow(r.RemoteAddr); !allow {
				rateLimitExceededResponse(w, retryAfter)
				log.Warn("rate limit exceeded",
					zap.String("remote_addr", r.RemoteAddr),
					zap.Duration("retry_after", retryAfter))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func rateLimitExceededResponse(w http.ResponseWriter, retryAfter time.Duration) {
	w.Header().Set("Retry-After", retryAfter.String())
	utils.WriteError(w, http.StatusTooManyRequests, "rate limit exceeded! please try again later.")
}
