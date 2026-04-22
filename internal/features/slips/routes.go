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
	routes.Use(middleware.HydrateStudentIIRContext(db))
	routes.Use(middleware.HydrateStudentCORContext(db))
	routes.Use(middleware.AuditContextMiddleware())

	adminOnly := routes.Group("")
	adminOnly.Use(middleware.RoleMiddleware(
		int(constants.AdminRoleID),
	))
	{
		adminOnly.GET("", h.GetSlips)
		adminOnly.GET("/urgent", h.GetSlipUrgent)
		adminOnly.PATCH("/id/:id/status", h.PatchSlipStatus)
	}

	studentOnly := routes.Group("")
	studentOnly.Use(middleware.RoleMiddleware(
		int(constants.StudentRoleID),
	))
	{
		studentOnly.GET("/me", h.GetSlipMe)
		studentOnly.POST("", h.PostSlip)
		studentOnly.PATCH("/id/:id", h.PatchSlip)
	}

	sharedRoutes := routes.Group("")
	sharedRoutes.Use(middleware.RoleMiddleware(
		int(constants.AdminRoleID),
		int(constants.StudentRoleID),
	))
	{
		sharedRoutes.GET("/id/:id", h.GetSlipByID)
		sharedRoutes.GET("/stats", h.GetSlipStats)
		sharedRoutes.GET(
			"/id/:id/attachments",
			h.GetSlipAttachments,
		)
		sharedRoutes.GET(
			"/id/:id/attachments/:attachmentId",
			h.GetSlipAttachmentContent,
		)
		sharedRoutes.GET(
			"/lookups/statuses",
			h.GetSlipStatuses,
		)
		sharedRoutes.GET(
			"/lookups/categories",
			h.GetSlipCategories,
		)
	}
}
