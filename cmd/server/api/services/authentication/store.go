package authentication

import "gorm.io/gorm"

type AuthenticationStore interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id string) (*User, error)
	UpdateUser(user *User) error
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

func (s *authenticationStore) GetUserByEmail(email string) (*User, error) {
	var user User
	if err := s.postgresDB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *authenticationStore) GetUserByID(id string) (*User, error) {
	var user User
	if err := s.postgresDB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *authenticationStore) UpdateUser(user *User) error {
	return s.postgresDB.Save(user).Error
}
