package appointments

import (
	"context"
	"fmt"
	"testing"

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/notes"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/ai/classifier"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
	"go.uber.org/mock/gomock"
)

func TestService_CreateAppointment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepositoryInterface(ctrl)
	mockAuditLog := audit.NewMockLogger(ctrl)
	mockAuditNotif := audit.NewMockNotifier(ctrl)
	mockUserSvc := users.NewMockServiceInterface(ctrl)
	mockNoteSvc := notes.NewMockServiceInterface(ctrl)
	mockStudentSvc := students.NewMockServiceInterface(ctrl)
	mockClassifier := classifier.NewMockServiceInterface(ctrl)

	cfg := config.NewTestConfig()

	svc := &Service{
		repo:           mockRepo,
		logService:     mockAuditLog,
		notifService:   mockAuditNotif,
		userService:    mockUserSvc,
		noteService:    mockNoteSvc,
		studentService: mockStudentSvc,
		classifier:     mockClassifier,
	}

	ctx := audit.WithContext(
		context.Background(),
		"127.0.0.1",
		"test-ua",
		"user-1",
		"test@email.com",
		"1",
		"trace-1",
	)
	iirID := "test-iir-id"
	req := AppointmentDTO{
		Reason:              structs.StringToNullableString("I feel anxious"),
		WhenDate:            "2026-04-20T14:00:00Z",
		TimeSlot:            TimeSlot{ID: 1},
		AppointmentCategory: AppointmentCategory{ID: 1},
	}

	t.Run("should fail if student record is locked", func(t *testing.T) {
		mockStudentSvc.EXPECT().
			IsStudentLocked(gomock.Any(), iirID).
			Return(true, nil)

		appt, err := svc.CreateAppointment(ctx, iirID, req, cfg)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if appt != nil {
			t.Error("expected nil appointment, got one")
		}
		if err.Error() != "cannot create appointment: student record is locked (Graduated/Archived)" {
			t.Errorf("unexpected error message: %s", err.Error())
		}
	})

	t.Run("should succeed when valid data is provided", func(t *testing.T) {
		mockStudentSvc.EXPECT().
			IsStudentLocked(gomock.Any(), iirID).
			Return(false, nil)

		mockClassifier.EXPECT().
			Classify(gomock.Any(), req.Reason.String, cfg).
			Return(&classifier.ClassifyResponse{
				Level:      "HIGH",
				Confidence: 0.95,
			}, nil)

		mockRepo.EXPECT().
			WithTransaction(gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(ctx context.Context, fn func(datastore.DB) error) error {
					return fn(nil)
				},
			)

		mockRepo.EXPECT().
			CreateAppointment(gomock.Any(), nil, gomock.Any()).
			Return(nil)

		mockUserSvc.EXPECT().
			GetUserByID(gomock.Any(), "user-1").
			Return(nil, nil).
			AnyTimes()
		mockUserSvc.EXPECT().
			GetUserIDsByRole(gomock.Any(), gomock.Any()).
			Return([]string{"admin-1"}, nil).
			AnyTimes()

		mockAuditLog.EXPECT().
			Record(gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes()
		mockAuditNotif.EXPECT().Send(gomock.Any(), gomock.Any()).AnyTimes()

		appt, err := svc.CreateAppointment(ctx, iirID, req, cfg)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if appt == nil {
			t.Fatal("expected appointment, got nil")
		}
		if appt.UrgencyLevel != "HIGH" {
			t.Errorf("expected HIGH urgency, got %s", appt.UrgencyLevel)
		}
	})

	t.Run(
		"should handle classifier failure by falling back to MEDIUM",
		func(t *testing.T) {
			mockStudentSvc.EXPECT().
				IsStudentLocked(gomock.Any(), iirID).
				Return(false, nil)

			mockClassifier.EXPECT().
				Classify(gomock.Any(), req.Reason.String, cfg).
				Return(nil, fmt.Errorf("classifier down"))

			mockRepo.EXPECT().
				WithTransaction(gomock.Any(), gomock.Any()).
				DoAndReturn(
					func(ctx context.Context, fn func(datastore.DB) error) error {
						return fn(nil)
					},
				)

			mockRepo.EXPECT().
				CreateAppointment(gomock.Any(), nil, gomock.Any()).
				Return(nil)

			mockUserSvc.EXPECT().
				GetUserByID(gomock.Any(), gomock.Any()).
				Return(nil, nil).
				AnyTimes()
			mockUserSvc.EXPECT().
				GetUserIDsByRole(gomock.Any(), gomock.Any()).
				Return(nil, nil).
				AnyTimes()
			mockAuditLog.EXPECT().
				Record(gomock.Any(), gomock.Any(), gomock.Any()).
				AnyTimes()
			mockAuditNotif.EXPECT().Send(gomock.Any(), gomock.Any()).AnyTimes()

			appt, err := svc.CreateAppointment(ctx, iirID, req, cfg)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if appt.UrgencyLevel != "MEDIUM" {
				t.Errorf(
					"expected MEDIUM urgency on fallback, got %s",
					appt.UrgencyLevel,
				)
			}
		},
	)
}

