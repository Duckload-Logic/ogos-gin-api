package notifications

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service // Injected service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetUserNotifications(c *gin.Context) {
	userID := c.Param("userId")

	notifications, err := h.service.GetUserNotifications(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch notifications",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Notifications retrieved successfully",
		"data":    notifications,
	})
}