package analytics

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	analyticsRoutes := rg.Group("/analytics")

	analyticsRoutes.Use(middleware.AuthMiddleware())

	analyticsRoutes.GET("/dashboard", 
		middleware.RoleMiddleware(int(constants.CounselorRoleID)), 
		h.GetAnalyticsDashboard,
	)
}