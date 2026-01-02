package appointments

import (
    "github.com/gin-gonic/gin"
    "github.com/olazo-johnalbert/duckload-api/internal/core/constants"
    "github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
    appointmentRoutes := rg.Group("/appointments")
    appointmentRoutes.PATCH("/:id/status", h.HandleUpdateStatus)
    appointmentRoutes.Use(middleware.RoleMiddleware(
        int(constants.StudentRoleID),
        int(constants.CounselorRoleID),
        int(constants.FrontDeskRoleID),
    ))
    {
        appointmentRoutes.POST("", h.Create)
        appointmentRoutes.GET("", h.HandleListAppointments)
        appointmentRoutes.GET("/:id", h.HandleGetAppointment)
        appointmentRoutes.GET("/student/:studentID", h.HandleGetStudentAppointments)
    }
}