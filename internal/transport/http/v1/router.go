package v1

import (
	"mytro-backend-auth/internal/app"
	"mytro-backend-auth/internal/transport/http/v1/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, app *app.App) {
	router.POST("/sign-up", handler.SignUpHandler(app))
}
