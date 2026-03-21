package appointments

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/database"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

func RegisterRoutes(db *sqlx.DB, rg *gin.RouterGroup, h *Handler, redis *database.RedisClient) {
	routes := rg.Group("/appointments")
	routes.Use(middleware.AuthMiddleware(redis))
	routes.Use(middleware.HydrateStudentContext(db))
	routes.Use(middleware.AuditContextMiddleware())

	adminOnly := routes.Group("")
	adminOnly.Use(middleware.RoleMiddleware(
		int(constants.CounselorRoleID),
	))
	{
		adminOnly.GET("", h.GetAppointmentList)
		adminOnly.GET(
			"/calendar/stats",
			h.GetDailyStatusCountList,
		)
	}

	studentOnly := routes.Group("")
	studentOnly.Use(middleware.RoleMiddleware(
		int(constants.StudentRoleID),
	))
	{
		studentOnly.GET("/me", h.GetAppointmentListByIIR)
		studentOnly.POST("", h.PostAppointment)
	}

	sharedRoutes := routes.Group("")
	sharedRoutes.Use(middleware.RoleMiddleware(
		int(constants.StudentRoleID),
		int(constants.CounselorRoleID),
	))
	{
		sharedRoutes.GET("/id/:id", h.GetAppointmentByID)
		sharedRoutes.GET("/stats", h.GetAppointmentStatsList)
		sharedRoutes.GET(
			"/lookups/categories",
			h.GetAppointmentCategoryList,
		)
		sharedRoutes.GET(
			"/lookups/slots",
			h.GetAvailableTimeSlotList,
		)
		sharedRoutes.GET(
			"/lookups/statuses",
			h.GetAppointmentStatusList,
		)
		sharedRoutes.PATCH(
			"/id/:id",
			h.PatchAppointment,
		)
	}
}
