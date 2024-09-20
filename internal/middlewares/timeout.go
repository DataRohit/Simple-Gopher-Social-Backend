package middlewares

import (
	"context"
	"gopher-social-backend-server/pkg/utils"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func TimeoutMiddleware(duration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), duration)
			defer cancel()

			r = r.WithContext(ctx)

			finished := make(chan struct{})
			go func() {
				next.ServeHTTP(w, r)
				close(finished)
			}()

			select {
			case <-finished:
				return
			case <-ctx.Done():
				if ctx.Err() == context.DeadlineExceeded {
					utils.WriteError(w, http.StatusGatewayTimeout, "request timed out")
					log.Warn("request timed out",
						zap.String("method", r.Method),
						zap.String("url", r.URL.Path),
						zap.Duration("timeout", duration))
				}
			}
		})
	}
}
