package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerMiddleware trả về một middleware để logging request
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		logger.Info("Request handled",
			zap.String("time", start.UTC().Format(time.RFC3339)),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.String("latency", latency.String()),
			zap.String("service", "golang-app"),
			zap.String("clientIP", c.ClientIP()),
		)
	}
}
