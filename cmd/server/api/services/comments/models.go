package comments

import (
	"gopher-social-backend-server/cmd/server/api/services/authentication"
	"gopher-social-backend-server/cmd/server/api/services/posts"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID           `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	AuthorID  uuid.UUID           `json:"author_id" gorm:"type:uuid;not null;index"`
	Author    authentication.User `json:"author" gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	PostID    uuid.UUID           `json:"post_id" gorm:"type:uuid;not null;index"`
	Post      posts.Post          `json:"post" gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Content   string              `json:"content" gorm:"type:text;not null"`
	CreatedAt int64               `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt int64               `json:"updated_at" gorm:"autoUpdateTime"`
}

type CommentLike struct {
	ID        uuid.UUID           `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID           `json:"user_id" gorm:"type:uuid;not null;index"`
	User      authentication.User `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CommentID uuid.UUID           `json:"comment_id" gorm:"type:uuid;not null;index"`
	Comment   Comment             `json:"comment" gorm:"foreignKey:CommentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type CommentDislike struct {
	ID        uuid.UUID           `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID           `json:"user_id" gorm:"type:uuid;not null;index"`
	User      authentication.User `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CommentID uuid.UUID           `json:"comment_id" gorm:"type:uuid;not null;index"`
	Comment   Comment             `json:"comment" gorm:"foreignKey:CommentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
