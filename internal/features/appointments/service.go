package appointments

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"

	"github.com/google/uuid"
)

type Service struct {
	repo         RepositoryInterface
	notifService audit.Notifier
	logService   audit.Logger
}

func NewService(
	repo RepositoryInterface,
	notifService audit.Notifier,
	logService audit.Logger,
) *Service {
	return &Service{
		repo:         repo,
		notifService: notifService,
		logService:   logService,
	}
}

func (s *Service) GetConcernCategories(
	ctx context.Context,
) ([]AppointmentCategory, error) {
	return s.repo.GetCategories(ctx)
}

func (s *Service) CreateAppointment(
	ctx context.Context,
	iirID string,
	req AppointmentDTO,
) (*Appointment, error) {
	appt := &Appointment{
		ID:                    uuid.New().String(),
		IIRID:                 iirID,
		Reason:                structs.ToSqlNull(req.Reason),
		WhenDate:              req.WhenDate,
		TimeSlotID:            req.TimeSlot.ID,
		AppointmentCategoryID: req.AppointmentCategory.ID,
		StatusID:              1,
	}

	err := datastore.RunInTransaction(
		ctx,
		s.repo.GetDB(),
		func(tx datastore.DB) error {
			return s.repo.CreateAppointment(ctx, tx, appt)
		},
	)
	if err != nil {
		audit.Dispatch(ctx, s.logService, s.notifService, audit.DispatchParams{
			Log: &audit.LogParams{
				Level:    audit.LevelError,
				Category: audit.CategoryAudit,
				Action:   audit.ActionAppointmentFailed,
				Message: fmt.Sprintf(
					"Failed to create appointment for IIR #%s",
					iirID,
				),
				Metadata: &audit.LogMetadata{
					EntityType: constants.AppointmentEntityType,
					NewValues:  req,
				},
			},
		})
		return nil, err
	}

	audit.Dispatch(ctx, s.logService, s.notifService, audit.DispatchParams{
		Log: &audit.LogParams{
			Level:    audit.LevelInfo,
			Category: audit.CategoryAudit,
			Action:   audit.ActionAppointmentCreated,
			Message:  fmt.Sprintf("Appointment #%s created", appt.ID),
			Metadata: &audit.LogMetadata{
				EntityType: constants.AppointmentEntityType,
				EntityID:   appt.ID,
				NewValues:  req,
			},
		},
		Notification: &audit.NotificationParams{
			TargetID: structs.StringToNullableString(appt.ID),
			TargetType: structs.StringToNullableString(
				constants.AppointmentEntityType,
			),
			Title:   fmt.Sprintf("Appointment #%s Created", appt.ID),
			Message: "Your appointment has been created and is pending approval.",
			Type:    constants.AppointmentEntityType,
		},
	})

	return appt, nil
}

func (s *Service) GetAppointmentByID(
	ctx context.Context,
	id string,
) (*Appointment, error) {
	return s.repo.GetAppointment(ctx, id)
}

func (s *Service) GetDailyStatusCount(
	ctx context.Context,
	startDate string,
) ([]DailyStatusCount, error) {
	layout := "2006-01-02"
	t, err := time.Parse(layout, startDate)
	if err != nil {
		return nil, err
	}

	startOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())

	endOfMonth := time.Date(
		t.Year(),
		t.Month()+2,
		0,
		23,
		59,
		59,
		0,
		t.Location(),
	)

	startStr := startOfMonth.Format(layout)
	endStr := endOfMonth.Format(layout)

	return s.repo.GetDailyStatusCount(ctx, startStr, endStr)
}

