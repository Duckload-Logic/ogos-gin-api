package notifications

import (
	"context"

	"github.com/jmoiron/sqlx"
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
	MarkAsRead(ctx context.Context, id string) error
}

// RepositoryInterface defines the data access layer for managing notifications.
type RepositoryInterface interface {
	GetDB() *sqlx.DB
	GetByUserID(ctx context.Context, userID string) ([]NotificationModel, error)
	MarkAsRead(ctx context.Context, tx datastore.DB, id string) error
	Create(
		ctx context.Context,
		tx datastore.DB,
		notif *NotificationModel,
	) error
}
