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
	// Construct the model from the DTO
	appt := &Appointment{
		UserID:          &userID, // taking the address since your model uses *int
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

func (s *Service) ListAppointments(ctx context.Context, status string, startDate string, endDate string) ([]Appointment, error) {
	return s.repo.List(ctx, status, startDate, endDate)
}

func (s *Service) GetAppointmentsByUserID(ctx context.Context, userID int) ([]Appointment, error) {
	return s.repo.GetByUserID(ctx, userID)
}

// GetAvailableTimeSlots calculates which slots are free for a specific date
func (s *Service) GetAvailableTimeSlots(ctx context.Context, date string) (map[string]bool, error) {
	// Define the standard schedule
	allSlots := []string{
		"08:00:00", "09:00:00", "10:00:00", "11:00:00",
		"13:00:00", "14:00:00", "15:00:00", "16:00:00", "17:00:00",
	}

	// Fetch existing appointments for that date
	bookedAppts, err := s.repo.GetTimeSlots(ctx, date)
	if err != nil {
		return nil, err
	}

	// Initialize map: Assume all slots are TRUE (available)
	availableSlots := make(map[string]bool)
	for _, slot := range allSlots {
		availableSlots[slot] = true
	}

	// Mark booked slots as FALSE (unavailable)
	for _, appt := range bookedAppts {
		// Ensure strict matching of time strings
		if _, exists := availableSlots[appt.ScheduledTime]; exists {
			availableSlots[appt.ScheduledTime] = false
		}
	}

	return availableSlots, nil
}

// UpdateAppointmentStatus handles Status updates AND Rescheduling
func (s *Service) UpdateAppointmentStatus(ctx context.Context, id int, req UpdateStatusRequest) error {
	// Validate Status
	validStatuses := map[string]bool{
		"Pending":     true,
		"Approved":    true,
		"Rejected":    true,
		"Completed":   true,
		"Cancelled":   true,
		"Rescheduled": true,
	}

	if req.Status != "" && !validStatuses[req.Status] {
		return fmt.Errorf("invalid status provided: %s", req.Status)
	}

	// Call Repo (uses UpdateAppointment to handle potential date/time changes too)
	return s.repo.UpdateAppointment(ctx, id, req)
}