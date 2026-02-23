package excuseslips

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, slip *ExcuseSlip) error {

	dateStr := slip.Date_of_absence.Format("2006-01-02")

	query := `
		INSERT INTO excuse_slips (
			student_record_id,
			reason,
			date_of_absence,
			file_path,
			excuse_slip_status,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		slip.StudentRecordID,
		slip.Reason,
		dateStr,
		slip.FilePath,
		slip.Status,
	)
	if err != nil {
		return fmt.Errorf("failed to insert excuse slip: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	slip.ID = int(id)
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

func (r *Repository) GetAll(ctx context.Context) ([]*ExcuseSlip, error) {
	// Define the query
	query := `
        SELECT
            excuse_slip_id,
            student_record_id,
            reason,
            date_of_absence,
            file_path,
            excuse_slip_status,
            created_at,
            updated_at
        FROM excuse_slips
        ORDER BY created_at DESC
    `
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query excuse slips: %w", err)
	}
	defer rows.Close()

	var slips []*ExcuseSlip

	for rows.Next() {
		var slip ExcuseSlip
		err := rows.Scan(
			&slip.ID,
			&slip.StudentRecordID,
			&slip.Reason,
			&slip.Date_of_absence,
			&slip.FilePath,
			&slip.Status,
			&slip.CreatedAt,
			&slip.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan excuse slip row: %w", err)
		}
		slips = append(slips, &slip)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return slips, nil
}

// GetByID retrieves an excuse slip by its ID
func (r *Repository) GetByID(ctx context.Context, id int) (*ExcuseSlip, error) {
	query := `
        SELECT
            excuse_slip_id,
            student_record_id,
            reason,
            date_of_absence,
            file_path,
            excuse_slip_status,
            created_at,
            updated_at
        FROM excuse_slips
        WHERE excuse_slip_id = ?
    `

	var slip ExcuseSlip

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&slip.ID,
		&slip.StudentRecordID,
		&slip.Reason,
		&slip.Date_of_absence,
		&slip.FilePath,
		&slip.Status,
		&slip.CreatedAt,
		&slip.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get excuse slip by id: %w", err)
	}

	return &slip, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, id int, status string) error {
	query := `
        UPDATE excuse_slips
        SET excuse_slip_status = ?, updated_at = NOW()
        WHERE excuse_slip_id = ?
    `

	result, err := r.db.ExecContext(ctx, query, status, id)
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
