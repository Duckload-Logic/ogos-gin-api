package notifications

import (
	"context"

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

// ServiceInterface defines the business logic for managing notifications.
type ServiceInterface interface {
	Send(
		ctx context.Context,
		notif audit.NotificationEntry,
	) error
	GetUserNotifications(
		ctx context.Context,
		userID string,
	) ([]audit.NotificationEntry, error)
	Subscribe(
		ctx context.Context,
		userID string,
	) (<-chan audit.NotificationEntry, func())
	Unsubscribe(
		ctx context.Context,
		userID string,
		ch chan audit.NotificationEntry,
	)
	Broadcast(
		ctx context.Context,
		notif audit.NotificationEntry,
	) error
	MarkAsRead(ctx context.Context, id string, userID string) error
	DeleteOldNotifications(ctx context.Context, days int) (int64, error)
}

// RepositoryInterface defines the data access layer for managing notifications.
type RepositoryInterface interface {
	WithTransaction(ctx context.Context, fn func(datastore.DB) error) error
	GetByUserID(ctx context.Context, userID string) ([]Notification, error)
	MarkAsRead(ctx context.Context, tx datastore.DB, id string, userID string) error
	Create(
		ctx context.Context,
		tx datastore.DB,
		notif Notification,
	) error
	DeleteOldNotifications(ctx context.Context, days int) (int64, error)
}
