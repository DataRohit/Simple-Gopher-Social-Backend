package api

import "gopher-social-backend-server/cmd/server/api/services/health"

type Handlers struct {
	HealthHandler *health.HealthHandler
}

func NewHandlers() *Handlers {
	return &Handlers{
		HealthHandler: &health.HealthHandler{},
	}
}
