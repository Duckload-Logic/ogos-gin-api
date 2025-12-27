package bootstrap

import (
	"github.com/olazo-johnalbert/duckload-api/internal/features/appointments"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students"
	"github.com/olazo-johnalbert/duckload-api/internal/features/excuseslips"
)

type Handlers struct {
	StudentHandler     *students.Handler
	AppointmentHandler *appointments.Handler
	ExcuseSlipHandler  *excuseslips.Handler
}

func getHandlers(repos *Repositories) *Handlers {
	studentService := students.NewService(repos.StudentRepo)
	appointmentService := appointments.NewService(repos.AppointmentRepo)
	excuseSlipService := excuseslips.NewService(repos.ExcuseSlipRepo)

	return &Handlers{
		StudentHandler:     students.NewHandler(studentService),
		AppointmentHandler: appointments.NewHandler(appointmentService),
		ExcuseSlipHandler:  excuseslips.NewHandler(excuseSlipService),
	}
}
