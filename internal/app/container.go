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
	SessionService *service.SessionService
	TokensService  *service.TokensService
	UserService    *service.UserService
}

type AppRepositories struct {
	SessionRepository repository.SessionRepository
	UserRepository    repository.UserRepository
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
	app.repositories.SessionRepository = repository.NewSessionRepository(app.DB, app.Logger)
	app.repositories.UserRepository = repository.NewUserRepository(app.DB, app.Logger)
}

func (app *App) createServices() {
	app.Services.TokensService = service.NewTokensService(&app.Config.Keys, app.Logger)
	app.Services.SessionService = service.NewSessionService(
		app.repositories.SessionRepository,
		app.repositories.UserRepository,
		app.Services.TokensService,
		app.Logger,
	)
	app.Services.UserService = service.NewUserService(app.repositories.UserRepository, app.Logger)
}
