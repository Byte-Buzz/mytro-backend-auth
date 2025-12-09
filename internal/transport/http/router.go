package http

import (
	"mytro-backend-auth/internal/app"
	"mytro-backend-auth/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(app *app.App) *gin.Engine {
	// Create a new Gin router
	router := gin.Default()

	// Apply middlewares
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger(app.Logger))
	router.Use(middleware.CORS(app.Config.CORS))
	router.Use(middleware.SecurityHeaders())
	router.Use(gin.Recovery())

	// Register health check endpoint
	registerHealth(router, app)

	return router
}
