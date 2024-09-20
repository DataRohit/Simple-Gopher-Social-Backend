package main

import (
	"fmt"
	"gopher-social-backend-server/cmd/server/api"
	"gopher-social-backend-server/internal/database"
	"gopher-social-backend-server/pkg/logger"
	"gopher-social-backend-server/pkg/utils"

	"go.uber.org/zap"
)

var log = logger.GetLogger()

var DEBUG = utils.GetEnvAsBool("DEBUG", false)
var PORT = utils.GetEnvAsInt("PORT", 8080)

func main() {
	appConfig := api.Config{
		Address: fmt.Sprintf(":%d", PORT),
	}

	postgresDB, err := database.NewPostgresDB()
	if err != nil {
		log.Error("failed to connect to the database", zap.Error(err))
	}

	app := &api.Application{
		Config:   appConfig,
		Handlers: api.NewHandlers(),
		DB:       postgresDB,
	}

	app.Run()
}
