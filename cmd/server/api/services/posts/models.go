package posts

import (
	"gopher-social-backend-server/cmd/server/api/services/authentication"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID           `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	AuthorID  uuid.UUID           `json:"author_id" gorm:"type:uuid;not null;index"`
	Author    authentication.User `json:"author" gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Title     string              `json:"title" gorm:"type:varchar(255);not null"`
	Content   string              `json:"content" gorm:"type:text;not null"`
	CreatedAt int64               `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64               `json:"updated_at" gorm:"autoUpdateTime"`
}

type PostLike struct {
	ID     uuid.UUID           `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID uuid.UUID           `json:"user_id" gorm:"type:uuid;not null"`
	User   authentication.User `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	PostID uuid.UUID           `json:"post_id" gorm:"type:uuid;not null"`
	Post   Post                `json:"post" gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type PostDislike struct {
	ID     uuid.UUID           `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID uuid.UUID           `json:"user_id" gorm:"type:uuid;not null"`
	User   authentication.User `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	PostID uuid.UUID           `json:"post_id" gorm:"type:uuid;not null"`
	Post   Post                `json:"post" gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
