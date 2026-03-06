package notifications

import "time"

type NotificationDTO struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"` // e.g., "Appointment", "System"
	IsRead    bool      `json:"isRead"`
	CreatedAt time.Time `json:"createdAt"`
}

type NotificationResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    []NotificationDTO `json:"data"`
}