package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
)

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := uuid.New().String()
		c.Set(audit.TraceIDKey, traceID)
		c.Writer.Header().Set("X-Trace-ID", traceID)
		c.Next()
	}
}
