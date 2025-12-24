package appointments

import (
	"context"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateAppointment(ctx context.Context, req CreateAppointmentRequest) (*Appointment, error) {

	parsedTime, err := time.Parse("2006-01-02 15:04:05", req.ScheduledAt)
	if err != nil {
		return nil, err
	}

	appt := &Appointment{
		StudentRecordID:   req.StudentRecordID,
		AppointmentTypeID: req.AppointmentTypeID,
		ScheduledAt:       parsedTime,
		ConcernCategory:   req.ConcernCategory,
		Status:            "Pending",
	}
	if err := s.repo.Create(ctx, appt); err != nil {
		return nil, err
	}

	return appt, nil
}
