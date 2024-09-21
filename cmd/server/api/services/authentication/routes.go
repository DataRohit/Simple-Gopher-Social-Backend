package authentication

import "github.com/go-chi/chi/v5"

func RegisterAuthenticationRoutes(router chi.Router, handler *AuthenticationHandler) {
	router.Post("/auth/register", handler.RegisterUserHandler)
	router.Get("/auth/activate/{token}", handler.ActivateUserHandler)
	router.Post("/auth/login", handler.LoginUserHandler)
	router.Post("/auth/logout", handler.LogoutUserHandler)
	router.Post("/auth/forgot-password", handler.ForgotPasswordHandler)
	router.Post("/auth/reset-password/{token}", handler.ResetPasswordHandler)
	router.Get("/auth/google/login", handler.GoogleLoginHandler)
	router.Get("/auth/google/callback", handler.GoogleCallbackHandler)
	router.Get("/auth/github/login", handler.GitHubLoginHandler)
	router.Get("/auth/github/callback", handler.GitHubCallbackHandler)
}
