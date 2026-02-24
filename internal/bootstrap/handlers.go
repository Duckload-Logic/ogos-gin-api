package bootstrap

import (
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/features/appointments"
	"github.com/olazo-johnalbert/duckload-api/internal/features/auth"
	"github.com/olazo-johnalbert/duckload-api/internal/features/excuseslips"
	"github.com/olazo-johnalbert/duckload-api/internal/features/locations"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
)

type Handlers struct {
	DB                 *sqlx.DB
	AuthHandler        *auth.Handler
	UserHandler        *users.Handler
	LocationsHandler   *locations.Handler
	StudentHandler     *students.Handler
	AppointmentHandler *appointments.Handler
	ExcuseSlipHandler  *excuseslips.Handler
}

func getHandlers(repos *Repositories) *Handlers {
	authService := auth.NewService(repos.UserRepo, auth.NewTokenService())
	userService := users.NewService(repos.UserRepo)
	locationsService := locations.NewService(repos.LocationsRepo)
	studentService := students.NewService(repos.StudentRepo, locationsService)
	appointmentService := appointments.NewService(repos.AppointmentRepo)
	excuseSlipService := excuseslips.NewService(repos.ExcuseSlipRepo)

	return &Handlers{
		AuthHandler:        auth.NewHandler(authService),
		UserHandler:        users.NewHandler(userService),
		LocationsHandler:   locations.NewHandler(locationsService),
		StudentHandler:     students.NewHandler(studentService),
		AppointmentHandler: appointments.NewHandler(appointmentService),
		ExcuseSlipHandler:  excuseslips.NewHandler(excuseSlipService),
	}
}
