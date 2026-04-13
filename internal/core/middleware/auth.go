package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/sessions"
	"github.com/olazo-johnalbert/duckload-api/internal/core/tokens"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

func AuthMiddleware(redis *datastore.RedisClient) gin.HandlerFunc {
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
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "No authentication token provided"},
			)
			return
		}

		// Validate JWT signature and expiration first
		claims, err := tokens.NewService().ValidateToken(tokenString)
		if err != nil {
			log.Printf(
				"[AuthMiddleware] {Token}: Invalid or expired token: %v",
				err,
			)
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "Invalid or expired token"},
			)
			return
		}

		// Validate against Redis using the Token ID (jti)
		if redis != nil {
			jti := sessions.NewJTI(claims.ID)
			tokenKey := jti.ToSessionKey()
			val, err := redis.Get(c.Request.Context(), tokenKey)
			if err != nil {
				log.Printf(
					"[AuthMiddleware] {Redis}: Session %s missing or "+
						"revoked for ID: %s",
					claims.ID,
					func() string {
						if claims.M2MClientID != "" {
							return "M2M:" + claims.M2MClientID
						}
						return "User:" + claims.UserID
					}(),
				)
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					gin.H{"error": "Session has been revoked or expired"},
				)
				return
			}

			// Parse session data
			var sessionData map[string]string
			if err := json.Unmarshal([]byte(val), &sessionData); err == nil {
				if idpToken, ok := sessionData["idpAccessToken"]; ok {
					c.Set("idpAccessToken", idpToken)
				}
			}
		}

		// Set info in context
		if claims.M2MClientID != "" {
			c.Set("m2mClientID", claims.M2MClientID)
			c.Set("isM2M", true)
		} else {
			c.Set("userID", claims.UserID)
			c.Set("userEmail", claims.UserEmail)
			c.Set("roleID", claims.RoleID)
		}
		c.Set("tokenType", claims.TokenType)

		c.Next()
	}
}
