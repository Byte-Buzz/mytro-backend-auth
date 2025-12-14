package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GlobalMiddleware applies global middleware to the Gin router.
func Logger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info("Incoming request",
			zap.String("request_id", c.GetString("request_id")),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("ip", c.ClientIP()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
		)
		c.Next()
	}
}
