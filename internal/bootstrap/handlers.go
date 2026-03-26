package bootstrap

import (
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/features/analytics"
	"github.com/olazo-johnalbert/duckload-api/internal/features/apikeys"
	"github.com/olazo-johnalbert/duckload-api/internal/features/appointments"
	"github.com/olazo-johnalbert/duckload-api/internal/features/auth"
	"github.com/olazo-johnalbert/duckload-api/internal/features/consents"
	"github.com/olazo-johnalbert/duckload-api/internal/features/locations"
	"github.com/olazo-johnalbert/duckload-api/internal/features/logs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/notes"
	"github.com/olazo-johnalbert/duckload-api/internal/features/notifications"
	"github.com/olazo-johnalbert/duckload-api/internal/features/slips"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students/external"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/storage"
)

type Handlers struct {
	DB                     *sqlx.DB
	AuthHandler            *auth.Handler
	UserHandler            *users.Handler
	LocationsHandler       *locations.Handler
	StudentHandler         *students.Handler
	NoteHandler            *notes.Handler
	ExternalStudentHandler *external.Handler
	AppointmentHandler     *appointments.Handler
	SlipHandler            *slips.Handler
	AnalyticsHandler       *analytics.Handler
	APIKeyHandler          *apikeys.Handler
	APIKeyService          *apikeys.Service
	NotificationsHandler   *notifications.Handler
	SystemLogHandler       *logs.Handler
	SystemLogService       *logs.Service
	ConsentHandler         *consents.Handler
	Redis                  *datastore.RedisClient
}

func getHandlers(
	repos *Repositories,
	fileStorage storage.FileStorage,
	cfg *config.Config,
	redis *datastore.RedisClient,
) *Handlers {
	systemLogService := logs.NewService(repos.SystemLogRepo)
	systemLogHandler := logs.NewHandler(systemLogService)
	apiKeyService := apikeys.NewService(repos.APIKeyRepo, systemLogService)

	notificationsService := notifications.NewService(repos.NotificationRepo)
	notificationsHandler := notifications.NewHandler(notificationsService)
	authService := auth.NewService(repos.UserRepo, redis)
	userService := users.NewService(repos.UserRepo)
	locationsService := locations.NewService(repos.LocationsRepo)
	studentService := students.NewService(repos.StudentRepo, locationsService)
	noteService := notes.NewService(repos.NoteRepo)
	externalStudentService := external.NewService(repos.ExternalStudentRepo)
	appointmentService := appointments.NewService(
		repos.AppointmentRepo,
		notificationsService,
		systemLogService,
	)
	slipService := slips.NewService(
		repos.SlipRepo,
		systemLogService,
		fileStorage,
	)
	analyticsService := analytics.NewService(repos.AnalyticsRepo)
	analyticsHandler := analytics.NewHandler(analyticsService)
	consentService := consents.NewService(
		repos.ConsentRepo,
		systemLogService,
		fileStorage,
	)

	return &Handlers{
		AuthHandler: auth.NewHandler(
			authService,
			systemLogService,
			cfg,
		),
		UserHandler:            users.NewHandler(userService),
		LocationsHandler:       locations.NewHandler(locationsService),
		StudentHandler:         students.NewHandler(studentService),
		NoteHandler:            notes.NewHandler(noteService),
		ExternalStudentHandler: external.NewHandler(externalStudentService),
		AppointmentHandler:     appointments.NewHandler(appointmentService),
		SlipHandler:            slips.NewHandler(slipService),
		AnalyticsHandler:       analyticsHandler,
		APIKeyHandler:          apikeys.NewHandler(apiKeyService),
		APIKeyService:          apiKeyService,
		NotificationsHandler:   notificationsHandler,
		SystemLogHandler:       systemLogHandler,
		SystemLogService:       systemLogService,
		ConsentHandler:         consents.NewHandler(consentService),
		Redis:                  redis,
	}
}
