package notes

import (
	"context"
	"testing"

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"go.uber.org/mock/gomock"
)

func TestService_GetStudentSignificantNotes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepositoryInterface(ctrl)
	mockLog := audit.NewMockLogger(ctrl)
	mockNotif := audit.NewMockNotifier(ctrl)
	svc := NewService(mockRepo, mockLog, mockNotif)

	ctx := context.Background()
	iirID := "test-iir-id"

	t.Run("success", func(t *testing.T) {
		notes := []SignificantNote{
			{ID: "note-1", Note: "Note 1"},
		}

		mockRepo.EXPECT().
			GetStudentSignificantNotes(ctx, iirID).
			Return(notes, nil)

		resp, err := svc.GetStudentSignificantNotes(ctx, iirID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resp) != 1 {
			t.Errorf("got %d, want 1", len(resp))
		}
	})
}

func TestService_HasNoteForAppointment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepositoryInterface(ctrl)
	mockLog := audit.NewMockLogger(ctrl)
	mockNotif := audit.NewMockNotifier(ctrl)
	svc := NewService(mockRepo, mockLog, mockNotif)

	ctx := context.Background()
	appointmentID := "apt-123"

	t.Run("true", func(t *testing.T) {
		mockRepo.EXPECT().
			HasNoteForAppointment(ctx, appointmentID).
			Return(true, nil)

		has, err := svc.HasNoteForAppointment(ctx, appointmentID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !has {
			t.Error("expected true")
		}
	})
}
