package bootstrap

import (
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/pdf"
	"github.com/olazo-johnalbert/duckload-api/internal/core/sessions"
	"github.com/olazo-johnalbert/duckload-api/internal/core/tokens"
	"github.com/olazo-johnalbert/duckload-api/internal/features/analytics"
	"github.com/olazo-johnalbert/duckload-api/internal/features/appointments"
	"github.com/olazo-johnalbert/duckload-api/internal/features/auth"
	"github.com/olazo-johnalbert/duckload-api/internal/features/locations"
	"github.com/olazo-johnalbert/duckload-api/internal/features/logs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/m2mclients"
	"github.com/olazo-johnalbert/duckload-api/internal/features/notes"
	"github.com/olazo-johnalbert/duckload-api/internal/features/notifications"
	"github.com/olazo-johnalbert/duckload-api/internal/features/slips"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students/integrations"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/email"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/gotenberg"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/storage"
)

type Services struct {
	AuthService               auth.ServiceInterface
	UserService               users.ServiceInterface
	LocationsService          locations.ServiceInterface
	StudentService            students.ServiceInterface
	NoteService               notes.ServiceInterface
	IntegrationStudentService integrations.ServiceInterface
	AppointmentService        appointments.ServiceInterface
	SlipService               slips.ServiceInterface
	AnalyticsService          analytics.ServiceInterface
	M2MClientService          m2mclients.ServiceInterface
	NotificationsService      notifications.ServiceInterface
	SystemLogService          logs.ServiceInterface
	SessionService            *sessions.Service
}

func getServices(
	repos *Repositories,
	fileStorage storage.FileStorage,
	cfg *config.Config,
	redis *datastore.RedisClient,
	emailer email.Emailer,
) *Services {
	notificationsService := notifications.NewService(repos.NotificationRepo)
	userService := users.NewService(repos.UserRepo)
	systemLogService := logs.NewService(
		repos.SystemLogRepo,
		notificationsService,
		userService,
	)
	tokenService := tokens.NewService()
	sessionService := sessions.NewService(redis)
	m2mClientService := m2mclients.NewService(
		repos.M2MClientRepo,
		systemLogService,
		notificationsService,
		tokenService,
		sessionService,
	)
	authService := auth.NewService(
		repos.UserRepo,
		redis,
		sessionService,
		emailer,
	)
	locationsService := locations.NewService(repos.LocationsRepo)

	gotenbergClient := gotenberg.NewClient(cfg.GotenbergURL)
	pdfService := pdf.NewService(gotenbergClient)

	studentService := students.NewService(
		repos.StudentRepo,
		locationsService,
		userService,
		systemLogService,
		notificationsService,
		cfg,
		pdfService,
	)
	noteService := notes.NewService(
		repos.NoteRepo,
		systemLogService,
		notificationsService,
	)
	integrationStudentService := integrations.NewService(
		repos.IntegrationStudentRepo,
	)
	appointmentService := appointments.NewService(
		repos.AppointmentRepo,
		notificationsService,
		systemLogService,
		userService,
		noteService,
	)
	slipService := slips.NewService(
		repos.SlipRepo,
		systemLogService,
		notificationsService,
		fileStorage,
		userService,
	)
	analyticsService := analytics.NewService(repos.AnalyticsRepo)

	return &Services{
		AuthService:               authService,
		UserService:               userService,
		LocationsService:          locationsService,
		StudentService:            studentService,
		NoteService:               noteService,
		IntegrationStudentService: integrationStudentService,
		AppointmentService:        appointmentService,
		SlipService:               slipService,
		AnalyticsService:          analyticsService,
		M2MClientService:          m2mClientService,
		NotificationsService:      notificationsService,
		SystemLogService:          systemLogService,
		SessionService:            sessionService,
	}
}
