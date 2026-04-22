package notifications

import (
	"database/sql"
	"time"
)

// NotificationDB represents the database model for the notifications table.
type NotificationDB struct {
	ID         string         `db:"id"`
	ReceiverID sql.NullString `db:"receiver_id"`
	ActorID    sql.NullString `db:"actor_id"`
	TargetID   sql.NullString `db:"target_id"`
	TargetType sql.NullString `db:"target_type"`
	Title      string         `db:"title"`
	Message    string         `db:"message"`
	Type       string         `db:"type"`
	IsRead     bool           `db:"is_read"`
	CreatedAt  time.Time      `db:"created_at"`
	UpdatedAt  time.Time      `db:"updated_at"`
}
