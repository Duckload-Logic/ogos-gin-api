package slips

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/middleware"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

func RegisterRoutes(
	db *sqlx.DB,
	rg *gin.RouterGroup,
	h *Handler,
	redis *datastore.RedisClient,
) {
	routes := rg.Group("/slips")
	routes.Use(middleware.AuthMiddleware(redis))
	routes.Use(middleware.HydrateStudentContext(db))
	routes.Use(middleware.AuditContextMiddleware())

	adminOnly := routes.Group("")
	adminOnly.Use(middleware.RoleMiddleware(
		int(constants.CounselorRoleID),
	))
	{
		adminOnly.GET("", h.GetSlipList)
		adminOnly.GET("/urgent", h.GetUrgentSlipList)
		adminOnly.PATCH("/id/:id/status", h.PatchSlipStatus)
	}

	studentOnly := routes.Group("")
	studentOnly.Use(middleware.RoleMiddleware(
		int(constants.StudentRoleID),
	))
	{
		studentOnly.GET("/me", h.GetSlipListByIIR)
		studentOnly.POST("", h.PostSlip)
	}

	sharedRoutes := routes.Group("")
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
