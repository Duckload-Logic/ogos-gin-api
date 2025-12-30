package appointments

import (
	"context"
	"database/sql"
	"fmt"
	"time"
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

func (r *Repository) GetByID(ctx context.Context, id int) (*Appointment, error) {
    query := `
        SELECT 
            id, student_record_id, counselor_user_id, appointment_type_id, 
            scheduled_date, scheduled_time, concern_category, 
            status, created_at, updated_at
        FROM appointments
        WHERE id = ?
    `

    var appt Appointment
    var dateStr, timeStr string

    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &appt.ID,
        &appt.StudentRecordID,
        &appt.CounselorUserID,
        &appt.AppointmentTypeID,
        &dateStr, 
        &timeStr,
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

    fullTimeStr := dateStr + " " + timeStr
    parsedTime, _ := time.Parse("2006-01-02 15:04:05", fullTimeStr)
    appt.ScheduledAt = parsedTime

    return &appt, nil
}

func (r *Repository) List(ctx context.Context, status string, date string) ([]Appointment, error) {
    query := `
        SELECT 
            id, student_record_id, counselor_user_id, appointment_type_id, 
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
        var dateStr, timeStr string

        if err := rows.Scan(
            &appt.ID,
            &appt.StudentRecordID,
            &appt.CounselorUserID,
            &appt.AppointmentTypeID,
            &dateStr,
            &timeStr,
            &appt.ConcernCategory,
            &appt.Status,
            &appt.CreatedAt,
            &appt.UpdatedAt,
        ); err != nil {
            return nil, err
        }

        // Combine date and time
        fullTimeStr := dateStr + " " + timeStr
        parsedTime, _ := time.Parse("2006-01-02 15:04:05", fullTimeStr)
        appt.ScheduledAt = parsedTime

        appts = append(appts, appt)
    }

    return appts, nil
}

func (r *Repository) GetByStudentID(ctx context.Context, studentID int) ([]Appointment, error) {
    query := `
        SELECT 
            id, student_record_id, counselor_user_id, appointment_type_id, 
            scheduled_date, scheduled_time, concern_category, 
            status, created_at, updated_at
        FROM appointments
        WHERE student_record_id = ?
    `

    rows, err := r.db.QueryContext(ctx, query, studentID)
    if err != nil {
        return nil, fmt.Errorf("failed to get student appointments: %w", err)
    }
    defer rows.Close()

    var appts []Appointment
    for rows.Next() {
        var appt Appointment
        var dateStr, timeStr string

        if err := rows.Scan(
            &appt.ID,
            &appt.StudentRecordID,
            &appt.CounselorUserID,
            &appt.AppointmentTypeID,
            &dateStr,
            &timeStr,
            &appt.ConcernCategory,
            &appt.Status,
            &appt.CreatedAt,
            &appt.UpdatedAt,
        ); err != nil {
            return nil, err
        }

        fullTimeStr := dateStr + " " + timeStr
        parsedTime, _ := time.Parse("2006-01-02 15:04:05", fullTimeStr)
        appt.ScheduledAt = parsedTime

        appts = append(appts, appt)
    }

    return appts, nil
}