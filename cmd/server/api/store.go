package api

import (
	"gopher-social-backend-server/cmd/server/api/services/authentication"
	"gopher-social-backend-server/cmd/server/api/services/comments"
	"gopher-social-backend-server/cmd/server/api/services/posts"

	"gorm.io/gorm"
)

type Store struct {
	AuthenticationStore authentication.AuthenticationStore
	PostsStore          posts.PostsStore
	CommentsStore       comments.CommentsStore
}

func NewStore(postgresDB *gorm.DB) *Store {
	return &Store{
		AuthenticationStore: authentication.NewAuthenticationStore(postgresDB),
		PostsStore:          posts.NewPostsStore(postgresDB),
		CommentsStore:       comments.NewCommentStore(postgresDB),
	}
}
