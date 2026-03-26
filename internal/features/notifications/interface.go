package notifications

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

// ServiceInterface defines the business logic for managing notifications.
type ServiceInterface interface {
	Send(
		ctx context.Context,
		userID string,
		title, message, notifType string,
	) error
	GetUserNotifications(
		ctx context.Context,
		userID string,
	) ([]NotificationDTO, error)
	MarkAsRead(ctx context.Context, id int) error
}

// RepositoryInterface defines the data access layer for managing notifications.
type RepositoryInterface interface {
	GetDB() *sqlx.DB
	GetByUserID(ctx context.Context, userID string) ([]NotificationModel, error)
	MarkAsRead(ctx context.Context, tx datastore.DB, id int) error
	Create(
		ctx context.Context,
		tx datastore.DB,
		userID string,
		title, message, notifType string,
	) error
}
