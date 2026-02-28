package appointments

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/database"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetTimeSlots(ctx context.Context, date string) ([]TimeSlot, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM time_slots
	`, database.GetColumns(TimeSlot{}))

	var slots []TimeSlot
	err := r.db.SelectContext(ctx, &slots, query, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get time slots: %w", err)
	}

	return slots, nil
}

func (r *Repository) GetCategories(ctx context.Context) ([]AppointmentCategory, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM appointment_categories
	`, database.GetColumns(AppointmentCategory{}))

	var categories []AppointmentCategory
	err := r.db.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get concern categories: %w", err)
	}

	return categories, nil
}

func (r *Repository) GetAppointment(ctx context.Context, id int) (*Appointment, error) {
	query := fmt.Sprintf(`
        SELECT %s
        FROM appointments
        WHERE appointment_id = ?
	`, database.GetColumns(Appointment{}))

	var appt Appointment

	err := r.db.GetContext(ctx, &appt, query, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get appointment: %w", err)
	}

	return &appt, nil
}

func (r *Repository) GetDailyStatusCount(ctx context.Context, startDate, endDate string) ([]DailyStatusCount, error) {
	query := `
		SELECT
			DATE(a.when_date) as date,
			COUNT(CASE WHEN s.name = 'Pending' THEN 1 END) as pending_count,
			COUNT(CASE WHEN s.name = 'Scheduled' THEN 1 END) as scheduled_count,
			COUNT(CASE WHEN s.name = 'Rescheduled' THEN 1 END) as rescheduled_count
		FROM appointments a
		JOIN appointment_statuses s ON a.status_id = s.id
		WHERE when_date BETWEEN ? AND ?
		GROUP BY DATE(a.when_date);
	`

	var dsc []DailyStatusCount
	err := r.db.SelectContext(ctx, &dsc, query, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return dsc, nil
}

func (r *Repository) GetTotalAppointmentsCount(ctx context.Context, statusID, startDate, endDate string, userID *int) (int, error) {
	query := `SELECT COUNT(*) FROM appointments WHERE 1=1`
	var args []interface{}

	if statusID != "" {
		query += " AND status_id = ?"
		args = append(args, statusID)
	}
	if startDate != "" {
		query += " AND when_date >= ?"
		args = append(args, startDate)
	}
	if endDate != "" {
		query += " AND when_date <= ?"
		args = append(args, endDate)
	}
	if userID != nil {
		query += " AND user_id = ?"
		args = append(args, *userID)
	}

	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to get appointment count: %w", err)
	}
	return count, nil
}

func (r *Repository) List(
	ctx context.Context, offset, limit int, search, orderBy, statusIDs, startDate, endDate string,
) ([]AppointmentWithDetailsView, error) {
	query := (`
		SELECT
			a.id,
			u.id AS user_id,
			u.first_name AS user_first_name,
			u.middle_name AS user_middle_name,
			u.last_name AS user_last_name,
			u.email AS user_email,
			a.reason AS reason,
			a.admin_notes AS admin_notes,
			a.when_date AS when_date,
			a.created_at AS created_at,
			a.updated_at AS updated_at,
			ts.id AS time_slot_id,
			ts.time AS time_slot_time,
			ac.id AS category_id,
			ac.name AS category_name,
			as2.id AS status_id,
			as2.name AS status_name,
			as2.color_key AS status_color_key
		FROM appointments a
		JOIN users u ON a.user_id = u.id
		JOIN time_slots ts ON a.time_slot_id = ts.id
		JOIN appointment_categories ac ON a.appointment_category_id = ac.id
		JOIN appointment_statuses as2 ON a.status_id = as2.id
		WHERE 1=1
	`)
	var args []interface{}

	if statusIDs != "" {
		statusIDList := strings.Split(statusIDs, ",")
		query += " AND a.status_id IN (?)"
		args = append(args, statusIDList) // Pass the slice directly here
	}
	if startDate != "" {
		query += " AND a.when_date >= ?"
		args = append(args, startDate)
	}
	if endDate != "" {
		query += " AND a.when_date <= ?"
		args = append(args, endDate)
	}
	if search != "" {
		query += " AND (u.first_name LIKE ? OR u.middle_name LIKE ? OR u.last_name LIKE ? OR u.email LIKE ?)"
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm, searchTerm)
	}

	query += fmt.Sprintf(" ORDER BY a.%s LIMIT %d OFFSET %d", orderBy, limit, offset)

	expandedQuery, expandedArgs, err := sqlx.In(query, args...)
	if err != nil {
		return nil, err
	}

	finalQuery := r.db.Rebind(expandedQuery)

	var appts []AppointmentWithDetailsView
	err = r.db.SelectContext(ctx, &appts, finalQuery, expandedArgs...)
	if err != nil {
		return nil, fmt.Errorf("failed to list appointments: %w", err)
	}

	return appts, nil
}

