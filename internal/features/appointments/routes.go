package appointments

import (
    "github.com/gin-gonic/gin"
    "github.com/olazo-johnalbert/duckload-api/internal/core/constants"
    "github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
    appointmentRoutes := rg.Group("/appointments")
    
    appointmentRoutes.Use(middleware.AuthMiddleware())

    appointmentRoutes.POST("", middleware.RoleMiddleware(
        int(constants.StudentRoleID),
        int(constants.FrontDeskRoleID), 
    ), h.Create)

    appointmentRoutes.GET("", middleware.RoleMiddleware(
        int(constants.StudentRoleID), 
        int(constants.CounselorRoleID), 
        int(constants.FrontDeskRoleID),
    ), h.HandleListAppointments)

    appointmentRoutes.GET("/:id", middleware.RoleMiddleware(
        int(constants.StudentRoleID), 
        int(constants.CounselorRoleID), 
        int(constants.FrontDeskRoleID),
    ), h.HandleGetAppointment)

    appointmentRoutes.GET("/student/:studentID", middleware.RoleMiddleware(
        int(constants.StudentRoleID), 
        int(constants.CounselorRoleID), 
        int(constants.FrontDeskRoleID),
    ), h.HandleGetStudentAppointments)

    appointmentRoutes.PATCH("/:id/status", middleware.RoleMiddleware(
        int(constants.CounselorRoleID),
        int(constants.FrontDeskRoleID),
    ), h.HandleUpdateStatus)
}