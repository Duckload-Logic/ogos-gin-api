package appointments

import (
	"context"
	"strings"
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetConcernCategories(ctx context.Context) ([]AppointmentCategory, error) {
	return s.repo.GetCategories(ctx)
}

func (s *Service) CreateAppointment(ctx context.Context, userID int, req AppointmentDTO) (*Appointment, error) {
	appt := &Appointment{
		UserID:                userID,
		Reason:                structs.ToSqlNull(req.Reason),
		WhenDate:              req.WhenDate,
		TimeSlotID:            req.TimeSlot.ID,
		AppointmentCategoryID: req.AppointmentCategory.ID,
		StatusID:              1,
	}

	if err := s.repo.CreateAppointment(ctx, appt); err != nil {
		return nil, err
	}

	return appt, nil
}

func (s *Service) GetAppointmentByID(ctx context.Context, id int) (*Appointment, error) {
	return s.repo.GetAppointment(ctx, id)
}

func (s *Service) GetDailyStatusCount(ctx context.Context, startDate string) ([]DailyStatusCount, error) {
	layout := "2006-01-02"
	t, err := time.Parse(layout, startDate)
	if err != nil {
		return nil, err
	}

	startOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())

	endOfMonth := time.Date(t.Year(), t.Month()+2, 0, 23, 59, 59, 0, t.Location())

	startStr := startOfMonth.Format(layout)
	endStr := endOfMonth.Format(layout)

	return s.repo.GetDailyStatusCount(ctx, startStr, endStr)
}

func (s *Service) ListAppointments(ctx context.Context, req ListAppointmentsRequest) (*ListAppointmentsDTO, error) {
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	if req.OrderBy == "" {
		req.OrderBy = "created_at"
	}

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
			ID:         appt.UserID,
			Role:       users.Role{ID: 0, Name: ""},
			FirstName:  appt.UserFirstName,
			MiddleName: structs.FromSqlNull(appt.UserMiddleName),
			LastName:   appt.UserLastName,
			Email:      appt.UserEmail,
		}
		dtos = append(dtos, AppointmentDTO{
			ID:         appt.ID,
			User:       userDTO,
			Reason:     structs.FromSqlNull(appt.Reason),
			AdminNotes: structs.FromSqlNull(appt.AdminNotes),
			WhenDate:   appt.WhenDate,
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

	total, err := s.repo.GetTotalAppointmentsCount(ctx, req.StatusID, req.StartDate, req.EndDate, nil)
	if err != nil {
		return nil, err
	}

	return &ListAppointmentsDTO{
		Appointments: dtos,
		Total:        total,
		Page:         req.Page,
		PageSize:     req.PageSize,
		TotalPages:   (total + req.PageSize - 1) / req.PageSize,
	}, nil
}

func (s *Service) GetAppointmentsByUserID(ctx context.Context, userID int, req ListAppointmentsRequest) (*ListAppointmentsDTO, error) {
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	if req.OrderBy == "" {
		req.OrderBy = "created_at"
	}

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
			ID:         appt.UserID,
			Role:       users.Role{ID: 0, Name: ""},
			FirstName:  appt.UserFirstName,
			MiddleName: structs.FromSqlNull(appt.UserMiddleName),
			LastName:   appt.UserLastName,
			Email:      appt.UserEmail,
		}
		dtos = append(dtos, AppointmentDTO{
			ID:         appt.ID,
			User:       userDTO,
			Reason:     structs.FromSqlNull(appt.Reason),
			AdminNotes: structs.FromSqlNull(appt.AdminNotes),
			WhenDate:   appt.WhenDate,
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

	total, err := s.repo.GetTotalAppointmentsCount(ctx, req.StatusID, req.StartDate, req.EndDate, &userID)
	if err != nil {
		return nil, err
	}

	return &ListAppointmentsDTO{
		Appointments: dtos,
		Total:        total,
		Page:         req.Page,
		PageSize:     req.PageSize,
		TotalPages:   (total + req.PageSize - 1) / req.PageSize,
	}, nil
}

func (s *Service) GetAppointmentStats(ctx context.Context, req ListAppointmentsRequest, userID *int) ([]StatusCount, error) {
	return s.repo.GetAppointmentStats(ctx, req.StatusID, req.StartDate, req.EndDate, userID)
}

func (s *Service) GetAvailableTimeSlots(ctx context.Context, date string) ([]AvailableTimeSlotView, error) {
	availableSlots, err := s.repo.GetAvailableTimeSlots(ctx, date)
	if err != nil {
		return nil, err
	}

	return availableSlots, nil
}

func (s *Service) GetAppointmentStatuses(ctx context.Context) ([]AppointmentStatus, error) {
	return s.repo.GetStatuses(ctx)
}

// handles Status updates AND Rescheduling
func (s *Service) UpdateAppointment(ctx context.Context, id int, req AppointmentDTO) error {
	appt := Appointment{
		ID:                    id,
		Reason:                structs.ToSqlNull(req.Reason),
		WhenDate:              req.WhenDate,
		TimeSlotID:            req.TimeSlot.ID,
		AppointmentCategoryID: req.AppointmentCategory.ID,
	}

	err := s.repo.UpdateAppointment(ctx, appt)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateAppointmentStatus(ctx context.Context, id int, req AppointmentDTO) error {
	appt := Appointment{
		ID:                    id,
		StatusID:              req.Status.ID,
		Reason:                structs.ToSqlNull(req.Reason),
		AdminNotes:            structs.ToSqlNull(req.AdminNotes),
		WhenDate:              req.WhenDate,
		TimeSlotID:            req.TimeSlot.ID,
		AppointmentCategoryID: req.AppointmentCategory.ID,
	}

	err := s.repo.UpdateAppointment(ctx, appt)
	if err != nil {
		return err
	}

	return nil
}
