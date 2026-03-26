package notifications

import (
	"context"
)

// ServiceInterface defines the business logic for managing notifications.
type ServiceInterface interface {
	Send(ctx context.Context, userID string, title, message, notifType string) error
	GetUserNotifications(ctx context.Context, userID string) ([]NotificationDTO, error)
	MarkAsRead(ctx context.Context, id int) error
}

// RepositoryInterface defines the data access layer for managing notifications.
type RepositoryInterface interface {
	GetByUserID(ctx context.Context, userID string) ([]NotificationModel, error)
	MarkAsRead(ctx context.Context, id int) error
	Create(ctx context.Context, userID string, title, message, notifType string) error
}
