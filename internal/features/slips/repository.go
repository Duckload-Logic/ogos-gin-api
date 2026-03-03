package slips

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/database"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateSlip(ctx context.Context, slip *Slip) (*int, error) {
	cols, vals := database.GetInsertStatement(slip, []string{"id", "updated_at"})
	query := fmt.Sprintf(`
		INSERT INTO admission_slips (%s)
		VALUES (%s)
	`, cols, vals)

	result, err := r.db.NamedExecContext(ctx, query, slip)
	if err != nil {
		return nil, fmt.Errorf("failed to insert excuse slip: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	slip.ID = int(id)
	return &slip.ID, nil
}

func (r *Repository) SaveSlipAttachment(ctx context.Context, attachment *SlipAttachment) error {
	cols, vals := database.GetInsertStatement(attachment, []string{"id", "updated_at"})
	query := fmt.Sprintf(`
		INSERT INTO slip_attachments (%s)
		VALUES (%s)
	`, cols, vals)
	_, err := r.db.NamedExecContext(ctx, query, attachment)
	if err != nil {
		return fmt.Errorf("failed to save slip attachment: %w", err)
	}

	return nil
}

func (r *Repository) CheckStudentExistence(ctx context.Context, studentID int) (bool, error) {
	query := `SELECT count(*) FROM student_records WHERE student_record_id = ?`

	var count int
	err := r.db.QueryRowContext(ctx, query, studentID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("database error checking student: %w", err)
	}

	return count > 0, nil
}

func (r *Repository) GetSlipStatuses(ctx context.Context) ([]SlipStatus, error) {
	var statuses []SlipStatus
	query := fmt.Sprintf(`
		SELECT %s FROM statuses WHERE status_type IN ('slip', 'both')
	`, database.GetColumns(SlipStatus{}))
	err := r.db.SelectContext(ctx, &statuses, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get slip statuses: %w", err)
	}

	return statuses, nil
}

func (r *Repository) GetSlipCategories(ctx context.Context) ([]SlipCategory, error) {
	var categories []SlipCategory
	query := fmt.Sprintf(`
		SELECT %s FROM admission_slip_categories
	`, database.GetColumns(SlipCategory{}))
	err := r.db.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get slip categories: %w", err)
	}

	return categories, nil
}

func (r *Repository) GetSlipStats(ctx context.Context, userID *int, req *ListSlipRequest) ([]SlipStatusCount, error) {
	var args []interface{}
	filterConditions := "1=1"

	if req.StatusID != 0 {
		filterConditions += " AND es.status_id = ?"
		args = append(args, req.StatusID)
	}

	if req.StatusID != 0 {
		filterConditions += " AND es.status_id = ?"
		args = append(args, req.StatusID)
	}

	if req.StartDate != "" {
		filterConditions += " AND es.created_at >= ?"
		args = append(args, req.StartDate)
	}

	if req.EndDate != "" {
		filterConditions += " AND es.created_at <= ?"
		args = append(args, req.EndDate)
	}

	if userID != nil {
		filterConditions += " AND es.user_id = ?"
		args = append(args, userID)
	}

	var counts []SlipStatusCount
	query := fmt.Sprintf(`
		SELECT
			s.id AS id,
			s.name AS name,
			s.color_key AS color_key,
			COUNT(es.id) AS count
		FROM statuses s
		LEFT JOIN admission_slips es ON s.id = es.status_id AND %s
		WHERE s.status_type IN ('slip', 'both')
		GROUP BY s.id, s.name, s.color_key
	`, filterConditions)
	err := r.db.SelectContext(ctx, &counts, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get slip status counts: %w", err)
	}

	return counts, nil
}

