package notifications

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, handler *Handler) {
	notificationRoutes := r.Group("/notifications")
	{
		notificationRoutes.GET("/:userId", handler.GetUserNotifications)
	}
}