func TestService_GetAppointmentByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepositoryInterface(ctrl)
	mockNoteSvc := notes.NewMockServiceInterface(ctrl)
	svc := &Service{repo: mockRepo, noteService: mockNoteSvc}

	ctx := context.Background()
	id := "apt-123"

	t.Run("success", func(t *testing.T) {
		appt := &AppointmentWithDetailsView{
			ID:            id,
			UserFirstName: "John",
			UserLastName:  "Doe",
			StatusName:    "Pending",
		}
		mockRepo.EXPECT().GetAppointment(ctx, id).Return(appt, nil)
		mockNoteSvc.EXPECT().HasNoteForAppointment(ctx, id).Return(true, nil)

		resp, err := svc.GetAppointmentByID(ctx, id)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.ID != id || !resp.HasSignificantNote {
			t.Errorf("resp mismatch: %+v", resp)
		}
	})
}

func TestService_UpdateAppointment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepositoryInterface(ctrl)
	mockAuditLog := audit.NewMockLogger(ctrl)
	mockAuditNotif := audit.NewMockNotifier(ctrl)
	mockNoteSvc := notes.NewMockServiceInterface(ctrl)

	svc := &Service{
		repo:         mockRepo,
		logService:   mockAuditLog,
		notifService: mockAuditNotif,
		noteService:  mockNoteSvc,
	}

	ctx := audit.WithContext(
		context.Background(),
		"127.0.0.1",
		"test-ua",
		"admin-1",
		"admin@email.com",
		"2",
		"trace-1",
	)
	id := "apt-123"
	req := AppointmentDTO{
		Status: AppointmentStatus{ID: 3}, // Completed
	}

	t.Run(
		"success completed without note should trigger notification",
		func(t *testing.T) {
			oldAppt := &AppointmentWithDetailsView{ID: id}
			newAppt := &AppointmentWithDetailsView{
				ID:           id,
				StatusName:   "Completed",
				WhenDate:     "2026-04-20",
				TimeSlotTime: "09:00",
			}

			mockRepo.EXPECT().GetAppointment(ctx, id).Return(oldAppt, nil)
			mockRepo.EXPECT().
				WithTransaction(ctx, gomock.Any()).
				DoAndReturn(
					func(ctx context.Context, fn func(datastore.DB) error) error {
						return fn(nil)
					},
				)
			mockRepo.EXPECT().
				UpdateAppointment(ctx, nil, gomock.Any()).
				Return(nil)
			mockRepo.EXPECT().GetAppointment(ctx, id).Return(newAppt, nil)
			mockRepo.EXPECT().
				GetUserIDByAppointmentID(ctx, id).
				Return("student-1", nil)

			mockNoteSvc.EXPECT().
				HasNoteForAppointment(ctx, id).
				Return(false, nil)

			// Expectations for audit dispatch
			mockAuditLog.EXPECT().
				Record(gomock.Any(), gomock.Any(), gomock.Any()).
				AnyTimes()
			mockAuditNotif.EXPECT().Send(gomock.Any(), gomock.Any()).AnyTimes()

			err := svc.UpdateAppointment(ctx, id, req)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		},
	)
}
