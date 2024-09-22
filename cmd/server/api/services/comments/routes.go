package comments

import (
	"gopher-social-backend-server/internal/middlewares"

	"github.com/go-chi/chi/v5"
)

func RegisterCommentsRoutes(router chi.Router, handler *CommentsHandler) {
	router.Get("/comments/{commentID}", handler.GetCommentByIDHandler)
	router.With(middlewares.PaginationMiddleware, middlewares.OrderingMiddleware).Get("/posts/{postID}/comments", handler.GetCommentsForPostHandler)
	router.With(middlewares.AuthMiddleware).Post("/posts/{postID}/comments", handler.CreateCommentHandler)
	router.With(middlewares.AuthMiddleware).Put("/comments/{commentID}", handler.UpdateCommentHandler)
	router.With(middlewares.AuthMiddleware).Delete("/comments/{commentID}", handler.DeleteCommentHandler)
	router.With(middlewares.AuthMiddleware).Post("/comments/{commentID}/like", handler.LikeCommentHandler)
	router.With(middlewares.AuthMiddleware).Delete("/comments/{commentID}/like", handler.UnlikeCommentHandler)
	router.With(middlewares.AuthMiddleware).Post("/comments/{commentID}/dislike", handler.DislikeCommentHandler)
	router.With(middlewares.AuthMiddleware).Delete("/comments/{commentID}/dislike", handler.UndislikeCommentHandler)
}
