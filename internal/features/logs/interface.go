package logs

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
		req ListSystemLogsRequest,
	) (*ListSystemLogsDTO, error)
	GetStats(
		ctx context.Context,
		startDate, endDate string,
	) ([]LogStatsDTO, error)
}

type RepositoryInterface interface {
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
	) ([]LogStatsDTO, error)
}