func (s *Service) ListAppointments(
	ctx context.Context,
	req ListAppointmentsRequest,
) (*ListAppointmentsDTO, error) {
	req.SetDefaults("created_at")

	statusIDs := []string{}
	if req.StatusID != "" {
		statusIDs = strings.Split(req.StatusID, ",")
	}

	appts, err := s.repo.List(
		ctx,
		req.GetOffset(),
		req.PageSize,
		req.Search,
		req.OrderBy,
		strings.Join(statusIDs, ","),
		req.StartDate,
		req.EndDate,
	)
	if err != nil {
		return nil, err
	}

	dtos := make([]AppointmentDTO, 0, len(appts))
	for _, appt := range appts {
		userDTO := users.GetUserResponse{
			ID:        "",
			Role:      users.Role{ID: 0, Name: ""},
			FirstName: appt.UserFirstName,
			MiddleName: structs.FromSqlNull(
				appt.UserMiddleName,
			),
			LastName: appt.UserLastName,
		}
		dtos = append(dtos, AppointmentDTO{
			ID:     appt.ID,
			User:   userDTO,
			Reason: structs.FromSqlNull(appt.Reason),
			AdminNotes: structs.FromSqlNull(
				appt.AdminNotes,
			),
			WhenDate: appt.WhenDate,
			TimeSlot: TimeSlot{
				ID:   appt.TimeSlotID,
				Time: appt.TimeSlotTime,
			},
			AppointmentCategory: AppointmentCategory{
				ID:   appt.CategoryID,
				Name: appt.CategoryName,
			},
			Status: AppointmentStatus{
				ID:       appt.StatusID,
				Name:     appt.StatusName,
				ColorKey: appt.StatusColorKey,
			},
			CreatedAt: appt.CreatedAt,
			UpdatedAt: appt.UpdatedAt,
		})
	}

	total, err := s.repo.GetTotalAppointmentsCount(
		ctx,
		req.StatusID,
		req.StartDate,
		req.EndDate,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &ListAppointmentsDTO{
		Appointments: dtos,
		Meta:         structs.CalculateMetadata(total, req.Page, req.PageSize),
	}, nil
}

func (s *Service) GetAppointmentsByUserID(
	ctx context.Context,
	userID string,
	req ListAppointmentsRequest,
) (*ListAppointmentsDTO, error) {
	req.SetDefaults("created_at")

	appts, err := s.repo.ListByUserID(
		ctx,
		userID,
		req.GetOffset(),
		req.PageSize,
		req.OrderBy,
		req.StatusID,
		req.StartDate,
		req.EndDate,
	)
	if err != nil {
		return nil, err
	}

	dtos := make([]AppointmentDTO, 0, len(appts))
	for _, appt := range appts {
		userDTO := users.GetUserResponse{
			Role:      users.Role{ID: 0, Name: ""},
			ID:        "",
			FirstName: appt.UserFirstName,
			MiddleName: structs.FromSqlNull(
				appt.UserMiddleName,
			),
			LastName: appt.UserLastName,
		}
		dtos = append(dtos, AppointmentDTO{
			ID:     appt.ID,
			User:   userDTO,
			Reason: structs.FromSqlNull(appt.Reason),
			AdminNotes: structs.FromSqlNull(
				appt.AdminNotes,
			),
			WhenDate: appt.WhenDate,
			TimeSlot: TimeSlot{
				ID:   appt.TimeSlotID,
				Time: appt.TimeSlotTime,
			},
			AppointmentCategory: AppointmentCategory{
				ID:   appt.CategoryID,
				Name: appt.CategoryName,
			},
			Status: AppointmentStatus{
				ID:       appt.StatusID,
				Name:     appt.StatusName,
				ColorKey: appt.StatusColorKey,
			},
			CreatedAt: appt.CreatedAt,
			UpdatedAt: appt.UpdatedAt,
		})
	}

	total, err := s.repo.GetTotalAppointmentsCount(
		ctx,
		req.StatusID,
		req.StartDate,
		req.EndDate,
		&userID,
	)
	if err != nil {
		return nil, err
	}

	return &ListAppointmentsDTO{
		Appointments: dtos,
		Meta:         structs.CalculateMetadata(total, req.Page, req.PageSize),
	}, nil
}

func (s *Service) GetAppointmentsByIIRID(
	ctx context.Context,
	iirID string,
	req ListAppointmentsRequest,
) (*ListAppointmentsDTO, error) {
	req.SetDefaults("created_at")

	appts, err := s.repo.ListByIIRID(
		ctx,
		iirID,
		req.GetOffset(),
		req.PageSize,
		req.OrderBy,
		req.StatusID,
		req.StartDate,
		req.EndDate,
	)
	if err != nil {
		return nil, err
	}

	dtos := make([]AppointmentDTO, 0, len(appts))
	for _, appt := range appts {
		userDTO := users.GetUserResponse{
			Role:      users.Role{ID: 0, Name: ""},
			ID:        "",
			FirstName: appt.UserFirstName,
			MiddleName: structs.FromSqlNull(
				appt.UserMiddleName,
			),
			LastName: appt.UserLastName,
		}
		dtos = append(dtos, AppointmentDTO{
			ID:     appt.ID,
			User:   userDTO,
			Reason: structs.FromSqlNull(appt.Reason),
			AdminNotes: structs.FromSqlNull(
				appt.AdminNotes,
			),
			WhenDate: appt.WhenDate,
			TimeSlot: TimeSlot{
				ID:   appt.TimeSlotID,
				Time: appt.TimeSlotTime,
			},
			AppointmentCategory: AppointmentCategory{
				ID:   appt.CategoryID,
				Name: appt.CategoryName,
			},
			Status: AppointmentStatus{
				ID:       appt.StatusID,
				Name:     appt.StatusName,
				ColorKey: appt.StatusColorKey,
			},
			CreatedAt: appt.CreatedAt,
			UpdatedAt: appt.UpdatedAt,
		})
	}

	total, err := s.repo.GetTotalAppointmentsCount(
		ctx,
		req.StatusID,
		req.StartDate,
		req.EndDate,
		&iirID,
	)
	if err != nil {
		return nil, err
	}

	return &ListAppointmentsDTO{
		Appointments: dtos,
		Meta:         structs.CalculateMetadata(total, req.Page, req.PageSize),
	}, nil
}

func (s *Service) GetAppointmentStats(
	ctx context.Context,
	req ListAppointmentsRequest,
	iirID *string,
) ([]StatusCount, error) {
	return s.repo.GetAppointmentStats(
		ctx,
		req.StatusID,
		req.StartDate,
		req.EndDate,
		iirID,
	)
}

func (s *Service) GetAvailableTimeSlots(
	ctx context.Context,
	date string,
) ([]AvailableTimeSlotView, error) {
	availableSlots, err := s.repo.GetAvailableTimeSlots(ctx, date)
	if err != nil {
		return nil, err
	}

	return availableSlots, nil
}

func (s *Service) GetAppointmentStatuses(
	ctx context.Context,
) ([]AppointmentStatus, error) {
	return s.repo.GetStatuses(ctx)
}

// handles Status updates AND Rescheduling
func (s *Service) UpdateAppointment(
	ctx context.Context,
	id string,
	req AppointmentDTO,
) error {
	// Fetch old state for audit trail
	oldAppt, _ := s.repo.GetAppointment(ctx, id)

	appt := Appointment{
		ID:                    id,
		StatusID:              req.Status.ID,
		Reason:                structs.ToSqlNull(req.Reason),
		WhenDate:              req.WhenDate,
		TimeSlotID:            req.TimeSlot.ID,
		AppointmentCategoryID: req.AppointmentCategory.ID,
	}

	err := datastore.RunInTransaction(
		ctx,
		s.repo.GetDB(),
		func(tx datastore.DB) error {
			return s.repo.UpdateAppointment(ctx, tx, appt)
		},
	)
	if err != nil {
		audit.Dispatch(ctx, s.logService, s.notifService, audit.DispatchParams{
			Log: &audit.LogParams{
				Level:    audit.LevelError,
				Category: audit.CategoryAudit,
				Action:   audit.ActionAppointmentUpdateFailed,
				Message:  fmt.Sprintf("Failed to update appointment #%s", id),
				Metadata: &audit.LogMetadata{
					EntityType: "appointment",
					EntityID:   id,
					OldValues:  oldAppt,
					NewValues:  req,
					Error:      err.Error(),
				},
			},
		})

		return err
	}

	audit.Dispatch(ctx, s.logService, s.notifService, audit.DispatchParams{
		Log: &audit.LogParams{
			Level:    audit.LevelInfo,
			Category: audit.CategoryAudit,
			Action:   audit.ActionAppointmentUpdated,
			Message:  fmt.Sprintf("Appointment #%s updated", id),
			Metadata: &audit.LogMetadata{
				EntityType: constants.AppointmentEntityType,
				EntityID:   id,
				OldValues:  oldAppt,
				NewValues:  req,
			},
		},
		Notification: &audit.NotificationParams{
			ReceiverID: structs.StringToNullableString(req.User.ID),
			TargetID:   structs.StringToNullableString(id),
			TargetType: structs.StringToNullableString(
				constants.AppointmentEntityType,
			),
			Title: fmt.Sprintf("Appointment #%s Updated", id),
			Message: fmt.Sprintf(
				"Your appointment has been updated. New status: %s",
				req.Status.Name,
			),
			Type: constants.AppointmentEntityType,
		},
	})

	return nil
}
