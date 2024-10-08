package authentication

import (
	"gopher-social-backend-server/pkg/constants"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID               `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	FirstName   string                  `json:"first_name" gorm:"type:varchar(100);not null"`
	LastName    string                  `json:"last_name" gorm:"type:varchar(100);not null"`
	Email       string                  `json:"email" gorm:"type:varchar(100);unique;not null"`
	Password    string                  `json:"password" gorm:"type:varchar(100);not null"`
	Role        constants.UserRole      `json:"role" gorm:"type:varchar(20);not null;default:'user'"`
	IsActivated bool                    `json:"is_activated" gorm:"default:false"`
	OAuth       constants.OAuthProvider `json:"oauth" gorm:"type:varchar(20);not null;default:'none'"`
	CreatedAt   int64                   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   int64                   `json:"updated_at" gorm:"autoUpdateTime"`
}
