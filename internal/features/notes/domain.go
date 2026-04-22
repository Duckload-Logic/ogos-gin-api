package notes

import (
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

// SignificantNote represents a pure business entity for professional notes.
type SignificantNote struct {
	ID              string
	IIRID           structs.NullableString
	AppointmentID   structs.NullableString
	AdmissionSlipID structs.NullableString
	Note            string
	Remarks         string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
