package api

import "gorm.io/gorm"

type Config struct {
	Address string
}

type Application struct {
	Config   Config
	Handlers *Handlers
	DB       *gorm.DB
}
