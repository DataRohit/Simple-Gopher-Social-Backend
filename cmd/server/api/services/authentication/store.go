package authentication

import "gorm.io/gorm"

type AuthenticationStore interface {
	CreateUser(user *User) error
}

type authenticationStore struct {
	postgresDB *gorm.DB
}

func NewAuthenticationStore(postgresDB *gorm.DB) AuthenticationStore {
	return &authenticationStore{
		postgresDB: postgresDB,
	}
}

func (s *authenticationStore) CreateUser(user *User) error {
	return s.postgresDB.Create(user).Error
}
