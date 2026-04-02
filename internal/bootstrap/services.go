package bootstrap

import (
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/pdf"
	"github.com/olazo-johnalbert/duckload-api/internal/features/analytics"
	"github.com/olazo-johnalbert/duckload-api/internal/features/apikeys"
	"github.com/olazo-johnalbert/duckload-api/internal/features/appointments"
	"github.com/olazo-johnalbert/duckload-api/internal/features/auth"
	"github.com/olazo-johnalbert/duckload-api/internal/features/locations"
	"github.com/olazo-johnalbert/duckload-api/internal/features/logs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/notes"
	"github.com/olazo-johnalbert/duckload-api/internal/features/notifications"
	"github.com/olazo-johnalbert/duckload-api/internal/features/slips"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students/external"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/gotenberg"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/storage"
)

type Services struct {
	AuthService            auth.ServiceInterface
	UserService            users.ServiceInterface
	LocationsService       locations.ServiceInterface
	StudentService         students.ServiceInterface
	NoteService            notes.ServiceInterface
	ExternalStudentService external.ServiceInterface
	AppointmentService     appointments.ServiceInterface
	SlipService            slips.ServiceInterface
	AnalyticsService       analytics.ServiceInterface
	APIKeyService          apikeys.ServiceInterface
	NotificationsService   notifications.ServiceInterface
	SystemLogService       logs.ServiceInterface
}

func getServices(
	repos *Repositories,
	fileStorage storage.FileStorage,
	cfg *config.Config,
	redis *datastore.RedisClient,
) *Services {
	notificationsService := notifications.NewService(repos.NotificationRepo)
	userService := users.NewService(repos.UserRepo)
	systemLogService := logs.NewService(
		repos.SystemLogRepo,
		notificationsService,
		userService,
	)
	apiKeyService := apikeys.NewService(
		repos.APIKeyRepo,
		systemLogService,
		notificationsService,
	)
	authService := auth.NewService(repos.UserRepo, redis)
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
	externalStudentService := external.NewService(repos.ExternalStudentRepo)
	appointmentService := appointments.NewService(
		repos.AppointmentRepo,
		notificationsService,
		systemLogService,
		userService,
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
		AuthService:            authService,
		UserService:            userService,
		LocationsService:       locationsService,
		StudentService:         studentService,
		NoteService:            noteService,
		ExternalStudentService: externalStudentService,
		AppointmentService:     appointmentService,
		SlipService:            slipService,
		AnalyticsService:       analyticsService,
		APIKeyService:          apiKeyService,
		NotificationsService:   notificationsService,
		SystemLogService:       systemLogService,
	}
}
