package trails

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// Record inserts a new audit trail entry into the database.
// This is the core method called by the service layer.
func (r *Repository) Record(ctx context.Context, trail *AuditTrail) error {
	query := `
		INSERT INTO audit_trails (user_email, action, entity_type, entity_id, old_values, new_values, ip_address, user_agent)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		trail.UserEmail,
		trail.Action,
		trail.EntityType,
		trail.EntityID,
		trail.OldValues,
		trail.NewValues,
		trail.IPAddress,
		trail.UserAgent,
	)
	if err != nil {
		return fmt.Errorf("failed to insert audit trail: %w", err)
	}

	return nil
}

// List retrieves audit trail entries with filtering and pagination
func (r *Repository) List(
	ctx context.Context, offset, limit int,
	action, entityType string, entityID int, userEmail string,
	search, startDate, endDate, orderBy string,
) ([]AuditTrailWithUserView, error) {
	query := `
		SELECT
			at.id,
			at.user_email,
			u.first_name AS user_first_name,
			u.middle_name AS user_middle_name,
			u.last_name AS user_last_name,
			at.action,
			at.entity_type,
			at.entity_id,
			at.old_values,
			at.new_values,
			at.ip_address,
			at.user_agent,
			at.created_at
		FROM audit_trails at
		LEFT JOIN users u ON at.user_email = u.email
		WHERE 1=1
	`

	var args []interface{}

	if action != "" {
		query += " AND at.action = ?"
		args = append(args, action)
	}

	if entityType != "" {
		query += " AND at.entity_type = ?"
		args = append(args, entityType)
	}

	if entityID != 0 {
		query += " AND at.entity_id = ?"
		args = append(args, entityID)
	}

	if userEmail != "" {
		query += " AND at.user_email = ?"
		args = append(args, userEmail)
	}

	if search != "" {
		query += " AND (u.first_name LIKE ? OR u.last_name LIKE ? OR at.entity_type LIKE ?)"
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
	}

	if startDate != "" {
		query += " AND at.created_at >= ?"
		args = append(args, startDate)
	}

	if endDate != "" {
		query += " AND at.created_at <= ?"
		args = append(args, endDate)
	}

	if orderBy == "" {
		orderBy = "created_at"
	}

	query += fmt.Sprintf(" ORDER BY at.%s DESC LIMIT %d OFFSET %d", orderBy, limit, offset)

	var trails []AuditTrailWithUserView
	err := r.db.SelectContext(ctx, &trails, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list audit trails: %w", err)
	}

	return trails, nil
}

// GetTotalCount returns the total count of audit trail entries matching filters
func (r *Repository) GetTotalCount(
	ctx context.Context,
	action, entityType string, entityID int, userEmail string,
	search, startDate, endDate string,
) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM audit_trails at
		LEFT JOIN users u ON at.user_email = u.email
		WHERE 1=1
	`

	var args []interface{}

	if action != "" {
		query += " AND at.action = ?"
		args = append(args, action)
	}

	if entityType != "" {
		query += " AND at.entity_type = ?"
		args = append(args, entityType)
	}

	if entityID != 0 {
		query += " AND at.entity_id = ?"
		args = append(args, entityID)
	}

	if userEmail != "" {
		query += " AND at.user_email = ?"
		args = append(args, userEmail)
	}

	if search != "" {
		query += " AND (u.first_name LIKE ? OR u.last_name LIKE ? OR at.entity_type LIKE ?)"
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
	}

	if startDate != "" {
		query += " AND at.created_at >= ?"
		args = append(args, startDate)
	}

	if endDate != "" {
		query += " AND at.created_at <= ?"
		args = append(args, endDate)
	}

	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to count audit trails: %w", err)
	}

	return count, nil
}

// GetByEntityRef retrieves audit trails for a specific entity
func (r *Repository) GetByEntityRef(
	ctx context.Context, entityType string, entityID int,
) ([]AuditTrailWithUserView, error) {
	query := `
		SELECT
			at.id,
			at.user_email,
			u.first_name AS user_first_name,
			u.middle_name AS user_middle_name,
			u.last_name AS user_last_name,
			at.action,
			at.entity_type,
			at.entity_id,
			at.old_values,
			at.new_values,
			at.ip_address,
			at.user_agent,
			at.created_at
		FROM audit_trails at
		LEFT JOIN users u ON at.user_email = u.email
		WHERE at.entity_type = ? AND at.entity_id = ?
		ORDER BY at.created_at DESC
	`

	var trails []AuditTrailWithUserView
	err := r.db.SelectContext(ctx, &trails, query, entityType, entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit trails for entity: %w", err)
	}

	return trails, nil
}

// toNullString converts a nullable JSON value to sql.NullString
func toNullString(v interface{}) sql.NullString {
	if v == nil {
		return sql.NullString{Valid: false}
	}

	bytes, err := json.Marshal(v)
	if err != nil {
		return sql.NullString{Valid: false}
	}

	return sql.NullString{String: string(bytes), Valid: true}
}
