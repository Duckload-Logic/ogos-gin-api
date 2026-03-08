package logs

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	routes := rg.Group("/system-logs")
	routes.Use(middleware.AuthMiddleware())
	routes.Use(middleware.RoleMiddleware(
		int(constants.SuperAdminRoleID),
	))
	{
		routes.GET("", h.HandleListSystemLogs)
		routes.GET("/audit", h.HandleGetAuditLogs)
		routes.GET("/system", h.HandleGetSystemLogs)
		routes.GET("/security", h.HandleGetSecurityLogs)
		routes.GET("/stats", h.HandleGetLogStats)
	}
}
