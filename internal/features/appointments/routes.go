package appointments

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	routes := rg.Group("/appointments")
	routes.Use(middleware.AuthMiddleware())

	adminOnly := routes.Group("")
	adminOnly.Use(middleware.RoleMiddleware(
		int(constants.CounselorRoleID),
	))
	{
		adminOnly.GET("", h.HandleListAppointments)
		adminOnly.PATCH("/:id/status", h.HandleUpdateAppointmentStatus)
		adminOnly.GET("/calendar/stats", h.HandleGetDailyStatusCount)
	}

	studentOnly := routes.Group("")
	studentOnly.Use(middleware.RoleMiddleware(
		int(constants.StudentRoleID),
	))
	{
		studentOnly.GET("/me", h.HandleGetAppointmentsByUserID)
	}

	sharedRoutes := routes.Group("")
	sharedRoutes.Use(middleware.RoleMiddleware(
		int(constants.StudentRoleID),
		int(constants.CounselorRoleID),
	))
	{
		sharedRoutes.GET("/id/:id", h.GetAppointmentByID)
		sharedRoutes.GET("/stats", h.HandleGetAppointmentStats)
		sharedRoutes.GET("/lookups/categories", h.HandleGetConcernCategories)
		sharedRoutes.GET("/lookups/slots", h.HandleGetAvailableTimeSlots)
		sharedRoutes.GET("/lookups/statuses", h.HandleGetAppointmentStatuses)

		sharedRoutes.POST("", h.HandleCreateAppointment)
		sharedRoutes.PATCH("/id/:id", h.HandleUpdateAppointment)
	}
}
