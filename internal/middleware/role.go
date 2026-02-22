package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware
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
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{"error": "Access denied"},
			)
			return
		}

		c.Next()
	}
}
