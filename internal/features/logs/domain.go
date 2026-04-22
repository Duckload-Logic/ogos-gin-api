package logs

import (
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

// SystemLog represents a pure business entity for a system log entry.
type SystemLog struct {
	ID          int
	Level       string
	Category    string
	Action      string
	Message     string
	UserID      structs.NullableString
	TargetID    structs.NullableString
	TargetType  structs.NullableString
	UserEmail   structs.NullableString
	TargetEmail structs.NullableString
	IPAddress   structs.NullableString
	UserAgent   structs.NullableString
	Metadata    structs.NullableString
	TraceID     structs.NullableString
	CreatedAt   time.Time
}
