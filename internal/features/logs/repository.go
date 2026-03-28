package logs

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetDB() *sqlx.DB {
	return r.db
}

// Record inserts a new system log entry
func (r *Repository) Record(
	ctx context.Context,
	tx datastore.DB,
	log *SystemLog,
) error {
	cols, vals := datastore.GetInsertStatement(log, []string{"created_at"})
	query := fmt.Sprintf(`
		INSERT INTO system_logs (%s)
		VALUES (%s)
	`, cols, vals)

	_, err := tx.NamedExecContext(ctx, query, log)
	if err != nil {
		return fmt.Errorf("failed to insert system log: %w", err)
	}

	return nil
}

// List retrieves system log entries with filtering and pagination
func (r *Repository) List(
	ctx context.Context, offset, limit int,
	category, action, userEmail, targetType, targetEmail,
	search, startDate, endDate, orderBy string,
) ([]SystemLog, error) {
	query, args := r.applyLogFilters(
		fmt.Sprintf(
			"SELECT %s FROM system_logs WHERE 1=1",
			datastore.GetColumns(&SystemLog{}),
		),
		nil,
		category,
		action,
		userEmail,
		targetType,
		targetEmail,
		search,
		startDate,
		endDate,
	)

	if orderBy == "" {
		orderBy = "created_at"
	}

	query += fmt.Sprintf(" ORDER BY %s DESC LIMIT ? OFFSET ?", orderBy)
	args = append(args, limit, offset)

	var logs []SystemLog
	err := r.db.SelectContext(ctx, &logs, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list system logs: %w", err)
	}

	return logs, nil
}

// GetTotalCount returns the total count of system log entries matching filters
func (r *Repository) GetTotalCount(
	ctx context.Context,
	category, action, userEmail, targetType, targetEmail,
	search, startDate, endDate string,
) (int, error) {
	query, args := r.applyLogFilters(
		"SELECT COUNT(*) FROM system_logs WHERE 1=1",
		nil,
		category, action, userEmail, targetType, targetEmail,
		search, startDate, endDate,
	)

	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to count system logs: %w", err)
	}

	return count, nil
}

func (r *Repository) applyLogFilters(
	query string,
	args []interface{},
	category, action, userEmail, targetType, targetEmail,
	search, startDate, endDate string,
) (string, []interface{}) {
	if args == nil {
		args = []interface{}{}
	}

	if category != "" {
		query += " AND category = ?"
		args = append(args, category)
	}

	if action != "" {
		query += " AND action = ?"
		args = append(args, action)
	}

	if userEmail != "" {
		query += " AND (user_email = ? OR target_email = ?)"
		args = append(args, userEmail, userEmail)
	}

	if targetType != "" {
		query += " AND target_type = ?"
		args = append(args, targetType)
	}

	if targetEmail != "" {
		query += " AND target_email = ?"
		args = append(args, targetEmail)
	}

	if search != "" {
		query += " AND (message LIKE ? OR action LIKE ? OR user_email LIKE ?)"
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
	}

	if startDate != "" {
		query += " AND created_at >= ?"
		args = append(args, startDate)
	}

	if endDate != "" {
		query += " AND created_at <= ?"
		args = append(args, endDate+" 23:59:59")
	}

	return query, args
}

// GetStats returns log counts grouped by category
func (r *Repository) GetStats(
	ctx context.Context,
	startDate, endDate string,
) ([]LogStatsDTO, error) {
	query, args := r.applyLogFilters(
		"SELECT category, COUNT(*) as count FROM system_logs WHERE 1=1",
		nil,
		"", "", "", "", "", "", startDate, endDate,
	)

	query += " GROUP BY category ORDER BY category"

	var stats []LogStatsDTO
	err := r.db.SelectContext(ctx, &stats, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get log stats: %w", err)
	}

	return stats, nil
}

// toNullString converts a value to sql.NullString
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
