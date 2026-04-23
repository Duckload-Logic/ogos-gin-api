package m2mclients

import (
	"context"
	"testing"

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
	"go.uber.org/mock/gomock"
)

func TestService_ListClients(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepositoryInterface(ctrl)
	svc := NewService(mockRepo, nil, nil, nil, nil)

	ctx := context.Background()
	userID := "user-123"

	t.Run("success", func(t *testing.T) {
		clients := []M2MClient{
			{ID: 1, ClientName: "Client 1", ClientID: "cid-1"},
		}
		mockRepo.EXPECT().List(ctx, userID, false).Return(clients, nil)

		resp, err := svc.ListClients(ctx, userID, false, []int{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resp) != 1 {
			t.Errorf("got %d, want 1", len(resp))
		}
		if resp[0].ClientName != "Client 1" {
			t.Errorf("got %s, want Client 1", resp[0].ClientName)
		}
	})
}

func TestService_VerifyClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepositoryInterface(ctrl)
	mockAuditLog := audit.NewMockLogger(ctrl)
	mockAuditNotif := audit.NewMockNotifier(ctrl)

	svc := NewService(mockRepo, mockAuditLog, mockAuditNotif, nil, nil)

	ctx := context.Background()
	id := 1

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().
			WithTransaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(datastore.DB) error) error {
				return fn(nil)
			})
		mockRepo.EXPECT().
			UpdateVerificationStatus(ctx, nil, id, true).
			Return(nil)

		mockRepo.EXPECT().
			GetByID(ctx, nil, id).
			Return(&M2MClient{ID: id, ClientName: "Test", UserID: "user-1"}, nil)

		// Mock audit record
		mockAuditLog.EXPECT().
			Record(gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()
		mockAuditNotif.EXPECT().Send(gomock.Any(), gomock.Any()).AnyTimes()

		err := svc.VerifyClient(ctx, id)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
