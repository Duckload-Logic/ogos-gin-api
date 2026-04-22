package notifications

import (
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

// Notification represents a pure business entity for system notifications.
type Notification struct {
	ID         string
	ReceiverID structs.NullableString
	ActorID    structs.NullableString
	TargetID   structs.NullableString
	TargetType structs.NullableString
	Title      string
	Message    string
	Type       string
	IsRead     bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
