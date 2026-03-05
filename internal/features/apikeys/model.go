package apikeys

import (
	"database/sql"
	"time"
)

type APIKey struct {
	ID         int            `db:"id"`
	Name       string         `db:"name"`
	KeyHash    string         `db:"key_hash"`
	KeyPrefix  string         `db:"key_prefix"`
	Scopes     sql.NullString `db:"scopes"`
	IsActive   bool           `db:"is_active"`
	LastUsedAt sql.NullTime   `db:"last_used_at"`
	ExpiresAt  sql.NullTime   `db:"expires_at"`
	CreatedAt  time.Time      `db:"created_at"`
	UpdatedAt  time.Time      `db:"updated_at"`
}
