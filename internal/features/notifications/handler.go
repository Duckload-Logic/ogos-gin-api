package notifications

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
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

	data := ListNotificationsResponse{
		Notifications: notifications,
		Total:         len(notifications),
		Page:          1,
		PageSize:      len(notifications),
		TotalPages:    1,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Notifications retrieved successfully",
		"data":    data,
	})
}

func (h *Handler) MarkAsRead(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid notification ID",
		})
		return
	}

	if err := h.service.MarkAsRead(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to mark notification as read",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Notification marked as read",
	})
}
