package handler

import (
	"mytro-backend-auth/internal/app"
	"mytro-backend-auth/internal/transport/http/v1/request"
	"mytro-backend-auth/internal/transport/http/v1/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SignInByEmailHandler(app *app.App) gin.HandlerFunc {
	sessionService := app.Services.SessionService

	return func(c *gin.Context) {
		var signInRequest request.SignInRequest
		if err := c.ShouldBindJSON(&signInRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tokens, err := sessionService.SignInByEmail(c.Request.Context(), signInRequest.UsernameOrEmail, signInRequest.Password, c.ClientIP(), c.Request.UserAgent())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		maxAge := tokens.RefreshTokenExpiresAt.Unix() - time.Now().Unix()

		// ToDO: Set refresh token domain
		c.SetCookie("refresh_token", tokens.RefreshToken, int(maxAge), "/", "", false, true)

		c.JSON(http.StatusOK, response.TokenResponse{
			AccessToken: tokens.AccessToken,
		})
	}
}
