package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
)

// AuditContextMiddleware enriches the request context with audit metadata
// (user email, IP address, User-Agent). Must be placed after AuthMiddleware
// so that userEmail is available.
func AuditContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userEmail := ""
		if email, exists := c.Get("userEmail"); exists {
			if e, ok := email.(string); ok {
				userEmail = e
			}
		}

		ipAddress := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")

		ctx := audit.WithContext(c.Request.Context(), userEmail, ipAddress, userAgent)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
