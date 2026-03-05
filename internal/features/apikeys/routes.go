package apikeys

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

// RegisterRoutes sets up management endpoints for API keys (counselor-only).
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	routes := rg.Group("/api-keys")
	routes.Use(middleware.AuthMiddleware())
	routes.Use(middleware.RoleMiddleware(
		int(constants.CounselorRoleID),
	))
	{
		routes.POST("", h.HandleCreateAPIKey)
		routes.GET("", h.HandleListAPIKeys)
		routes.DELETE("/:id", h.HandleRevokeAPIKey)
	}
}
