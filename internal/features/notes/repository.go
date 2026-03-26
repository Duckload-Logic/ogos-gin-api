package notes

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetStudentSignificantNotes(
	ctx context.Context,
	iirID string,
) ([]SignificantNote, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM significant_notes
		WHERE iir_id = ?
	`, datastore.GetColumns(SignificantNote{}))

	var notes []SignificantNote
	err := r.db.SelectContext(ctx, &notes, query, iirID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get student significant notes: %w",
			err,
		)
	}

	log.Printf(
		"[GetStudentSignificantNotes] {Database Query}: Retrieved %d notes for IIR ID %s",
		len(notes),
		iirID,
	)

	return notes, nil
}

func (r *Repository) CreateSignificantNote(
	ctx context.Context,
	sn *SignificantNote,
) (string, error) {
	return datastore.NewRunInTransaction(
		ctx,
		r.db,
		func(tx datastore.DB) (string, error) {
			cols, vals := datastore.GetInsertStatement(
				SignificantNote{},
				[]string{"created_at", "updated_at"},
			)

			query := fmt.Sprintf(`
				INSERT INTO significant_notes (id, %s)
				VALUES (:id, %s)
			`, cols, vals)

			_, err := tx.NamedExecContext(
				ctx,
				query,
				sn,
			)
			if err != nil {
				return "", fmt.Errorf(
					"failed to create significant note: %w",
					err,
				)
			}

			return sn.ID, nil
		},
	)
}
