package app

import (
	"mytro-backend-auth/internal/infrastructure/config"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	Config *config.Config
	Logger *zap.Logger
	DB     *gorm.DB

	services     *AppServices
	repositories *AppRepositories
}

type AppServices struct {
}

type AppRepositories struct {
}

func NewApp(config *config.Config, db *gorm.DB, logger *zap.Logger) *App {
	return &App{
		Config: config,
		Logger: logger,
		DB:     db,

		services:     &AppServices{},
		repositories: &AppRepositories{},
	}
}
