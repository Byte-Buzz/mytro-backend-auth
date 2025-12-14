package handler

import (
	"mytro-backend-auth/internal/app"
	"mytro-backend-auth/internal/domain/models"
	"mytro-backend-auth/internal/transport/http/v1/request"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SignUpHandler(app *app.App) gin.HandlerFunc {
	userService := app.Services.UserService

	return func(c *gin.Context) {
		var signUpRequest request.SignUpRequest
		if err := c.ShouldBindJSON(&signUpRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user := toUserModel(signUpRequest)

		err := userService.SignUpByEmail(c.Request.Context(), &user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
	}
}

func toUserModel(req request.SignUpRequest) models.User {
	return models.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}
}
