package api

import (
	"context"
	"gopher-social-backend-server/cmd/server/api/services/authentication"
	"gopher-social-backend-server/internal/database"
	"gopher-social-backend-server/internal/middlewares"
	"gopher-social-backend-server/pkg/logger"
	"gopher-social-backend-server/pkg/ratelimiter"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var log = logger.GetLogger()

func (app *Application) mountRoutes(router chi.Router) {
	router.Get("/health/router", app.Handlers.HealthHandler.GetRouterHealthHandler)

	router.Post("/auth/register", app.Handlers.AuthenticationHandler.RegisterUserHandler)
	router.Get("/auth/activate/{token}", app.Handlers.AuthenticationHandler.ActivateUserHandler)
	router.Post("/auth/login", app.Handlers.AuthenticationHandler.LoginUserHandler)
	router.Post("/auth/logout", app.Handlers.AuthenticationHandler.LogoutUserHandler)
	router.Post("/auth/forgot-password", app.Handlers.AuthenticationHandler.ForgotPasswordHandler)
	router.Post("/auth/reset-password/{token}", app.Handlers.AuthenticationHandler.ResetPasswordHandler)
}

func (app *Application) prepareDatabase() {
	if err := app.PostgresDB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		log.Error("could not create extension", zap.String("extension", "uuid-ossp"), zap.Error(err))
	}

	if err := app.PostgresDB.Exec(`DO $$
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
			CREATE TYPE user_role AS ENUM ('user', 'staff', 'admin');
		END IF;
	END $$;`).Error; err != nil {
		log.Error("could not create type", zap.String("type", "user_role"), zap.Error(err))
	}
}

func (app *Application) makeMigrations() {
	app.prepareDatabase()

	if err := database.MigrateModel(&authentication.User{}); err != nil {
		log.Error("could not migrate model", zap.String("model", "User"), zap.Error(err))
	}
}

func (app *Application) configureRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middlewares.LoggingMiddleware)
	router.Use(middlewares.RequestIDMiddleware)
	router.Use(middlewares.RealIPMiddleware)
	router.Use(middlewares.RecovererMiddleware)
	router.Use(middlewares.CORSMiddleware(
		[]string{},
		[]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		[]string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		[]string{"Link"},
		false,
		300,
	))

	rateLimiter := ratelimiter.NewRateLimiter(time.Second)
	router.Use(middlewares.RateLimiterMiddleware(rateLimiter))
	router.Use(middlewares.TimeoutMiddleware(time.Minute))

	app.makeMigrations()
	app.mountRoutes(router)

	return router
}

func (app *Application) Run() {
	router := app.configureRouter()

	server := &http.Server{
		Addr:         app.Config.Address,
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  2 * time.Minute,
	}

	go func() {
		log.Info("starting server", zap.String("address", app.Config.Address))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("could not listen on address", zap.String("address", app.Config.Address), zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("could not gracefully shutdown the server", zap.Error(err))
	} else {
		log.Info("server shutdown gracefully")
	}
}
