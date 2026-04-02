package notifications

import (
	"database/sql"
	"time"
)

type NotificationModel struct {
	ID string `db:"id" json:"id"`

	// Entities involved in the notification
	ReceiverID sql.NullString `db:"receiver_id"     json:"receiverId"`
	ActorID    sql.NullString `db:"actor_id"    json:"actorId,omitempty"`
	TargetID   sql.NullString `db:"target_id"   json:"targetId,omitempty"`
	TargetType sql.NullString `db:"target_type" json:"targetType,omitempty"`

	Title     string    `db:"title"      json:"title"`
	Message   string    `db:"message"    json:"message"`
	Type      string    `db:"type"       json:"type"`
	IsRead    bool      `db:"is_read"    json:"isRead"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
