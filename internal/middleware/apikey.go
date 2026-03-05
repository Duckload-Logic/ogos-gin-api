package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIKeyValidator is a function that validates a plaintext API key
// and returns the key ID, key name, and an error if invalid.
type APIKeyValidator func(ctx context.Context, plaintext string) (id int, name string, err error)

// APIKeyMiddleware validates the X-API-Key header using the provided validator.
// On success it sets "apiKeyID" and "apiKeyName" in the gin context.
func APIKeyMiddleware(validate APIKeyValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")

		if key == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "API key is required (X-API-Key header)"},
			)
			return
		}

		id, name, err := validate(c.Request.Context(), key)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": err.Error()},
			)
			return
		}

		c.Set("apiKeyID", id)
		c.Set("apiKeyName", name)
		c.Next()
	}
}
