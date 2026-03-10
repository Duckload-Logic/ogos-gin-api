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
		userID := 0
		if id, exists := c.Get("userID"); exists {
			if e, ok := id.(int); ok {
				userID = e
			}
		}

		userEmail, exists := c.Get("userEmail")
		if !exists {
			userEmail = ""
		} else if e, ok := userEmail.(string); ok {
			userEmail = e
		}

		ipAddress := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")

		ctx := audit.WithContext(c.Request.Context(), ipAddress, userAgent, userID, userEmail.(string))
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
