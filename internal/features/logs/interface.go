package logs

//go:generate mockgen -source=interface.go -package=logs -destination=mock_interfaces.go

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

// ServiceInterface defines the business logic for system logging.
type ServiceInterface interface {
	GetDB() datastore.DB
	Record(ctx context.Context, tx datastore.DB, entry audit.LogEntry)
	RecordSecurity(
		ctx context.Context,
		tx datastore.DB,
		action, message string,
		userEmail, userID, ipAddress, userAgent structs.NullableString,
	)
	ListLogs(
		ctx context.Context,
		req audit.ListSystemLogsRequest,
	) (*audit.ListSystemLogsDTO, error)
	GetStats(
		ctx context.Context,
		startDate, endDate string,
	) ([]audit.LogStatsDTO, error)
	GetActivityStats(ctx context.Context) ([]audit.LogActivityDTO, error)
	DeleteLogsOlderThan(ctx context.Context, days int) (int64, error)
}

type RepositoryInterface interface {
	WithTransaction(ctx context.Context, fn func(datastore.DB) error) error
	GetDB() *sqlx.DB
	Record(ctx context.Context, tx datastore.DB, log *SystemLog) error
	List(
		ctx context.Context, offset, limit int,
		category, action, userEmail, targetType, targetEmail,
		search, startDate, endDate, orderBy string,
	) ([]SystemLog, error)
	GetTotalCount(
		ctx context.Context,
		category, action, userEmail, targetType, targetEmail,
		search, startDate, endDate string,
	) (int, error)
	GetStats(
		ctx context.Context,
		startDate, endDate string,
	) ([]audit.LogStatsDTO, error)
	GetActivityStats(ctx context.Context) ([]audit.LogActivityDTO, error)
	DeleteLogsOlderThan(ctx context.Context, days int) (int64, error)
}
