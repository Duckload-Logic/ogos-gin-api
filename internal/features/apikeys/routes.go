package apikeys

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/database"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

// RegisterRoutes sets up management endpoints for API keys (counselor-only).
func RegisterRoutes(rg *gin.RouterGroup, h *Handler, redis *database.RedisClient) {
	apiKeysRoutes := rg.Group("/api-keys")
	apiKeysRoutes.Use(middleware.AuthMiddleware(redis))
	apiKeysRoutes.Use(middleware.RoleMiddleware(
		int(constants.SuperAdminRoleID),
	))
	{
		apiKeysRoutes.DELETE("/:id", h.HandleRevokeAPIKey)
	}
}
