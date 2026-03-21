package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/tokens"
	"github.com/olazo-johnalbert/duckload-api/internal/database"
)

func AuthMiddleware(redis *database.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// ... (cookie and header logic)
		cookie, err := c.Cookie("access_token")
		if err == nil && cookie != "" {
			tokenString = cookie
		}

		if tokenString == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No authentication token provided"})
			return
		}

		// Validate JWT
		claims, err := tokens.NewService().ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Validate against Redis
		if redis != nil {
			tokenKey := "token:" + tokenString
			_, err := redis.Get(c.Request.Context(), tokenKey)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked or expired"})
				return
			}
		}

		// Set user info in context
		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.UserEmail)
		c.Set("roleID", claims.RoleID)
		c.Set("tokenType", claims.TokenType)

		c.Next()
	}
}
