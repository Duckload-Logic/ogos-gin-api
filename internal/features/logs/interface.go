package logs

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

// ServiceInterface defines the business logic for system logging.
type ServiceInterface interface {
	GetDB() datastore.DB
	Record(ctx context.Context, tx datastore.DB, entry LogEntry)
	RecordSecurity(
		ctx context.Context,
		tx datastore.DB,
		userEmail, action, message, ipAddress, userAgent string,
		userID string,
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
		category, action, userEmail, search, startDate, endDate, orderBy string,
	) ([]SystemLog, error)
	GetTotalCount(
		ctx context.Context,
		category, action, userEmail, search, startDate, endDate string,
	) (int, error)
	GetStats(
		ctx context.Context,
		startDate, endDate string,
	) ([]LogStatsDTO, error)
}