func (r *Repository) GetTimeSlotByID(ctx context.Context, id int) (*TimeSlot, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM time_slots
		WHERE id = ?
	`, database.GetColumns(TimeSlot{}))
	var slot TimeSlot
	err := r.db.GetContext(ctx, &slot, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get time slot by ID: %w", err)
	}

	return &slot, nil
}

func (r *Repository) GetAppointmentCategoryByID(ctx context.Context, id int) (*AppointmentCategory, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM appointment_categories
		WHERE id = ?
	`, database.GetColumns(AppointmentCategory{}))
	var category AppointmentCategory
	err := r.db.GetContext(ctx, &category, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get concern category by ID: %w", err)
	}

	return &category, nil
}

func (r *Repository) GetStatusByID(ctx context.Context, id int) (*AppointmentStatus, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM appointment_statuses
		WHERE id = ?
	`, database.GetColumns(AppointmentStatus{}))

	var status AppointmentStatus
	err := r.db.GetContext(ctx, &status, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointment status by ID: %w", err)
	}

	return &status, nil
}

func (r *Repository) GetAvailableTimeSlots(ctx context.Context, date string) ([]AvailableTimeSlotView, error) {
	query := `
		SELECT
			ts.id as time_slot_id,
            ts.time,
            (a.id IS NULL) as is_available
        FROM time_slots ts
        LEFT JOIN appointments a ON ts.id = a.time_slot_id
            AND a.when_date = ?
            AND a.status_id != (SELECT id FROM appointment_statuses WHERE name = 'Cancelled')
	`

	var slots []AvailableTimeSlotView
	err := r.db.SelectContext(ctx, &slots, query, date)
	if err != nil {
		return nil, fmt.Errorf("failed to get available time slots: %w", err)
	}

	return slots, nil
}

func (r *Repository) GetStatuses(ctx context.Context) ([]AppointmentStatus, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM appointment_statuses
	`, database.GetColumns(AppointmentStatus{}))
	var statuses []AppointmentStatus
	err := r.db.SelectContext(ctx, &statuses, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointment statuses: %w", err)
	}

	return statuses, nil
}

func (r *Repository) ListByUserID(ctx context.Context, userID int, offset, limit int, orderBy string, statusID, startDate, endDate string) ([]AppointmentWithDetailsView, error) {
	query := `
		SELECT
			a.id,
			u.id AS user_id,
			u.first_name AS user_first_name,
			u.middle_name AS user_middle_name,
			u.last_name AS user_last_name,
			u.email AS user_email,
			a.reason AS reason,
			a.admin_notes AS admin_notes,
			a.when_date AS when_date,
			a.created_at AS created_at,
			a.updated_at AS updated_at,
			ts.id AS time_slot_id,
			ts.time AS time_slot_time,
			ac.id AS category_id,
			ac.name AS category_name,
			as2.id AS status_id,
			as2.name AS status_name,
			as2.color_key AS status_color_key
		FROM appointments a
		JOIN users u ON a.user_id = u.id
		JOIN time_slots ts ON a.time_slot_id = ts.id
		JOIN appointment_categories ac ON a.appointment_category_id = ac.id
		JOIN appointment_statuses as2 ON a.status_id = as2.id
		WHERE a.user_id = ?
	`
	args := []interface{}{userID}

	if statusID != "" {
		query += " AND a.status_id = ?"
		args = append(args, statusID)
	}

	if startDate != "" {
		query += " AND a.when_date >= ?"
		args = append(args, startDate)
	}

	if endDate != "" {
		query += " AND a.when_date <= ?"
		args = append(args, endDate)
	}

	query += fmt.Sprintf(" ORDER BY a.%s LIMIT %d OFFSET %d", orderBy, limit, offset)

	var appts []AppointmentWithDetailsView
	err := r.db.SelectContext(ctx, &appts, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointments by user ID: %w", err)
	}

	return appts, nil
}

