package appointments

import (
	"time"
)

type Appointment struct {
	ID                int       `json:"id" db:"id"`
	StudentRecordID   int       `json:"studentRecordId" db:"student_record_id"`
	CounselorUserID   *int      `json:"counselorUserId" db:"counselor_user_id"`
	AppointmentTypeID int       `json:"appointmentTypeId" db:"appointment_type_id"`
	ConcernCategory   string    `json:"concernCategory" db:"concern_category"`
	Status            string    `json:"status" db:"status"`
	ScheduledAt       time.Time `json:"scheduledAt" db:"scheduled_at"`
	CreatedAt         time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt         time.Time `json:"updatedAt" db:"updated_at"`
}

type CreateAppointmentRequest struct {
	StudentRecordID   int    `json:"studentRecordId" binding:"required"`
	AppointmentTypeID int    `json:"appointmentTypeId" binding:"required"`
	ScheduledAt       string `json:"scheduledAt" binding:"required"`
	ConcernCategory   string `json:"concernCategory" binding:"required"`
}