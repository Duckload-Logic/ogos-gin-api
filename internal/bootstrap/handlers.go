package bootstrap

import (
	"github.com/olazo-johnalbert/duckload-api/internal/features/appointments"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students"
)

type Handlers struct {
	StudentHandler     *students.Handler
	AppointmentHandler *appointments.Handler
}

func getHandlers(repos *Repositories) *Handlers {
	studentService := students.NewService(repos.StudentRepo)
	appointmentService := appointments.NewService(repos.AppointmentRepo)

	return &Handlers{
		StudentHandler:     students.NewHandler(studentService),
		AppointmentHandler: appointments.NewHandler(appointmentService),
	}
}
