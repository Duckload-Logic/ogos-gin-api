package m2mclients

import (
	"database/sql"
	"time"
)

// M2MClientDB represents the database model for the m2m_clients table.
type M2MClientDB struct {
	ID                int            `db:"id"`
	UserID            string         `db:"user_id"`
	ClientName        string         `db:"client_name"`
	ClientID          string         `db:"client_id"`
	ClientSecretHash  string         `db:"client_secret_hash"`
	ClientDescription string         `db:"client_description"`
	Scopes            sql.NullString `db:"scopes"`
	IsActive          bool           `db:"is_active"`
	IsVerified        bool           `db:"is_verified"`
	LastUsedAt        sql.NullTime   `db:"last_used_at"`
	ExpiresAt         sql.NullTime   `db:"expires_at"`
	CreatedAt         time.Time      `db:"created_at"`
	UpdatedAt         time.Time      `db:"updated_at"`
}
