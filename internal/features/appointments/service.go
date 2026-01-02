package appointments

import (
    "fmt"
    "context"
    "time"
)

type Service struct {
    repo *Repository
}

func NewService(repo *Repository) *Service {
    return &Service{repo: repo}
}

// ==========================================
//              CREATE
// ==========================================

func (s *Service) CreateAppointment(ctx context.Context, req CreateAppointmentRequest) (*Appointment, error) {
    // Parse string date to time.Time
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

func (s *Service) GetAppointmentByID(ctx context.Context, id int) (*Appointment, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *Service) ListAppointments(ctx context.Context, status string, date string) ([]Appointment, error) {
    return s.repo.List(ctx, status, date)
}

func (s *Service) GetAppointmentsByStudentID(ctx context.Context, studentID int) ([]Appointment, error) {
    return s.repo.GetByStudentID(ctx, studentID)
}

func (s *Service) UpdateAppointmentStatus(ctx context.Context, id int, status string) error {
    validStatuses := map[string]bool{
        "Pending": true, "Approved": true, "Rejected": true, "Completed": true,
    }
    if !validStatuses[status] {
        return fmt.Errorf("invalid status")
    }
    return s.repo.UpdateStatus(ctx, id, status)
}