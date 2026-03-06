package trails

import (
	"database/sql"
	"time"
)

type AuditTrail struct {
	ID         int            `db:"id" json:"id"`
	UserEmail  sql.NullString `db:"user_email" json:"userEmail"`
	Action     string         `db:"action" json:"action"`
	EntityType string         `db:"entity_type" json:"entityType"`
	EntityID   int            `db:"entity_id" json:"entityId"`
	OldValues  sql.NullString `db:"old_values" json:"oldValues,omitempty"`
	NewValues  sql.NullString `db:"new_values" json:"newValues,omitempty"`
	IPAddress  sql.NullString `db:"ip_address" json:"ipAddress,omitempty"`
	UserAgent  sql.NullString `db:"user_agent" json:"userAgent,omitempty"`
	CreatedAt  time.Time      `db:"created_at" json:"createdAt"`
}

// AuditTrailWithUserView is a flattened view for listing with user info
type AuditTrailWithUserView struct {
	ID             int            `db:"id" json:"id"`
	UserEmail      sql.NullString `db:"user_email" json:"userEmail"`
	UserFirstName  sql.NullString `db:"user_first_name" json:"userFirstName,omitempty"`
	UserMiddleName sql.NullString `db:"user_middle_name" json:"userMiddleName,omitempty"`
	UserLastName   sql.NullString `db:"user_last_name" json:"userLastName,omitempty"`
	Action         string         `db:"action" json:"action"`
	EntityType     string         `db:"entity_type" json:"entityType"`
	EntityID       int            `db:"entity_id" json:"entityId"`
	OldValues      sql.NullString `db:"old_values" json:"oldValues,omitempty"`
	NewValues      sql.NullString `db:"new_values" json:"newValues,omitempty"`
	IPAddress      sql.NullString `db:"ip_address" json:"ipAddress,omitempty"`
	UserAgent      sql.NullString `db:"user_agent" json:"userAgent,omitempty"`
	CreatedAt      time.Time      `db:"created_at" json:"createdAt"`
}
