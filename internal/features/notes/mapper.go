package notes

import (
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

// MapSignificantNoteToDomain converts DB model to domain model.
func MapSignificantNoteToDomain(db SignificantNoteDB) SignificantNote {
	return SignificantNote{
		ID:              db.ID,
		IIRID:           structs.FromSqlNull(db.IIRID),
		AppointmentID:   structs.FromSqlNull(db.AppointmentID),
		AdmissionSlipID: structs.FromSqlNull(db.AdmissionSlipID),
		Note:            db.Note,
		Remarks:         db.Remarks,
		CreatedAt:       db.CreatedAt,
		UpdatedAt:       db.UpdatedAt,
	}
}

// MapSignificantNoteToDB converts domain model to DB model.
func MapSignificantNoteToDB(d SignificantNote) SignificantNoteDB {
	return SignificantNoteDB{
		ID:              d.ID,
		IIRID:           structs.ToSqlNull(d.IIRID),
		AppointmentID:   structs.ToSqlNull(d.AppointmentID),
		AdmissionSlipID: structs.ToSqlNull(d.AdmissionSlipID),
		Note:            d.Note,
		Remarks:         d.Remarks,
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}
}

// MapSignificantNotesToDomain maps a slice of DB models to domain models.
func MapSignificantNotesToDomain(db []SignificantNoteDB) []SignificantNote {
	domain := make([]SignificantNote, len(db))
	for i := range db {
		domain[i] = MapSignificantNoteToDomain(db[i])
	}
	return domain
}
