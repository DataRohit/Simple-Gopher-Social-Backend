package health

import "github.com/go-chi/chi/v5"

func RegisterHealthRoutes(router chi.Router, handler *HealthHandler) {
	router.Get("/health/router", handler.GetRouterHealthHandler)
}
