package logs

import (
	"context"
	"testing"

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"go.uber.org/mock/gomock"
)

func TestService_Record(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepositoryInterface(ctrl)
	mockNotif := audit.NewMockNotifier(ctrl)
	mockGetter := audit.NewMockUserGetter(ctrl)
	svc := NewService(mockRepo, mockNotif, mockGetter)

	ctx := context.Background()
	entry := audit.LogEntry{
		Category: "test",
		Action:   "create",
		Message:  "test message",
	}

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().Record(ctx, nil, gomock.Any()).Return(nil)

		svc.Record(ctx, nil, entry)
	})
}

func TestService_RecordSecurity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepositoryInterface(ctrl)
	mockNotif := audit.NewMockNotifier(ctrl)
	mockGetter := audit.NewMockUserGetter(ctrl)
	svc := NewService(mockRepo, mockNotif, mockGetter)

	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().Record(ctx, nil, gomock.Any()).Return(nil)

		svc.RecordSecurity(
			ctx,
			nil,
			"login",
			"success",
			structs.StringToNullableString("test@example.com"),
			structs.StringToNullableString("123"),
			structs.NullableString{Valid: false},
			structs.NullableString{Valid: false},
		)
	})
}

func TestService_ListLogs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepositoryInterface(ctrl)
	mockNotif := audit.NewMockNotifier(ctrl)
	mockGetter := audit.NewMockUserGetter(ctrl)
	svc := NewService(mockRepo, mockNotif, mockGetter)

	ctx := context.Background()
	req := audit.ListSystemLogsRequest{
		PaginationRequest: structs.PaginationRequest{
			Page:     1,
			PageSize: 10,
		},
	}

	t.Run("success search", func(t *testing.T) {
		systemLogs := []SystemLog{
			{ID: 1, Message: "Log 1"},
		}
		mockRepo.EXPECT().
			List(ctx, 0, 10, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(systemLogs, nil)
		mockRepo.EXPECT().
			GetTotalCount(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(1, nil)

		resp, err := svc.ListLogs(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resp.Logs) != 1 {
			t.Errorf("got %d, want 1", len(resp.Logs))
		}
	})
}
