package logs

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Record logs a system log entry. Fails silently (logs error)
// to avoid disrupting the parent operation.
//
// Usage:
//
//	logService.Record(ctx, logs.LogEntry{
//	    Category:  logs.CategorySecurity,
//	    Action:    logs.ActionLoginSuccess,
//	    Message:   "User john@example.com logged in successfully",
//	    UserEmail: "john@example.com",
//	    IPAddress: ipAddress,
//	    UserAgent: userAgent,
//	})
func (s *Service) Record(ctx context.Context, entry LogEntry) {
	sysLog := &SystemLog{
		Category:  entry.Category,
		Action:    entry.Action,
		Message:   entry.Message,
		UserID:    sql.NullInt64{Int64: int64(entry.UserID), Valid: entry.UserID != 0},
		UserEmail: sql.NullString{String: entry.UserEmail, Valid: entry.UserEmail != ""},
		IPAddress: sql.NullString{String: entry.IPAddress, Valid: entry.IPAddress != ""},
		UserAgent: sql.NullString{String: entry.UserAgent, Valid: entry.UserAgent != ""},
		Metadata:  toNullString(entry.Metadata),
	}

	if err := s.repo.Record(ctx, sysLog); err != nil {
		log.Printf("Failed to record system log: %v", err)
		return
	}
}

// RecordSecurity is a convenience method that satisfies the middleware.SecurityLogger interface.
// It records a security-category log entry with the given fields.
func (s *Service) RecordSecurity(ctx context.Context, userEmail, action, message, ipAddress, userAgent string, userID int) {
	s.Record(ctx, LogEntry{
		Category:  CategorySecurity,
		Action:    action,
		Message:   message,
		UserID:    userID,
		UserEmail: userEmail,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	})
}

// ListLogs returns a paginated list of system logs with filters
func (s *Service) ListLogs(ctx context.Context, req ListSystemLogsRequest) (*ListSystemLogsDTO, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}
	if req.OrderBy == "" {
		req.OrderBy = "created_at"
	}

	results, err := s.repo.List(
		ctx,
		req.GetOffset(), req.PageSize,
		req.Category, req.Action, req.UserEmail,
		req.Search, req.StartDate, req.EndDate, req.OrderBy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list system logs: %w", err)
	}

	dtos := s.mapLogsToDTOs(results)

	total, err := s.repo.GetTotalCount(
		ctx,
		req.Category, req.Action, req.UserEmail,
		req.Search, req.StartDate, req.EndDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to count system logs: %w", err)
	}

	return &ListSystemLogsDTO{
		Logs:       dtos,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: (total + req.PageSize - 1) / req.PageSize,
	}, nil
}

// GetStats returns log counts by category
func (s *Service) GetStats(ctx context.Context, startDate, endDate string) ([]LogStatsDTO, error) {
	return s.repo.GetStats(ctx, startDate, endDate)
}

func (s *Service) mapLogsToDTOs(logs []SystemLog) []SystemLogDTO {
	dtos := make([]SystemLogDTO, 0, len(logs))

	for _, l := range logs {
		dto := SystemLogDTO{
			ID:        l.ID,
			Category:  l.Category,
			Action:    l.Action,
			Message:   l.Message,
			UserID:    structs.FromSqlNullInt64(l.UserID),
			UserEmail: structs.FromSqlNull(l.UserEmail),
			IPAddress: structs.FromSqlNull(l.IPAddress),
			UserAgent: structs.FromSqlNull(l.UserAgent),
			CreatedAt: l.CreatedAt,
		}

		if l.Metadata.Valid {
			dto.Metadata = json.RawMessage(l.Metadata.String)
		}

		dtos = append(dtos, dto)
	}

	return dtos
}
