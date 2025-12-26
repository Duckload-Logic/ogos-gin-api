package excuseslips

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
