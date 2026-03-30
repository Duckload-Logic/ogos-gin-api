package auth

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

func RegisterDebugRoutes(
	rg *gin.RouterGroup,
	redis *datastore.RedisClient,
) {
	debugGroup := rg.Group("/debug")
	{
		debugGroup.GET("/redis", func(c *gin.Context) {
			// Basic security check via secret token in URL
			secret := c.Query("secret")
			expectedSecret := os.Getenv("REDIS_DEBUG_SECRET")

			// If environment secret isn't set, default to a fallback for the student so they aren't locked out
			if expectedSecret == "" {
				expectedSecret = "DLL-Admin-2027"
			}

			if secret != expectedSecret {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized debug access"})
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Get all keys
			keys, err := redis.Client.Keys(ctx, "*").Result()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Get details for a few keys (optional, but nice)
			details := make(map[string]interface{})
			for _, key := range keys {
				val, _ := redis.Client.Get(ctx, key).Result()
				details[key] = val
			}

			c.JSON(http.StatusOK, gin.H{
				"totalKeys": len(keys),
				"keys":      keys,
				"data":      details,
				"server":    "Production Redis (7-alpine)",
				"timestamp": time.Now().Format(time.RFC3339),
			})
		})
	}
}
