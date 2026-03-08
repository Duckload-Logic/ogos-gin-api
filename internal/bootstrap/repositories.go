package bootstrap

import (
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/features/analytics"
	"github.com/olazo-johnalbert/duckload-api/internal/features/apikeys"
	"github.com/olazo-johnalbert/duckload-api/internal/features/appointments"
	"github.com/olazo-johnalbert/duckload-api/internal/features/locations"
	"github.com/olazo-johnalbert/duckload-api/internal/features/notifications"
	"github.com/olazo-johnalbert/duckload-api/internal/features/slips"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students/external"
	"github.com/olazo-johnalbert/duckload-api/internal/features/trails"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
)

type Repositories struct {
	UserRepo            *users.Repository
	StudentRepo         *students.Repository
	ExternalStudentRepo *external.Repository
	AppointmentRepo     *appointments.Repository
	SlipRepo            *slips.Repository
	LocationsRepo       *locations.Repository
	AnalyticsRepo       *analytics.Repository
	AuditTrailRepo      *trails.Repository
	APIKeyRepo          *apikeys.Repository
	NotificationRepo    *notifications.Repository
}

func getRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		UserRepo:            users.NewRepository(db),
		StudentRepo:         students.NewRepository(db),
		ExternalStudentRepo: external.NewRepository(db),
		AppointmentRepo:     appointments.NewRepository(db),
		SlipRepo:            slips.NewRepository(db),
		LocationsRepo:       locations.NewRepository(db),
		AnalyticsRepo:       analytics.NewRepository(db),
		AuditTrailRepo:      trails.NewRepository(db),
		APIKeyRepo:          apikeys.NewRepository(db),
		NotificationRepo:    notifications.NewRepository(db),
	}
}
