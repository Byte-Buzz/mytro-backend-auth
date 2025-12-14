package app

import (
	"mytro-backend-auth/internal/app/repository"
	"mytro-backend-auth/internal/app/service"
	"mytro-backend-auth/internal/infrastructure/config"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	Config *config.Config
	Logger *zap.Logger
	DB     *gorm.DB

	Services     *AppServices
	repositories *AppRepositories
}

type AppServices struct {
	UserService *service.UserService
}

type AppRepositories struct {
	UserRepository repository.UserRepository
}

func NewApp(config *config.Config, db *gorm.DB, logger *zap.Logger) *App {
	app := &App{
		Config: config,
		Logger: logger,
		DB:     db,

		Services:     &AppServices{},
		repositories: &AppRepositories{},
	}

	app.createRepositories()
	app.createServices()

	return app
}

func (app *App) createRepositories() {
	app.repositories.UserRepository = repository.NewUserRepository(app.DB, app.Logger)
}

func (app *App) createServices() {
	app.Services.UserService = service.NewUserService(app.repositories.UserRepository, app.Logger)
}
