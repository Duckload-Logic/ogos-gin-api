package bootstrap

import (
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
	APIKeyService          *apikeys.Service
	NotificationsService   notifications.ServiceInterface
	SystemLogService       *logs.Service
	ConsentService         consents.ServiceInterface
}

func getServices(
	repos *Repositories,
	fileStorage storage.FileStorage,
	cfg *config.Config,
	redis *datastore.RedisClient,
) *Services {
	systemLogService := logs.NewService(repos.SystemLogRepo)
	apiKeyService := apikeys.NewService(repos.APIKeyRepo, systemLogService)
	notificationsService := notifications.NewService(repos.NotificationRepo)
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
	consentService := consents.NewService(
		repos.ConsentRepo,
		systemLogService,
		fileStorage,
	)

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
		ConsentService:         consentService,
	}
}
