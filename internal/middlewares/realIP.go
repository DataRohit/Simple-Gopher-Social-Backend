package middlewares

import "net/http"

func RealIPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		realIP := r.Header.Get("X-Forwarded-For")
		if realIP == "" {
			realIP = r.RemoteAddr
		}
		r.Header.Set("X-Real-IP", realIP)
		next.ServeHTTP(w, r)
	})
}
