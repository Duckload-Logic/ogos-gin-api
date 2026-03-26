package notes

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type Service struct {
	repo RepositoryInterface
}

func NewService(repo RepositoryInterface) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetStudentSignificantNotes(
	ctx context.Context,
	iirID string,
) ([]SignificantNoteDTO, error) {
	notes, err := s.repo.GetStudentSignificantNotes(
		ctx,
		iirID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get student significant notes: %w",
			err,
		)
	}

	var noteDTOs []SignificantNoteDTO
	for _, n := range notes {
		noteDTOs = append(noteDTOs, SignificantNoteDTO{
			ID:        n.ID,
			Note:      n.Note,
			Remarks:   n.Remarks,
			CreatedAt: n.CreatedAt,
			UpdatedAt: n.UpdatedAt,
		})
	}

	return noteDTOs, nil
}

func (s *Service) CreateSignificantNote(
	ctx context.Context,
	iirID string,
	noteReq SignificantNoteDTO,
) error {
	note := &SignificantNote{
		ID:      uuid.New().String(),
		IIRID:   sql.NullString{String: iirID, Valid: true},
		Note:    noteReq.Note,
		Remarks: noteReq.Remarks,
	}

	_, err := s.repo.CreateSignificantNote(ctx, note)
	if err != nil {
		return fmt.Errorf(
			"failed to create significant note: %w",
			err,
		)
	}

	return nil
}
