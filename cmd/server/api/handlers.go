package api

import (
	"gopher-social-backend-server/cmd/server/api/services/authentication"
	"gopher-social-backend-server/cmd/server/api/services/health"
)

type Handlers struct {
	HealthHandler         *health.HealthHandler
	AuthenticationHandler *authentication.AuthenticationHandler
}

func NewHandlers(store *Store) *Handlers {
	return &Handlers{
		HealthHandler:         &health.HealthHandler{},
		AuthenticationHandler: &authentication.AuthenticationHandler{AuthenticationStore: store.AuthenticationStore},
	}
}
