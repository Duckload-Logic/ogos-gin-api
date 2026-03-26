package logs

import (
	"context"
)

// ServiceInterface defines the business logic for system logging.
type ServiceInterface interface {
	Record(ctx context.Context, entry LogEntry)
	RecordSecurity(
		ctx context.Context,
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

// RepositoryInterface defines the data access layer for system logging.
type RepositoryInterface interface {
	Record(ctx context.Context, log *SystemLog) error
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
