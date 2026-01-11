package appointments

import "time"

type Appointment struct {
	ID              int       `json:"id"`
	UserID          *int      `json:"userId,omitempty"`
	Reason          string    `json:"reason"`
	ScheduledDate   string    `json:"scheduledDate"`
	ScheduledTime   string    `json:"scheduledTime"`
	ConcernCategory string    `json:"concernCategory"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type CreateAppointmentRequest struct {
	Reason          string `json:"reason" binding:"required"`
	ScheduledDate   string `json:"scheduledDate" binding:"required"`
	ScheduledTime   string `json:"scheduledTime" binding:"required"`
	ConcernCategory string `json:"concernCategory,omitempty"`
}

type UpdateStatusRequest struct {
	Status          string `json:"status" binding:"required"`
	Reason          string `json:"reason,omitempty"`
	ScheduledDate   string `json:"scheduledDate,omitempty"`
	ScheduledTime   string `json:"scheduledTime,omitempty"`
	ConcernCategory string `json:"concernCategory,omitempty"`
}
