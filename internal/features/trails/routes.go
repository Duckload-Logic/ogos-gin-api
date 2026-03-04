package trails

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	routes := rg.Group("/audit-trails")
	routes.Use(middleware.AuthMiddleware())
	routes.Use(middleware.RoleMiddleware(
		int(constants.CounselorRoleID),
	))
	{
		routes.GET("", h.HandleListAuditTrails)
		routes.GET("/:entityType/:entityId", h.HandleGetEntityAuditTrails)
	}
}
