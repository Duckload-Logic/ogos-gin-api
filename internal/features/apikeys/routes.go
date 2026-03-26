package apikeys

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/middleware"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

// RegisterRoutes sets up management endpoints for API keys (counselor-only).
func RegisterRoutes(
	rg *gin.RouterGroup,
	h *Handler,
	redis *datastore.RedisClient,
) {
	apiKeysRoutes := rg.Group("/access-tokens")
	apiKeysRoutes.Use(middleware.AuthMiddleware(redis))
	apiKeysRoutes.Use(middleware.RoleMiddleware(
		int(constants.SuperAdminRoleID),
	))
	{
		apiKeysRoutes.GET("", h.GetAPIKeys)
		apiKeysRoutes.POST("", h.PostAPIKey)
		apiKeysRoutes.DELETE("/:id", h.DeleteAPIKey)
	}
}
