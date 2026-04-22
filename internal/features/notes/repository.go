package notes

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) RepositoryInterface {
	return &Repository{db: db}
}

func (r *Repository) WithTransaction(
	ctx context.Context,
	fn func(datastore.DB) error,
) error {
	return datastore.RunInTransaction(ctx, r.db, fn)
}

func (r *Repository) GetStudentSignificantNotes(
	ctx context.Context,
	iirID string,
) ([]SignificantNote, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM significant_notes
		WHERE iir_id = ?
	`, datastore.GetColumns(SignificantNoteDB{}))

	var notes []SignificantNoteDB
	err := r.db.SelectContext(ctx, &notes, query, iirID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get student significant notes: %w",
			err,
		)
	}

	fmt.Printf(
		"[GetStudentSignificantNotes] {Query}: Retrieved %d notes for %s\n",
		len(notes),
		iirID,
	)

	return MapSignificantNotesToDomain(notes), nil
}

func (r *Repository) CreateSignificantNote(
	ctx context.Context,
	sn *SignificantNote,
) (string, error) {
	return datastore.NewRunInTransaction(
		ctx,
		r.db,
		func(tx datastore.DB) (string, error) {
			dbModel := MapSignificantNoteToDB(*sn)
			cols, vals := datastore.GetInsertStatement(
				SignificantNoteDB{},
				[]string{"created_at", "updated_at"},
			)

			query := fmt.Sprintf(`
				INSERT INTO significant_notes (id, %s)
				VALUES (:id, %s)
			`, cols, vals)

			_, err := tx.NamedExecContext(
				ctx,
				query,
				dbModel,
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

func (r *Repository) HasNoteForAppointment(
	ctx context.Context,
	appointmentID string,
) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM significant_notes WHERE appointment_id = ?"
	err := r.db.GetContext(ctx, &count, query, appointmentID)
	if err != nil {
		return false, fmt.Errorf(
			"failed to check for existing note: %w",
			err,
		)
	}

	return count > 0, nil
}
