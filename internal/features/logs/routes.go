package logs

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/middleware"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

func RegisterRoutes(
	rg *gin.RouterGroup,
	h *Handler,
	redis *datastore.RedisClient,
) {
	routes := rg.Group("/activity-meta")
	routes.Use(middleware.AuthMiddleware(redis))
	routes.Use(middleware.RoleMiddleware(
		int(constants.SuperAdminRoleID),
	))
	{
		routes.GET("", h.GetLogs)
		routes.GET("/audit", h.GetAuditLogs)
		routes.GET("/system", h.GetSystemLogs)
		routes.GET("/security", h.GetSecurityLogs)
		routes.GET("/stats", h.GetLogStats)
	}
}
