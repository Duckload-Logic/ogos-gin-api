package notes

import (
	"database/sql"
	"time"
)

// Significant Notes and Incidents
type SignificantNote struct {
	ID              string         `db:"id" json:"id"`
	IIRID           sql.NullString `db:"iir_id" json:"iirId"`
	AppointmentID   sql.NullString `db:"appointment_id" json:"appointmentId"`
	AdmissionSlipID sql.NullString `db:"admission_slip_id" json:"admissionSlipId"`
	Note            string         `db:"note" json:"note"`
	Remarks         string         `db:"remarks" json:"remarks"`
	CreatedAt       time.Time      `db:"created_at" json:"createdAt"`
	UpdatedAt       time.Time      `db:"updated_at" json:"updatedAt"`
}
