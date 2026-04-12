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
	// Base group for all activity logs
	activityGroup := rg.Group("/activity-meta")
	activityGroup.Use(middleware.AuthMiddleware(redis))

	// User-specific activity route (No role check, just auth)
	activityGroup.GET("/me", h.GetMyLogs)

	// Admin-only routes (Requires SuperAdmin role)
	adminOnly := activityGroup.Group("")
	adminOnly.Use(middleware.RoleMiddleware(int(constants.SuperAdminRoleID)))
	{
		adminOnly.GET("", h.GetLogs)
		adminOnly.GET("/audit", h.GetAuditLogs)
		adminOnly.GET("/system", h.GetSystemLogs)
		adminOnly.GET("/security", h.GetSecurityLogs)
		adminOnly.GET("/stats", h.GetLogStats)
		adminOnly.GET("/activity-stats", h.GetActivityStats)
	}
}
