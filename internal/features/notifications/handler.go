package notifications

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/response"
)

type Handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetNotifications(c *gin.Context) {
	userID := c.Param("userId")

	notifications, err := h.service.GetUserNotifications(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		response.SendError(c, "Failed to fetch notifications", http.StatusInternalServerError, nil)
		return
	}

	data := ListNotificationsResponse{
		Notifications: notifications,
		Total:         len(notifications),
		Page:          1,
		PageSize:      len(notifications),
		TotalPages:    1,
	}

	response.SendSuccess(c, data)
}

func (h *Handler) PatchNotificationRead(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.SendFail(c, gin.H{"error": "Invalid notification ID"})
		return
	}

	if err := h.service.MarkAsRead(c.Request.Context(), id); err != nil {
		response.SendError(c, "Failed to mark notification as read", http.StatusInternalServerError, nil)
		return
	}

	response.SendSuccess(c, gin.H{"message": "Notification marked as read"})
}
