package posts

import (
	"gopher-social-backend-server/cmd/server/api/services/authentication"

	"github.com/google/uuid"
)

type Post struct {
	ID       uuid.UUID           `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	AuthorID uuid.UUID           `json:"author_id" gorm:"type:uuid;not null;index"`
	Author   authentication.User `json:"author" gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Title    string              `json:"title" gorm:"type:varchar(255);not null"`
	Content  string              `json:"content" gorm:"type:text;not null"`
	// Likes     []Like              `json:"likes" gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	// Dislikes  []Dislike           `json:"dislikes" gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	// Comments  []Comment           `json:"comments" gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"`
	CreatedAt int64 `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64 `json:"updated_at" gorm:"autoUpdateTime"`
}
