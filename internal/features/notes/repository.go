package notes

import (
	"context"
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

func (r *Repository) GetStudentSignificantNotes(
	ctx context.Context,
	iirID int,
) ([]SignificantNote, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM significant_notes
		WHERE iir_id = ?
	`, database.GetColumns(SignificantNote{}))

	var notes []SignificantNote
	err := r.db.SelectContext(ctx, &notes, query, iirID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get student significant notes: %w",
			err,
		)
	}

	return notes, nil
}

func (r *Repository) CreateSignificantNote(
	ctx context.Context,
	sn *SignificantNote,
) (int, error) {
	return database.NewRunInTransaction(
		ctx,
		r.db,
		func(tx *sqlx.Tx) (int, error) {
			cols, vals := database.GetInsertStatement(
				SignificantNote{},
				[]string{"created_at", "updated_at"},
			)

			query := fmt.Sprintf(`
				INSERT INTO significant_notes (%s)
				VALUES (%s)
			`, cols, vals)

			result, err := tx.NamedExecContext(
				ctx,
				query,
				sn,
			)
			if err != nil {
				return 0, fmt.Errorf(
					"failed to create significant note: %w",
					err,
				)
			}

			lastID, err := result.LastInsertId()
			if err != nil {
				return 0, fmt.Errorf(
					"failed to get last insert ID: %w",
					err,
				)
			}

			return int(lastID), nil
		},
	)
}
