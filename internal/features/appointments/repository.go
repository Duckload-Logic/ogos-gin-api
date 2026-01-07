package appointments

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, appt *Appointment) error {

	query := `
		INSERT INTO appointments (
			user_id, 
			reason, 
			scheduled_date, 
			scheduled_time, 
			concern_category, 
			status, 
			created_at, 
			updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		appt.UserID,
		appt.Reason,
		appt.ScheduledDate, // scheduled_date
		appt.ScheduledTime, // scheduled_time
		appt.ConcernCategory,
		appt.Status,
	)
	if err != nil {
		return fmt.Errorf("failed to insert appointment: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	appt.ID = int(id)
	return nil
}

func (r *Repository) GetByID(ctx context.Context, id int) (*Appointment, error) {
	query := `
        SELECT 
            appointment_id, user_id, reason, 
            scheduled_date, scheduled_time, concern_category, 
            status, created_at, updated_at
        FROM appointments
        WHERE appointment_id = ?
    `

	var appt Appointment

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&appt.ID,
		&appt.UserID,
		&appt.Reason,
		&appt.ScheduledDate,
		&appt.ScheduledTime,
		&appt.ConcernCategory,
		&appt.Status,
		&appt.CreatedAt,
		&appt.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get appointment: %w", err)
	}

	return &appt, nil
}

func (r *Repository) List(ctx context.Context, status string, date string) ([]Appointment, error) {
	query := `
        SELECT 
            appointment_id, user_id, reason, 
            scheduled_date, scheduled_time, concern_category, 
            status, created_at, updated_at
        FROM appointments
        WHERE 1=1
    `
	var args []interface{}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}
	if date != "" {
		query += " AND scheduled_date = ?"
		args = append(args, date)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list appointments: %w", err)
	}
	defer rows.Close()

	var appts []Appointment
	for rows.Next() {
		var appt Appointment
		if err := rows.Scan(
			&appt.ID,
			&appt.UserID,
			&appt.Reason,
			&appt.ScheduledDate,
			&appt.ScheduledTime,
			&appt.ConcernCategory,
			&appt.Status,
			&appt.CreatedAt,
			&appt.UpdatedAt,
		); err != nil {
			return nil, err
		}

		appts = append(appts, appt)
	}

	return appts, nil
}

func (r *Repository) GetTimeSlots(ctx context.Context, date string) ([]Appointment, error) {
	query := `
		SELECT
			appointment_id, user_id, reason,
			scheduled_date, scheduled_time, concern_category,
			status, created_at, updated_at
		FROM appointments
		WHERE scheduled_date = ? AND status = 'Approved' OR status = 'Pending'
	`

	rows, err := r.db.QueryContext(ctx, query, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get available time slots: %w", err)
	}
	defer rows.Close()

	var appts []Appointment
	for rows.Next() {
		var appt Appointment

		if err := rows.Scan(
			&appt.ID,
			&appt.UserID,
			&appt.Reason,
			&appt.ScheduledDate,
			&appt.ScheduledTime,
			&appt.ConcernCategory,
			&appt.Status,
			&appt.CreatedAt,
			&appt.UpdatedAt,
		); err != nil {
			return nil, err
		}

		appts = append(appts, appt)
	}

	return appts, nil
}

func (r *Repository) GetByUserID(ctx context.Context, userID int) ([]Appointment, error) {
	query := `
        SELECT 
            appointment_id, user_id, reason, 
            scheduled_date, scheduled_time, concern_category, 
            status, created_at, updated_at
        FROM appointments
        WHERE user_id = ?
    `

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student appointments: %w", err)
	}
	defer rows.Close()

	var appts []Appointment
	for rows.Next() {
		var appt Appointment

		if err := rows.Scan(
			&appt.ID,
			&appt.UserID,
			&appt.Reason,
			&appt.ScheduledDate,
			&appt.ScheduledTime,
			&appt.ConcernCategory,
			&appt.Status,
			&appt.CreatedAt,
			&appt.UpdatedAt,
		); err != nil {
			return nil, err
		}

		appts = append(appts, appt)
	}

	return appts, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, id int, status string) error {
	query := `UPDATE appointments SET status = ? WHERE appointment_id = ?`

	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
