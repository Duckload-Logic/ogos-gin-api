package appointments

import "context"

type ServiceInterface interface {
	GetConcernCategories(ctx context.Context) ([]AppointmentCategory, error)
	CreateAppointment(
		ctx context.Context,
		iirID string,
		req AppointmentDTO,
	) (*Appointment, error)
	GetAppointmentByID(ctx context.Context, id string) (*Appointment, error)
	GetDailyStatusCount(
		ctx context.Context,
		startDate string,
	) ([]DailyStatusCount, error)
	ListAppointments(
		ctx context.Context,
		req ListAppointmentsRequest,
	) (*ListAppointmentsDTO, error)
	GetAppointmentsByUserID(
		ctx context.Context,
		userID string,
		req ListAppointmentsRequest,
	) (*ListAppointmentsDTO, error)
	GetAppointmentsByIIRID(
		ctx context.Context,
		iirID string,
		req ListAppointmentsRequest,
	) (*ListAppointmentsDTO, error)
	GetAppointmentStats(
		ctx context.Context,
		req ListAppointmentsRequest,
		iirID *string,
	) ([]StatusCount, error)
	GetAvailableTimeSlots(
		ctx context.Context,
		date string,
	) ([]AvailableTimeSlotView, error)
	GetAppointmentStatuses(ctx context.Context) ([]AppointmentStatus, error)
	UpdateAppointment(ctx context.Context, id string, req AppointmentDTO) error
	ConfirmAppointment(
		ctx context.Context,
		appointmentID string,
		studentEmail string,
	) error
}

type RepositoryInterface interface {
	GetTimeSlots(ctx context.Context, date string) ([]TimeSlot, error)
	GetCategories(ctx context.Context) ([]AppointmentCategory, error)
	GetAppointment(ctx context.Context, id string) (*Appointment, error)
	GetDailyStatusCount(
		ctx context.Context,
		startDate, endDate string,
	) ([]DailyStatusCount, error)
	GetTotalAppointmentsCount(
		ctx context.Context,
		statusID, startDate, endDate string,
		iirID *string,
	) (int, error)
	List(
		ctx context.Context,
		offset, limit int,
		search, orderBy, statusIDs, startDate, endDate string,
	) ([]AppointmentWithDetailsView, error)
	GetTimeSlotByID(ctx context.Context, id int) (*TimeSlot, error)
	GetAppointmentCategoryByID(
		ctx context.Context,
		id int,
	) (*AppointmentCategory, error)
	GetStatusByID(ctx context.Context, id int) (*AppointmentStatus, error)
	GetAvailableTimeSlots(
		ctx context.Context,
		date string,
	) ([]AvailableTimeSlotView, error)
	GetStatuses(ctx context.Context) ([]AppointmentStatus, error)
	ListByUserID(
		ctx context.Context,
		userID string,
		offset, limit int,
		orderBy string,
		statusID, startDate, endDate string,
	) ([]AppointmentWithDetailsView, error)
	ListByIIRID(
		ctx context.Context,
		iirID string,
		offset, limit int,
		orderBy string,
		statusID, startDate, endDate string,
	) ([]AppointmentWithDetailsView, error)
	GetAppointmentStats(
		ctx context.Context,
		statusID, startDate, endDate string,
		iirID *string,
	) ([]StatusCount, error)
	CreateAppointment(ctx context.Context, appt *Appointment) error
	UpdateAppointment(ctx context.Context, appt Appointment) error
}
