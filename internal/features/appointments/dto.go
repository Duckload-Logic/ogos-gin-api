package appointments

import (
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/request"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
)

type ListAppointmentsRequest struct {
	request.PaginationParams
	Search    string `form:"search,omitempty"`
	StatusID  string `form:"status_id"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	OrderBy   string `form:"order_by" binding:"omitempty,oneof=created_at when_date"`
}

type ListAppointmentsDTO struct {
	Appointments []AppointmentDTO `json:"appointments"`
	Total        int              `json:"total"`
	Page         int              `json:"page"`
	PageSize     int              `json:"pageSize"`
	TotalPages   int              `json:"totalPages"`
}

type AppointmentDTO struct {
	ID                  int                    `db:"id" json:"id,omitempty"`
	User                users.GetUserResponse  `db:"user" json:"user,omitempty"`
	Reason              structs.NullableString `db:"reason" json:"reason,omitempty"`
	WhenDate            string                 `db:"when_date" json:"whenDate,omitempty"`
	TimeSlot            TimeSlot               `db:"time_slot" json:"timeSlot,omitempty"`
	AppointmentCategory AppointmentCategory    `db:"appointment_category" json:"appointmentCategory,omitempty"`
	AdminNotes          structs.NullableString `db:"admin_notes" json:"adminNotes,omitempty"`
	Status              AppointmentStatus      `db:"status" json:"status,omitempty"`
	CreatedAt           time.Time              `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt           time.Time              `db:"updated_at" json:"updatedAt,omitempty"`
}
