package appointments

import (
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
)

type ListAppointmentsRequest struct {
	structs.PaginationRequest
	StatusID  string `form:"status_id"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
}

type ListAppointmentsDTO struct {
	Appointments []AppointmentDTO           `json:"appointments"`
	Meta         structs.PaginationMetadata `json:"meta"`
}

type AppointmentDTO struct {
	ID                  string                 `db:"id"                   json:"id,omitempty"`
	User                users.GetUserResponse  `db:"user"                 json:"user,omitempty"`
	IIRID               string                 `db:"iir_id"               json:"iirId,omitempty"`
	StudentNumber       string                 `db:"student_number"       json:"studentNumber,omitempty"`
	Reason              structs.NullableString `db:"reason"               json:"reason,omitempty"`
	WhenDate            string                 `db:"when_date"            json:"whenDate,omitempty"`
	TimeSlot            TimeSlot               `db:"time_slot"            json:"timeSlot,omitempty"`
	AppointmentCategory AppointmentCategory    `db:"appointment_category" json:"appointmentCategory,omitempty"`
	AdminNotes          structs.NullableString `db:"admin_notes"          json:"adminNotes,omitempty"`
	Status              AppointmentStatus      `db:"status"               json:"status,omitempty"`
	HasSignificantNote  bool                   `                          json:"hasSignificantNote"`
	CreatedAt           time.Time              `db:"created_at"           json:"createdAt,omitempty"`
	UpdatedAt           time.Time              `db:"updated_at"           json:"updatedAt,omitempty"`
}
