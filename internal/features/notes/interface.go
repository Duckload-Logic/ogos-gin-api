package notes

import (
	"context"

	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

// ServiceInterface defines the business logic for managing student notes.
type ServiceInterface interface {
	GetStudentSignificantNotes(
		ctx context.Context,
		iirID string,
	) ([]SignificantNoteDTO, error)
	CreateSignificantNote(
		ctx context.Context,
		iirID string,
		noteReq SignificantNoteDTO,
	) error
	HasNoteForAppointment(
		ctx context.Context,
		appointmentID string,
	) (bool, error)
}

// RepositoryInterface defines the data access layer for managing student notes.
type RepositoryInterface interface {
	GetStudentSignificantNotes(
		ctx context.Context,
		iirID string,
	) ([]SignificantNote, error)
	CreateSignificantNote(
		ctx context.Context,
		sn *SignificantNote,
	) (string, error)
	HasNoteForAppointment(
		ctx context.Context,
		appointmentID string,
	) (bool, error)
	WithTransaction(ctx context.Context, fn func(datastore.DB) error) error
}
