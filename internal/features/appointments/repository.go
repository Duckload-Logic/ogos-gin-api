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

	dateStr := appt.ScheduledAt.Format("2006-01-02")
	timeStr := appt.ScheduledAt.Format("15:04:05")

	query := `
		INSERT INTO appointments (
			student_record_id, 
			appointment_type_id, 
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
		appt.StudentRecordID,
		appt.AppointmentTypeID,
		dateStr, // scheduled_date
		timeStr, // scheduled_time
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
