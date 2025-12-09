package http

import (
	"mytro-backend-auth/internal/app"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// registerHealth registers the health check endpoint to the router.
func registerHealth(router *gin.Engine, app *app.App) {
	router.GET("/health", func(c *gin.Context) {
		app.Logger.Info("Health check endpoint called")
		sqlDB, err := app.DB.DB()
		if err != nil {
			app.Logger.Error("Failed to get database connection", zap.Error(err))
			c.JSON(500, gin.H{"status": "error", "message": "Database connection error"})
		}
		if err := sqlDB.Ping(); err != nil {
			app.Logger.Error("Database ping failed", zap.Error(err))
			c.JSON(500, gin.H{"status": "error", "message": "Database not reachable"})
			return
		}
		app.Logger.Info("Health check successful")
		c.JSON(200, gin.H{"status": "ok"})
	})
}
