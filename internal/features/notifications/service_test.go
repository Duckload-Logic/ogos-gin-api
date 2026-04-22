package notifications

import (
	"context"
	"testing"

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"go.uber.org/mock/gomock"
)

func TestService_GetUserNotifications(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepositoryInterface(ctrl)
	svc := NewService(mockRepo)

	ctx := context.Background()
	userID := "test-user-id"

	t.Run("success", func(t *testing.T) {
		models := []Notification{
			{ID: "notif-1", Message: "Notif 1"},
		}

		mockRepo.EXPECT().GetByUserID(ctx, userID).Return(models, nil)

		resp, err := svc.GetUserNotifications(ctx, userID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resp) != 1 {
			t.Errorf("got %d, want 1", len(resp))
		}
	})
}

func TestService_Send(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepositoryInterface(ctrl)
	svc := NewService(mockRepo)

	ctx := context.Background()
	entry := audit.NotificationEntry{
		ReceiverID: structs.StringToNullableString("user-123"),
		Message:    "Hello",
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().Create(ctx, nil, gomock.Any()).Return(nil)

		err := svc.Send(ctx, entry)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
