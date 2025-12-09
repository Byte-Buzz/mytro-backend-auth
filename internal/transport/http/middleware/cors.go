package middleware

import (
	"mytro-backend-auth/internal/infrastructure/config"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS applies CORS middleware to the Gin router.
func CORS(cors config.CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Determine if the origin is allowed
		var allowedOrigin string
		if contains(cors.AllowedOrigins, "*") {
			allowedOrigin = "*"
		} else if origin != "" && contains(cors.AllowedOrigins, origin) {
			allowedOrigin = origin
		}

		if allowedOrigin == "" && len(cors.AllowedOrigins) > 0 {
			// Origin not allowed
			c.AbortWithStatus(403)
			return
		}

		// Set CORS headers
		if allowedOrigin != "" {
			c.Header("Access-Control-Allow-Origin", allowedOrigin)
		}
		c.Header("Access-Control-Allow-Methods", strings.Join(cors.AllowedMethods, ", "))
		c.Header("Access-Control-Allow-Headers", strings.Join(cors.AllowedHeaders, ", "))
		if cors.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // No Content
			return
		}

		c.Next()
	}
}

// contains checks if a slice contains a specific string.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
