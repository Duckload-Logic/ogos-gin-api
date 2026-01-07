package appointments

import (
	"context"
	"fmt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateAppointment(ctx context.Context, userID int, req CreateAppointmentRequest) (*Appointment, error) {
	appt := &Appointment{
		UserID:          &userID,
		Reason:          req.Reason,
		ScheduledTime:   req.ScheduledTime,
		ScheduledDate:   req.ScheduledDate,
		ConcernCategory: req.ConcernCategory,
		Status:          "Pending",
	}
	if err := s.repo.Create(ctx, appt); err != nil {
		return nil, err
	}

	return appt, nil
}

func (s *Service) GetAppointmentByID(ctx context.Context, id int) (*Appointment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListAppointments(ctx context.Context, status string, date string) ([]Appointment, error) {
	return s.repo.List(ctx, status, date)
}

func (s *Service) GetAppointmentsByUserID(ctx context.Context, userID int) ([]Appointment, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *Service) GetAvailableTimeSlots(ctx context.Context, date string) (map[string]bool, error) {
	return s.getAvailableSlots(ctx, date)
}

func (s *Service) UpdateAppointmentStatus(ctx context.Context, id int, status string) error {
	validStatuses := map[string]bool{
		"Pending": true, "Approved": true, "Rejected": true, "Completed": true, "Cancelled": true, "Rescheduled": true,
	}
	if !validStatuses[status] {
		return fmt.Errorf("invalid status")
	}
	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *Service) getAvailableSlots(ctx context.Context, date string) (map[string]bool, error) {
	// Define all possible time slots
	allSlots := []string{
		"08:00:00",
		"09:00:00", "10:00:00",
		"11:00:00", "13:00:00",
		"14:00:00", "15:00:00",
		"16:00:00", "17:00:00",
	}

	// Get booked appointments for the given date
	bookedAppts, err := s.repo.GetTimeSlots(ctx, date)
	if err != nil {
		return nil, err
	}

	// Create a map to track available slots
	availableSlots := make(map[string]bool)
	for _, slot := range allSlots {
		availableSlots[slot] = true
	}

	// Mark booked slots as unavailable
	for _, appt := range bookedAppts {
		slotTime := appt.ScheduledTime
		if _, exists := availableSlots[slotTime]; exists {
			availableSlots[slotTime] = false
		}
	}

	return availableSlots, nil
}
