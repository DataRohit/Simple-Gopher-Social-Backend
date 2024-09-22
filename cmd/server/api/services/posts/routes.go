package posts

import (
	"gopher-social-backend-server/internal/middlewares"

	"github.com/go-chi/chi/v5"
)

func RegisterPostsRoutes(router chi.Router, handler *PostsHandler) {
	router.With(middlewares.AuthMiddleware).Post("/posts", handler.CreatePostHandler)
	router.With(middlewares.AuthMiddleware).Patch("/posts/{postID}", handler.UpdatePostByIDHandler)
	router.With(middlewares.AuthMiddleware).Delete("/posts/{postID}", handler.DeletePostByIDHandler)
}
