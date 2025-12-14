package http

import (
	"mytro-backend-auth/internal/app"
	"mytro-backend-auth/internal/transport/http/middleware"
	v1 "mytro-backend-auth/internal/transport/http/v1"

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

	baseGroup := router.Group("/api/auth")

	// Register API v1 routes
	v1Group := baseGroup.Group("/v1")
	v1.RegisterRoutes(v1Group, app)

	// Register base API routes
	v1.RegisterRoutes(baseGroup, app)

	return router
}
