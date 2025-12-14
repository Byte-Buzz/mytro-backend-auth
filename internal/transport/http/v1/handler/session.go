package handler

import (
	"mytro-backend-auth/internal/app"
	"mytro-backend-auth/internal/transport/http/v1/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SignOutHandler(app *app.App) gin.HandlerFunc {
	sessionService := app.Services.SessionService

	return func(c *gin.Context) {
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
			return
		}

		err = sessionService.SignOut(c.Request.Context(), refreshToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Clear the refresh token cookie
		c.SetCookie("refresh_token", "", -1, "/", "", false, true)

		c.JSON(http.StatusOK, gin.H{"message": "Signed out successfully"})
	}
}

func UpdateSessionHandler(app *app.App) gin.HandlerFunc {
	sessionService := app.Services.SessionService

	return func(c *gin.Context) {
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
			return
		}

		tokens, err := sessionService.RefreshSession(c.Request.Context(), refreshToken, c.ClientIP(), c.Request.UserAgent())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		maxAge := tokens.RefreshTokenExpiresAt.Unix() - time.Now().Unix()

		c.SetCookie("refresh_token", tokens.RefreshToken, int(maxAge), "/", "", false, true)

		c.JSON(http.StatusOK, response.TokenResponse{
			AccessToken: tokens.AccessToken,
		})
	}
}
