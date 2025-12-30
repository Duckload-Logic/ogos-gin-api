package excuseslips

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup, h *Handler) {
    api.POST("/excuseslips", h.Submit)
    api.GET("/excuseslips", h.GetAll)
    api.GET("/excuseslips/:id", h.GetByID)
	api.PATCH("/excuseslips/:id/status", h.UpdateStatus)
	api.DELETE("/excuseslips/:id", h.Delete)
}