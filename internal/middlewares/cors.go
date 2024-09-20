package middlewares

import (
	"net/http"
	"strconv"
	"strings"
)

func CORSMiddleware(allowedOrigins []string, allowedMethods []string, allowedHeaders []string, exposedHeaders []string, allowCredentials bool, maxAge int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == "" {
				origin = "*"
			}

			if contains(allowedOrigins, origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ","))
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ","))
				w.Header().Set("Access-Control-Expose-Headers", strings.Join(exposedHeaders, ","))

				if allowCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}

				w.Header().Set("Access-Control-Max-Age", strconv.Itoa(maxAge))
			}

			if r.Method == http.MethodOptions {
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}
