package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/tokens"
)

var tokenService = tokens.NewService()

// AuthMiddleware
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
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "Invalid or expired token"},
			)
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("roleID", claims.RoleID)
		c.Next()
	}
}
