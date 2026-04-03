package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// APIKeyValidator is a function that validates a plaintext API key
// and returns the key ID, key name, and an error if invalid.
type APIKeyValidator func(
	ctx context.Context, plaintext string,
) (id int, name string, err error)

// APIKeyMiddleware validates the X-API-Key header using the provided validator.
// On success it sets "apiKeyID" and "apiKeyName" in the gin context.
func APIKeyMiddleware(validate APIKeyValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var logSvc SecurityLogger
		if svc, ok := c.Get(SecurityLoggerContextKey); ok {
			logSvc, _ = svc.(SecurityLogger)
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "API key is required (Authorization header)"},
			)
			return
		}

		// Extract the API key from the Authorization header
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "Invalid Authorization header format"},
			)
			return
		}

		key := strings.TrimPrefix(authHeader, "Bearer ")
		if key == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "API key cannot be empty"},
			)
			return
		}

		id, name, err := validate(c.Request.Context(), key)
		if err != nil {
			if logSvc != nil {
				prefix := key
				if len(prefix) > 8 {
					prefix = prefix[:8]
				}
				logSvc.RecordSecurity(
					c.Request.Context(),
					"API_KEY_INVALID",
					fmt.Sprintf(
						"Invalid API key attempt (prefix: %s...): %s",
						prefix,
						err.Error(),
					),
					"",
					c.ClientIP(),
					c.Request.UserAgent(),
				)
			}
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": err.Error()},
			)
			return
		}

		if logSvc != nil {
			logSvc.RecordSecurity(
				c.Request.Context(),
				"API_KEY_USED",
				fmt.Sprintf(
					"API key '%s' (ID: %d) used on %s %s",
					name,
					id,
					c.Request.Method,
					c.Request.URL.Path,
				),
				"",
				c.ClientIP(),
				c.Request.UserAgent(),
			)
		}

		c.Set("apiKeyID", id)
		c.Set("apiKeyName", name)
		c.Next()
	}
}
