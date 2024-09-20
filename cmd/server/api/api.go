package api

import (
	"context"
	"gopher-social-backend-server/internal/middlewares"
	"gopher-social-backend-server/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var log = logger.GetLogger()

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
