package bootstrap

import (
	"database/sql"

	"github.com/olazo-johnalbert/duckload-api/internal/features/appointments"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students"
	"github.com/olazo-johnalbert/duckload-api/internal/features/excuseslips"
)

type Repositories struct {
	StudentRepo     *students.Repository
	AppointmentRepo *appointments.Repository
	ExcuseSlipRepo  *excuseslips.Repository
}

func getRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		StudentRepo:     students.NewRepository(db),
		AppointmentRepo: appointments.NewRepository(db),
		ExcuseSlipRepo:  excuseslips.NewRepository(db),
	}
}
