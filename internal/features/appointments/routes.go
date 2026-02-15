package appointments

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	
	routes := rg.Group("/appointments")

	routes.Use(middleware.AuthMiddleware())

	routes.GET("/slots", middleware.RoleMiddleware(
		int(constants.StudentRoleID),
		int(constants.CounselorRoleID),
		int(constants.FrontDeskRoleID),
	), h.HandleGetAvailableTimeSlots)

	// List All (Counselor)
	routes.GET("", middleware.RoleMiddleware(
		int(constants.CounselorRoleID),
	), h.HandleListAppointments) 

	// List the Appointments (Student)
	routes.GET("/me", middleware.RoleMiddleware(
		int(constants.StudentRoleID),
	), h.GetAppointments) 

	// Create New
	routes.POST("", middleware.RoleMiddleware(
		int(constants.StudentRoleID),
		int(constants.FrontDeskRoleID),
	), h.Create)

	// Get Details
	routes.GET("/:id", middleware.RoleMiddleware(
		int(constants.StudentRoleID),
		int(constants.CounselorRoleID),
		int(constants.FrontDeskRoleID),
	), h.HandleGetAppointment) 

	// Update Status or Reschedule
	routes.PATCH("/:id", middleware.RoleMiddleware(
		int(constants.CounselorRoleID),
		int(constants.FrontDeskRoleID),
	), h.HandleUpdateStatus)
}