package slips

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

func RegisterRoutes(api *gin.RouterGroup, h *Handler) {
	excuseslipGroup := api.Group("/slips")
	excuseslipGroup.Use(middleware.AuthMiddleware())
	excuseslipGroup.Use(middleware.AuditContextMiddleware())

	adminOnly := excuseslipGroup.Group("")
	adminOnly.Use(middleware.RoleMiddleware(int(constants.CounselorRoleID)))
	{
		adminOnly.GET("", h.GetAll)
		adminOnly.GET("/urgent", h.HandleGetUrgentSlips)
		adminOnly.PATCH("/id/:id/status", h.UpdateStatus)
	}

	studentOnly := excuseslipGroup.Group("")
	studentOnly.Use(middleware.RoleMiddleware(int(constants.StudentRoleID)))
	{
		studentOnly.GET("/me", h.GetUserSlips)
		studentOnly.POST("", h.Submit)
	}

	sharedRoutes := excuseslipGroup.Group("")
	sharedRoutes.Use(middleware.RoleMiddleware(
		int(constants.CounselorRoleID),
		int(constants.StudentRoleID),
	))
	{
		sharedRoutes.GET("/stats", h.GetSlipStats)
		sharedRoutes.GET("/id/:id/attachments", h.GetSlipAttachments)
		sharedRoutes.GET("/id/:id/attachments/:attachmentId", h.HandleDownloadAttachment)
		sharedRoutes.GET("/lookups/statuses", h.HandleGetSlipStatuses)
		sharedRoutes.GET("/lookups/categories", h.HandleGetSlipCategories)
	}

	// excuseslipGroup.DELETE("/:id", middleware.RoleMiddleware(
	// 	int(constants.CounselorRoleID),
	// 	int(constants.FrontDeskRoleID),
	// ), h.Delete)
}
