package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

func (app *Application) configureRouter() *chi.Mux {
	router := chi.NewRouter()

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
		log.Printf("server started on %s", app.Config.Address)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not listen on %s: %v\n", app.Config.Address, err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("could not gracefully shutdown the server: %v\n", err)
	} else {
		log.Println("server shutdown gracefully")
	}
}
