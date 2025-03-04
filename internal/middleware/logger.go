package middleware

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger middleware để format log theo JSON
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		logEntry := map[string]interface{}{
			"time":     start.UTC().Format(time.RFC3339),
			"level":    "INFO",
			"message":  "Request handled",
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"status":   c.Writer.Status(),
			"latency":  latency.String(),
			"service":  "golang-app",
			"clientIP": c.ClientIP(),
		}

		// Convert log thành JSON
		logData, err := json.Marshal(logEntry)
		if err != nil {
			log.Printf("❌ Lỗi JSON: %v", err)
		} else {
			jsonLog := string(logData)
			jsonLog = strings.ReplaceAll(jsonLog, "\n", " ")
			jsonLog = strings.ReplaceAll(jsonLog, "\r", " ")
			log.Print(jsonLog)
		}
	}
}
