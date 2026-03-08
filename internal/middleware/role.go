package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware checks if the user's role is in the allowed list.
// Optionally accepts a log service to record ACCESS_DENIED events.
func RoleMiddleware(allowedRoles ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID, exists := c.Get("roleID")
		if !exists {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{"error": "Role not found"},
			)
			return
		}

		rid, ok := roleID.(int)
		if !ok {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{"error": "Invalid role type"},
			)
			return
		}

		isAuthorized := false
		for _, role := range allowedRoles {
			if role == rid {
				isAuthorized = true
				break
			}
		}

		if !isAuthorized {
			// Log access denied if log service is available
			if logSvc, ok := c.Get(SecurityLoggerContextKey); ok {
				if svc, ok := logSvc.(SecurityLogger); ok {
					userEmail, _ := c.Get("userEmail")
					svc.RecordSecurity(c.Request.Context(), "ACCESS_DENIED",
						fmt.Sprintf("Access denied for %v (role %d) on %s %s", userEmail, rid, c.Request.Method, c.Request.URL.Path),
						fmt.Sprintf("%v", userEmail), c.ClientIP(), c.Request.UserAgent())
				}
			}
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{"error": "Access denied"},
			)
			return
		}

		c.Next()
	}
}
