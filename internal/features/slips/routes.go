package slips

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

func RegisterRoutes(db *sqlx.DB, api *gin.RouterGroup, h *Handler) {
	excuseslipGroup := api.Group("/slips")
	excuseslipGroup.Use(middleware.AuthMiddleware())
	excuseslipGroup.Use(middleware.HydrateStudentContext(db))
	excuseslipGroup.Use(middleware.AuditContextMiddleware())

	adminOnly := excuseslipGroup.Group("")
	adminOnly.Use(middleware.RoleMiddleware(
		int(constants.CounselorRoleID),
	))
	{
		adminOnly.GET("", h.GetSlipList)
		adminOnly.GET("/urgent", h.GetUrgentSlipList)
		adminOnly.PATCH("/id/:id/status", h.PatchSlipStatus)
	}

	studentOnly := excuseslipGroup.Group("")
	studentOnly.Use(middleware.RoleMiddleware(
		int(constants.StudentRoleID),
	))
	{
		studentOnly.GET("/me", h.GetSlipListByIIR)
		studentOnly.POST("", h.PostSlip)
	}

	sharedRoutes := excuseslipGroup.Group("")
	sharedRoutes.Use(middleware.RoleMiddleware(
		int(constants.CounselorRoleID),
		int(constants.StudentRoleID),
	))
	{
		sharedRoutes.GET("/stats", h.GetSlipStatsList)
		sharedRoutes.GET(
			"/id/:id/attachments",
			h.GetSlipAttachmentList,
		)
		sharedRoutes.GET(
			"/id/:id/attachments/:attachmentId",
			h.GetAttachmentFile,
		)
		sharedRoutes.GET(
			"/lookups/statuses",
			h.GetSlipStatusList,
		)
		sharedRoutes.GET(
			"/lookups/categories",
			h.GetSlipCategoryList,
		)
	}
}
