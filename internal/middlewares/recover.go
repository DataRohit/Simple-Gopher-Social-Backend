package middlewares

import (
	"gopher-social-backend-server/pkg/utils"
	"net/http"

	"go.uber.org/zap"
)

func RecovererMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("recovered from panic", zap.Error(err.(error)))
				utils.WriteError(w, http.StatusInternalServerError, err.(error).Error())
			}
		}()
		next.ServeHTTP(w, r)
	})
}
