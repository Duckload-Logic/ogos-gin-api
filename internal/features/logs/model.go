package logs

import (
	"database/sql"
	"time"
)

// SystemLog represents a row in the system_logs table
type SystemLog struct {
	ID        int            `db:"id"`
	Category  string         `db:"category"`
	Action    string         `db:"action"`
	Message   string         `db:"message"`
	UserID    sql.NullInt64  `db:"user_id"`
	UserEmail sql.NullString `db:"user_email"`
	IPAddress sql.NullString `db:"ip_address"`
	UserAgent sql.NullString `db:"user_agent"`
	Metadata  sql.NullString `db:"metadata"`
	CreatedAt time.Time      `db:"created_at"`
}
