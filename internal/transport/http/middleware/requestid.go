package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID is a middleware that assigns a unique request ID to each incoming HTTP request.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := uuid.New().String()
		c.Header("X-Request-ID", uuid)
		c.Set("request_id", uuid)
		c.Next()
	}
}
