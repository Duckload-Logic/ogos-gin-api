package analytics

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/database"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, redis *database.RedisClient) {
	analyticsRoutes := rg.Group("/analytics")
	analyticsRoutes.Use(middleware.AuditContextMiddleware())
	analyticsRoutes.Use(middleware.AuthMiddleware(redis))

	analyticsRoutes.GET("/dashboard",
		middleware.RoleMiddleware(int(constants.CounselorRoleID)),
		h.GetAnalyticsDashboard,
	)
}
