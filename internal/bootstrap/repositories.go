package bootstrap

import (
	"database/sql"

	"github.com/olazo-johnalbert/duckload-api/internal/features/appointments"
)

type Repositories struct {
	AppointmentRepo *appointments.Repository
}

func getRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		AppointmentRepo: appointments.NewRepository(db),
	}
}
