package analytics

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAnalyticsDashboard(c *gin.Context) {
	ctx := c.Request.Context()

	dashboardData, err := h.service.GetDashboard(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate analytics dashboard",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Analytics retrieved successfully",
		"data":    dashboardData,
	})
}