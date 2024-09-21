package middlewares

import (
	"context"
	"fmt"
	"gopher-social-backend-server/pkg/constants"
	"gopher-social-backend-server/pkg/utils"
	"net/http"
	"time"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("AuthToken")
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, "authToken cookie not found")
			return
		}

		if !utils.ValidateAccessToken(cookie.Value) {
			http.SetCookie(w, &http.Cookie{
				Name:     "AuthToken",
				Value:    "",
				Expires:  time.Unix(0, 0),
				HttpOnly: true,
			})
			utils.WriteError(w, http.StatusUnauthorized, "token is expired or invalid")
			return
		}

		claims, err := utils.ParseAccessToken(cookie.Value)
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, fmt.Sprintf("error parsing token: %v", err))
			return
		}

		userID := claims.Subject
		if userID == "" {
			utils.WriteError(w, http.StatusUnauthorized, "invalid token claims: userID is empty")
			return
		}

		ctx := context.WithValue(r.Context(), constants.UserIDKey, userID)
		r = r.WithContext(ctx)

		userIDFromContext := r.Context().Value(constants.UserIDKey)
		if userIDFromContext == nil || userIDFromContext == "" {
			utils.WriteError(w, http.StatusUnauthorized, "userID not found in context")
			return
		}

		next.ServeHTTP(w, r)
	})
}
