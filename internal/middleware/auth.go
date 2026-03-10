package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/tokens"
)

var tokenService = tokens.NewService()

// AuthMiddleware validates JWT tokens and sets user context.
// Reads the log service from gin context ("logService") to record INVALID_TOKEN events.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "Authorization header is required"},
			)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := tokenService.ValidateToken(tokenString)
		if err != nil {
			if logSvc, ok := c.Get(SecurityLoggerContextKey); ok {
				if svc, ok := logSvc.(SecurityLogger); ok {
					svc.RecordSecurity(c.Request.Context(), "INVALID_TOKEN",
						fmt.Sprintf("Invalid or expired token used on %s %s", c.Request.Method, c.Request.URL.Path),
						"", c.ClientIP(), c.Request.UserAgent())
				}
			}
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "Invalid or expired token"},
			)
			return
		}

		newCtx := audit.WithContext(
			c.Request.Context(),
			c.ClientIP(),
			c.Request.UserAgent(),
			claims.UserID,
			claims.UserEmail, // Assuming Email exists in your claims
		)

		c.Request = c.Request.WithContext(newCtx)

		c.Set("userEmail", claims.UserEmail) // Set user email in context for logging
		c.Set("userID", claims.UserID)
		c.Set("roleID", claims.RoleID)
		c.Next()
	}
}
