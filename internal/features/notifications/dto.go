package notifications

import "time"

type NotificationDTO struct {
    ID        uint      `json:"id"`
    UserID    string    `json:"userId"` 
    Title     string    `json:"title"`
    Message   string    `json:"message"`
    Type      string    `json:"type"` 
    IsRead    bool      `json:"isRead"`
    CreatedAt time.Time `json:"createdAt"`
}

type NotificationResponse struct {
    Success bool              `json:"success"`
    Message string            `json:"message"`
    Data    []NotificationDTO `json:"data"`
}

type ListNotificationsResponse struct {
    Notifications []NotificationDTO `json:"notifications"`
    Total         int               `json:"total"`
    Page          int               `json:"page"`
    PageSize      int               `json:"pageSize"`
    TotalPages    int               `json:"totalPages"`
}