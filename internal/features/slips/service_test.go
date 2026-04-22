package slips

import (
	"context"
	"testing"

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
	"go.uber.org/mock/gomock"
)

func TestService_GetSlipByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepositoryInterface(ctrl)
	mockLog := audit.NewMockLogger(ctrl)
	mockNotif := audit.NewMockNotifier(ctrl)
	svc := NewService(mockRepo, mockLog, mockNotif, nil, nil, nil, nil)
	ctx := context.Background()
	id := "slip-123"

	t.Run("success", func(t *testing.T) {
		view := &SlipWithDetailsView{
			ID:         id,
			StatusName: "Approved",
		}
		mockRepo.EXPECT().GetSlipByIDWithDetails(ctx, id).Return(view, nil)

		resp, err := svc.GetSlipByID(ctx, id)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.ID != id {
			t.Errorf("got %s, want %s", resp.ID, id)
		}
	})
}

func TestService_UpdateExcuseSlipStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepositoryInterface(ctrl)
	mockAuditLog := audit.NewMockLogger(ctrl)
	mockAuditNotif := audit.NewMockNotifier(ctrl)
	mockUserSvc := users.NewMockServiceInterface(ctrl)

	svc := NewService(
		mockRepo,
		mockAuditLog,
		mockAuditNotif,
		nil,
		mockUserSvc,
		nil,
		nil,
	)

	ctx := audit.WithContext(
		context.Background(),
		"127.0.0.1",
		"ua",
		"admin-1",
		"admin@email.com",
		"2",
		"trace-1",
	)
	id := "slip-123"

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().
			GetSlipByID(gomock.Any(), id).
			Return(&Slip{ID: id}, nil)
		mockRepo.EXPECT().
			WithTransaction(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(datastore.DB) error) error {
				return fn(nil)
			})
		mockRepo.EXPECT().
			UpdateStatus(gomock.Any(), gomock.Any(), id, "Approved", "Looks good").
			Return(nil)
		mockRepo.EXPECT().
			GetUserIDBySlipID(gomock.Any(), id).
			Return("student-1", nil)

		// Mock notifications and logs
		mockAuditLog.EXPECT().
			Record(gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()
		mockAuditNotif.EXPECT().Send(gomock.Any(), gomock.Any()).AnyTimes()

		err := svc.UpdateExcuseSlipStatus(ctx, id, "Approved", "Looks good")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
