package middlewares

import (
	"gopher-social-backend-server/pkg/logger"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var log = logger.GetLogger()

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		traceID := uuid.New().String()

		log.Info("incoming request",
			zap.String("method", r.Method),
			zap.String("url", r.URL.Path),
			zap.String("trace_id", traceID))

		rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rec, r)

		log.Info("completed request",
			zap.String("method", r.Method),
			zap.String("url", r.URL.Path),
			zap.Int("status_code", rec.statusCode),
			zap.Duration("duration", time.Since(start)),
			zap.String("trace_id", traceID))
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}
