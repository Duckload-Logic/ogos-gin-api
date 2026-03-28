package logs

import (
	"database/sql"
	"time"
)

// SystemLog represents a row in the system_logs table
type SystemLog struct {
	ID       int    `db:"id"`
	Level    string `db:"level"`
	Category string `db:"category"`
	Action   string `db:"action"`
	Message  string `db:"message"`

	UserID     sql.NullString `db:"user_id"`
	TargetID   sql.NullString `db:"target_id"`
	TargetType sql.NullString `db:"target_type"`

	UserEmail   sql.NullString `db:"user_email"`
	TargetEmail sql.NullString `db:"target_email"`

	IPAddress sql.NullString `db:"ip_address"`
	UserAgent sql.NullString `db:"user_agent"`
	Metadata  sql.NullString `db:"metadata"`
	TraceID   sql.NullString `db:"trace_id"`
	CreatedAt time.Time      `db:"created_at"`
}