func (r *Repository) GetTotalSlipsCount(ctx context.Context, req *ListSlipRequest) (int, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM admission_slips es
		JOIN users u ON es.user_id = u.id
		WHERE 1=1
	`

	var args []interface{}
	if req.StatusID != 0 {
		query += " AND es.status_id = ?"
		args = append(args, req.StatusID)
	}

	if req.StartDate != "" {
		query += " AND es.created_at >= ?"
		args = append(args, req.StartDate)
	}

	if req.EndDate != "" {
		query += " AND es.created_at <= ?"
		args = append(args, req.EndDate)
	}

	if req.Search != "" {
		query += " AND (es.reason LIKE ? OR u.first_name LIKE ? OR u.last_name LIKE ?)"
		searchTerm := "%" + req.Search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
	}

	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to count slips: %w", err)
	}

	return count, nil
}

func (r *Repository) GetTotalUrgentSlipsCount(ctx context.Context, req *ListSlipRequest) (int, error) {
	var count int
	// Notice: We must match the JOINs from GetUrgentSlips if they filter the results
	query := `
        SELECT COUNT(*)
        FROM admission_slips es
        JOIN admission_slip_categories c ON es.category_id = c.id
        WHERE es.status_id IN (1, 9)
    `

	var args []interface{}

	if req.StartDate != "" {
		query += " AND es.date_needed >= ?"
		args = append(args, req.StartDate)
	}

	if req.EndDate != "" {
		query += " AND es.date_needed <= ?"
		args = append(args, req.EndDate)
	}

	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to count urgent slips: %w", err)
	}

	return count, nil
}

func (r *Repository) GetUrgentSlips(ctx context.Context, req *ListSlipRequest) ([]SlipWithDetailsView, error) {
	var slips []SlipWithDetailsView
	query := `
		SELECT
			es.id AS id,
			es.user_id AS user_id,
			u.first_name AS user_first_name,
			u.last_name AS user_last_name,
			u.email AS user_email,
			es.reason AS reason,
			es.date_of_absence AS date_of_absence,
			es.date_needed AS date_needed,
			es.admin_notes AS admin_notes,
			c.id AS category_id,
			c.name AS category_name,
			s.id AS status_id,
			s.name AS status_name,
			s.color_key AS status_color_key,
			es.created_at AS created_at,
			es.updated_at AS updated_at,
			(
				(1000 - DATEDIFF(es.date_needed, CURRENT_DATE)) * 10
				+
				CASE WHEN es.category_id = 1 THEN 500 ELSE 0 END
			) AS urgency_score
		FROM admission_slips es
		JOIN users u ON es.user_id = u.id
		JOIN admission_slip_categories c ON es.category_id = c.id
		JOIN statuses s ON es.status_id = s.id
		WHERE es.status_id IN (1, 9)
	`

	var args []interface{}
	if req.StartDate != "" {
		query += " AND es.date_needed >= ?"
		args = append(args, req.StartDate)
	}

	if req.EndDate != "" {
		query += " AND es.date_needed <= ?"
		args = append(args, req.EndDate)
	}

	query += `
		ORDER BY
			es.date_needed ASC,
			urgency_score DESC
		LIMIT ? OFFSET ?
	`
	args = append(args, req.PageSize, req.GetOffset())

	err := r.db.SelectContext(ctx, &slips, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get urgent slips: %w", err)
	}

	return slips, nil
}

func (r *Repository) GetAll(ctx context.Context, req *ListSlipRequest) ([]SlipWithDetailsView, error) {
	var slips []SlipWithDetailsView
	query := `
        SELECT
            es.id AS id,
            es.user_id AS user_id,
			u.first_name AS user_first_name,
			u.middle_name AS user_middle_name,
			u.last_name AS user_last_name,
			u.email AS user_email,
            es.reason AS reason,
            es.date_of_absence AS date_of_absence,
            es.date_needed AS date_needed,
			es.admin_notes AS admin_notes,
			c.id AS category_id,
			c.name AS category_name,
			s.id AS status_id,
			s.name AS status_name,
			s.color_key AS status_color_key,
            es.created_at AS created_at,
            es.updated_at AS updated_at
        FROM admission_slips es
		JOIN users u ON es.user_id = u.id
		JOIN admission_slip_categories c ON es.category_id = c.id
		JOIN statuses s ON es.status_id = s.id
		WHERE 1=1
    `
	var args []interface{}
	if req.StatusID != 0 {
		query += " AND es.status_id = ?"
		args = append(args, req.StatusID)
	}
	if req.StartDate != "" {
		query += " AND es.created_at >= ?"
		args = append(args, req.StartDate)
	}
	if req.EndDate != "" {
		query += " AND es.created_at <= ?"
		args = append(args, req.EndDate)
	}
	if req.Search != "" {
		query += " AND (es.reason LIKE ? OR u.first_name LIKE ? OR u.last_name LIKE ?)"
		searchTerm := "%" + req.Search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
	}

	query += `
		ORDER BY es.created_at DESC
		LIMIT ? OFFSET ?
	`
	args = append(args, req.PageSize, req.GetOffset())

	err := r.db.SelectContext(ctx, &slips, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get excuse slips: %w", err)
	}

	return slips, nil
}

func (r *Repository) GetByUserID(ctx context.Context, userID int, req *ListSlipRequest) ([]SlipWithDetailsView, error) {
	var slips []SlipWithDetailsView
	var args []interface{}
	args = append(args, userID)
	query := `
		SELECT
			es.id AS id,
			es.user_id AS user_id,
			u.first_name AS user_first_name,
			u.middle_name AS user_middle_name,
			u.last_name AS user_last_name,
			es.reason AS reason,
			es.date_of_absence AS date_of_absence,
			es.date_needed AS date_needed,
			es.admin_notes AS admin_notes,
			c.id AS category_id,
			c.name AS category_name,
			s.id AS status_id,
			s.name AS status_name,
			s.color_key AS status_color_key,
			es.created_at AS created_at,
			es.updated_at AS updated_at
		FROM admission_slips es
		JOIN users u ON es.user_id = u.id
		JOIN admission_slip_categories c ON es.category_id = c.id
		JOIN statuses s ON es.status_id = s.id
		WHERE es.user_id = ?
    `

	if req.StatusID != 0 {
		query += " AND es.status_id = ?"
		args = append(args, req.StatusID)
	}

	if req.StartDate != "" {
		query += " AND es.created_at >= ?"
		args = append(args, req.StartDate)
	}

	if req.EndDate != "" {
		query += " AND es.created_at <= ?"
		args = append(args, req.EndDate)
	}

	query += `
		ORDER BY es.created_at DESC
		LIMIT ? OFFSET ?
	`
	args = append(args, req.PageSize, req.GetOffset())

	err := r.db.SelectContext(ctx, &slips, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get excuse slips for user: %w", err)
	}

	return slips, nil
}

func (r *Repository) GetSlipAttachments(ctx context.Context, slipID int) ([]SlipAttachment, error) {
	var attachments []SlipAttachment
	query := fmt.Sprintf(`
		SELECT %s
		FROM slip_attachments
		WHERE admission_slip_id = ?
	`, database.GetColumns(SlipAttachment{}))
	err := r.db.SelectContext(ctx, &attachments, query, slipID)
	if err != nil {
		return nil, fmt.Errorf("failed to get slip attachments: %w", err)
	}

	return attachments, nil
}

func (r *Repository) GetAttachmentByID(ctx context.Context, attachmentID int) (*SlipAttachment, error) {
	var attachment SlipAttachment
	query := fmt.Sprintf(`
		SELECT %s
		FROM slip_attachments
		WHERE id = ?
	`, database.GetColumns(SlipAttachment{}))
	err := r.db.GetContext(ctx, &attachment, query, attachmentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found, return nil without error
		}
		return nil, fmt.Errorf("failed to get attachment by ID: %w", err)
	}

	return &attachment, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, id int, statusName string, adminNotes string) error {
	// First, get the status ID from the status name
	var statusID int
	query := `SELECT id FROM statuses WHERE name = ?`
	err := r.db.GetContext(ctx, &statusID, query, statusName)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("status '%s' not found", statusName)
		}
		return fmt.Errorf("failed to get status ID: %w", err)
	}

	// Now update the slip with the status ID and admin notes
	updateQuery := `
		UPDATE admission_slips
		SET status_id = ?, admin_notes = ?, updated_at = NOW()
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, updateQuery, statusID, adminNotes, id)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		// This usually means the ID doesn't exist
		return sql.ErrNoRows
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM excuse_slips WHERE excuse_slip_id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete excuse slip: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
