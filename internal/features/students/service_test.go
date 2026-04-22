package students

import (
	"context"
	"testing"

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"go.uber.org/mock/gomock"
)

func TestService_IsStudentLocked(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepositoryInterface(ctrl)
	mockLog := audit.NewMockLogger(ctrl)
	mockNotif := audit.NewMockNotifier(ctrl)

	svc := &Service{
		repo:         mockRepo,
		logService:   mockLog,
		notifService: mockNotif,
	}

	ctx := context.Background()
	iirID := "test-iir-id"

	t.Run("should return true when repo says locked", func(t *testing.T) {
		mockRepo.EXPECT().
			IsStudentLocked(gomock.Any(), iirID).
			Return(true, nil)

		locked, err := svc.IsStudentLocked(ctx, iirID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !locked {
			t.Error("expected true, got false")
		}
	})
}

func TestService_GetStudentBasicInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepositoryInterface(ctrl)
	svc := &Service{repo: mockRepo}

	ctx := context.Background()
	iirID := "test-iir-id"

	t.Run("success", func(t *testing.T) {
		info := &StudentBasicInfoView{
			Email:     "student@example.com",
			FirstName: "John",
			LastName:  "Doe",
		}

		mockRepo.EXPECT().GetStudentBasicInfo(ctx, iirID).Return(info, nil)

		resp, err := svc.GetStudentBasicInfo(ctx, iirID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Email != info.Email || resp.FirstName != info.FirstName {
			t.Errorf("resp mismatch: %+v", resp)
		}
	})
}

func TestService_ListStudents(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepositoryInterface(ctrl)
	svc := &Service{repo: mockRepo}

	ctx := context.Background()
	req := ListStudentsRequest{
		PaginationRequest: structs.PaginationRequest{
			Page:     1,
			PageSize: 10,
		},
	}

	t.Run("success empty list", func(t *testing.T) {
		mockRepo.EXPECT().
			ListStudents(ctx, gomock.Any(), 0, 10, gomock.Any(), 0, 0, 0, 0).
			Return([]StudentProfileView{}, nil)
		mockRepo.EXPECT().
			GetTotalStudentsCount(ctx, gomock.Any(), 0, 0, 0, 0).
			Return(0, nil)

		resp, err := svc.ListStudents(ctx, req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(resp.Students) != 0 {
			t.Error("expected empty students list")
		}
	})
}
