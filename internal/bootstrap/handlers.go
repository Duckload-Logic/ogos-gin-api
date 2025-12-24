package bootstrap

import (
	"github.com/olazo-johnalbert/duckload-api/internal/features/appointments"
)

type Handlers struct {
	AppointmentHandler *appointments.Handler
}

func getHandlers(repos *Repositories) *Handlers {
	appointmentService := appointments.NewService(repos.AppointmentRepo)

	return &Handlers{
		AppointmentHandler: appointments.NewHandler(appointmentService),
	}
}
