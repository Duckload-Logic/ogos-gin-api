package appointments

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
    appointmentRoutes := rg.Group("/appointments")
    {
        appointmentRoutes.POST("", h.Create)
        appointmentRoutes.GET("", h.HandleListAppointments)                
        appointmentRoutes.GET("/:id", h.HandleGetAppointment)                
        appointmentRoutes.GET("/student/:studentID", h.HandleGetStudentAppointments)
    }
}