func (r *Repository) GetAppointmentStats(ctx context.Context, statusID, startDate, endDate string, userID *int) ([]StatusCount, error) {
	// Build the JOIN condition - date filters go in ON clause to preserve LEFT JOIN behavior
	joinCondition := "a.status_id = as2.id"
	var args []interface{}

	if statusID != "" {
		joinCondition += " AND a.status_id = ?"
		args = append(args, statusID)
	}

	if startDate != "" {
		joinCondition += " AND a.when_date >= ?"
		args = append(args, startDate)
	}

	if endDate != "" {
		joinCondition += " AND a.when_date <= ?"
		args = append(args, endDate)
	}

	if userID != nil {
		joinCondition += " AND a.user_id = ?"
		args = append(args, *userID)
	}

	query := fmt.Sprintf(`
		SELECT
			as2.id AS id,
			as2.name AS name,
			COUNT(a.id) AS count
		FROM appointment_statuses as2
		LEFT JOIN appointments a ON %s
		GROUP BY as2.id, as2.name
		ORDER BY as2.id
	`, joinCondition)

	var counts []StatusCount
	err := r.db.SelectContext(ctx, &counts, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointment stats: %w", err)
	}

	return counts, nil
}

func (r *Repository) CreateAppointment(ctx context.Context, appt *Appointment) error {
	cols, vals := database.GetInsertStatement(appt, []string{"updated_at"})
	query := fmt.Sprintf(`
		INSERT INTO appointments (%s)
		VALUES (%s)
	`, cols, vals)

	result, err := r.db.NamedExecContext(ctx, query, appt)
	if err != nil {
		return fmt.Errorf("failed to upsert appointment: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	appt.ID = int(id)
	return nil
}

func (r *Repository) UpdateAppointment(ctx context.Context, appt Appointment) error {
	return database.RunInTransaction(ctx, r.db, func(txn *sqlx.Tx) error {
		query := `UPDATE appointments SET `

		var args []interface{}
		var setQuery []string
		if appt.Reason.Valid {
			setQuery = append(setQuery, " reason = ?")
			args = append(args, appt.Reason.String)
		}

		if appt.WhenDate != "" {
			setQuery = append(setQuery, " when_date = ?")
			args = append(args, appt.WhenDate)
		}

		if appt.TimeSlotID != 0 {
			setQuery = append(setQuery, " time_slot_id = ?")
			args = append(args, appt.TimeSlotID)
		}

		if appt.AppointmentCategoryID != 0 {
			setQuery = append(setQuery, " appointment_category_id = ?")
			args = append(args, appt.AppointmentCategoryID)
		}

		if appt.StatusID != 0 {
			setQuery = append(setQuery, " status_id = ?")
			args = append(args, appt.StatusID)
		}

		if appt.AdminNotes.Valid {
			setQuery = append(setQuery, " admin_notes = ?")
			args = append(args, appt.AdminNotes.String)
		}

		query += strings.Join(setQuery, ",")

		query += " WHERE id = ?"
		args = append(args, appt.ID)

		log.Println("Update query:", query, "Args:", args)

		_, err := txn.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}

		return nil
	})
}
