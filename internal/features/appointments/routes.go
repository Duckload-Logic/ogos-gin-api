package appointments

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	appointmentRoutes := rg.Group("/appointments")
	{
		appointmentRoutes.POST("", h.Create)
	}
}
