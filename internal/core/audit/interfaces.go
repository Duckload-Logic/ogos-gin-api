package audit

import (
	"context"

	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

// Logger defines the interface for recording system logs.
type Logger interface {
	Record(ctx context.Context, tx datastore.DB, entry LogEntry)
}

// Notifier defines the interface for sending notifications.
type Notifier interface {
	Send(ctx context.Context, notif NotificationEntry) error
}

// UserGetter defines the interface for fetching user IDs by role.
type UserGetter interface {
	GetUserIDsByRole(ctx context.Context, roleID int) ([]string, error)
}
