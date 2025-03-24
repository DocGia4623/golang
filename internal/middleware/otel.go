package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

// TracingMiddleware adds OpenTelemetry tracing to Gin requests
func TracingMiddleware() gin.HandlerFunc {
	tracer := otel.Tracer("my-tracer") // Tracer should be initialized inside function
	return func(c *gin.Context) {
		ctx, span := tracer.Start(c.Request.Context(), c.FullPath())
		defer span.End()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
