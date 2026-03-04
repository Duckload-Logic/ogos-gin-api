package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
)

// AuditContextMiddleware enriches the request context with audit metadata
// (user ID, IP address, User-Agent). Must be placed after AuthMiddleware
// so that userID is available.
func AuditContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := 0
		if id, exists := c.Get("userID"); exists {
			if uid, ok := id.(int); ok {
				userID = uid
			}
		}

		ipAddress := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")

		ctx := audit.WithContext(c.Request.Context(), userID, ipAddress, userAgent)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
