package api

import (
	"gopher-social-backend-server/cmd/server/api/services/authentication"
	"gopher-social-backend-server/cmd/server/api/services/comments"
	"gopher-social-backend-server/cmd/server/api/services/health"
	"gopher-social-backend-server/cmd/server/api/services/posts"
)

type Handlers struct {
	HealthHandler         *health.HealthHandler
	AuthenticationHandler *authentication.AuthenticationHandler
	PostsHandler          *posts.PostsHandler
	CommentsHandler       *comments.CommentsHandler
}

func NewHandlers(store *Store) *Handlers {
	return &Handlers{
		HealthHandler:         &health.HealthHandler{},
		AuthenticationHandler: &authentication.AuthenticationHandler{AuthenticationStore: store.AuthenticationStore},
		PostsHandler:          &posts.PostsHandler{PostsStore: store.PostsStore, AuthenticationStore: store.AuthenticationStore},
		CommentsHandler:       &comments.CommentsHandler{CommentsStore: store.CommentsStore, PostsStore: store.PostsStore, AuthenticationStore: store.AuthenticationStore},
	}
}
