package notes

import (
	"database/sql"
	"time"
)

// SignificantNoteDB represents the database model for professional notes.
type SignificantNoteDB struct {
	ID              string         `db:"id"`
	IIRID           sql.NullString `db:"iir_id"`
	AppointmentID   sql.NullString `db:"appointment_id"`
	AdmissionSlipID sql.NullString `db:"admission_slip_id"`
	Note            string         `db:"note"`
	Remarks         string         `db:"remarks"`
	CreatedAt       time.Time      `db:"created_at"`
	UpdatedAt       time.Time      `db:"updated_at"`
}
