package bootstrap

import (
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/features/appointments"
	"github.com/olazo-johnalbert/duckload-api/internal/features/excuseslips"
	"github.com/olazo-johnalbert/duckload-api/internal/features/locations"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
)

type Repositories struct {
	UserRepo        *users.Repository
	StudentRepo     *students.Repository
	AppointmentRepo *appointments.Repository
	ExcuseSlipRepo  *excuseslips.Repository
	LocationsRepo   *locations.Repository
}

func getRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		UserRepo:        users.NewRepository(db),
		StudentRepo:     students.NewRepository(db),
		AppointmentRepo: appointments.NewRepository(db),
		ExcuseSlipRepo:  excuseslips.NewRepository(db),
		LocationsRepo:   locations.NewRepository(db),
	}
}
