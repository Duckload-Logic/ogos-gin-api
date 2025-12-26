package bootstrap

import (
	"database/sql"

	"github.com/olazo-johnalbert/duckload-api/internal/features/appointments"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students"
)

type Repositories struct {
	StudentRepo     *students.Repository
	AppointmentRepo *appointments.Repository
}

func getRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		StudentRepo:     students.NewRepository(db),
		AppointmentRepo: appointments.NewRepository(db),
	}
}
