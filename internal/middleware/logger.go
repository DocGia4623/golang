package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	// code for logging
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		log.Println("Logger middleware called") // Add this line
		return fmt.Sprintf("%s - [%s] \"%s %s\" %d %s\n",
			param.ClientIP,
			param.TimeStamp.UTC().Format(time.RFC822),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
		)
	})
}
