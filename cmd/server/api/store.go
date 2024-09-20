package api

import (
	"gopher-social-backend-server/cmd/server/api/services/authentication"

	"gorm.io/gorm"
)

type Store struct {
	AuthenticationStore authentication.AuthenticationStore
}

func NewStore(postgresDB *gorm.DB) *Store {
	return &Store{
		AuthenticationStore: authentication.NewAuthenticationStore(postgresDB),
	}
}
