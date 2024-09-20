package main

import (
	"gopher-social-backend-server/cmd/server/api"
	"os"
)

var DEBUG = os.Getenv("DEBUG") == "true"
var PORT = os.Getenv("PORT")

func main() {
	appConfig := api.Config{
		Address: ":" + PORT,
	}

	app := &api.Application{
		Config: appConfig,
	}

	app.Run()
}
