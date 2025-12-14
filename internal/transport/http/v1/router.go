package v1

import (
	"mytro-backend-auth/internal/app"
	"mytro-backend-auth/internal/transport/http/v1/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, app *app.App) {
	emailGroup := router.Group("/email")
	{
		emailGroup.POST("/sign-up", handler.SignUpByEmailHandler(app))
		emailGroup.POST("/sign-in", handler.SignInByEmailHandler(app))
	}
	router.POST("/sign-out", handler.SignOutHandler(app))
	router.POST("/refresh", handler.UpdateSessionHandler(app))
}
