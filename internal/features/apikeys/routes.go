package apikeys

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/database"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

// RegisterRoutes sets up management endpoints for API keys (counselor-only).
func RegisterRoutes(rg *gin.RouterGroup, h *Handler, redis *database.RedisClient) {
	apiKeysRoutes := rg.Group("/access-tokens")
	apiKeysRoutes.Use(middleware.AuthMiddleware(redis))
	apiKeysRoutes.Use(middleware.RoleMiddleware(
		int(constants.SuperAdminRoleID),
	))
	{
		apiKeysRoutes.GET("", h.HandleListAPIKeys)
		apiKeysRoutes.POST("", h.HandleCreateAPIKey)
		apiKeysRoutes.DELETE("/:id", h.HandleRevokeAPIKey)
	}
}
