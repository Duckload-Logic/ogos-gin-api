package logs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Service struct {
	repo RepositoryInterface
}

func NewService(repo RepositoryInterface) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetDB() datastore.DB {
	return s.repo.GetDB()
}

// Record logs a system log entry. Fails silently (logs error)
// to avoid disrupting the parent operation.
func (s *Service) Record(
	ctx context.Context,
	tx datastore.DB,
	entry audit.LogEntry,
) {
	level := entry.Level
	if level == "" {
		level = audit.LevelInfo
	}

	// Safety net: extract trace ID from context if not provided
	if !entry.TraceID.Valid || entry.TraceID.String == "" {
		_, _, _, _, trace := audit.ExtractMeta(ctx)
		if trace != "" {
			entry.TraceID = structs.StringToNullableString(trace)
		}
	}

	sysLog := &SystemLog{
		Level:       level,
		Category:    entry.Category,
		Action:      entry.Action,
		Message:     entry.Message,
		UserID:      structs.ToSqlNull(entry.UserID),
		TargetID:    structs.ToSqlNull(entry.TargetID),
		TargetType:  structs.ToSqlNull(entry.TargetType),
		UserEmail:   structs.ToSqlNull(entry.UserEmail),
		TargetEmail: structs.ToSqlNull(entry.TargetEmail),
		IPAddress:   structs.ToSqlNull(entry.IPAddress),
		UserAgent:   structs.ToSqlNull(entry.UserAgent),
		TraceID:     structs.ToSqlNull(entry.TraceID),
		Metadata:    toNullString(entry.Metadata),
	}

	if err := s.repo.Record(ctx, tx, sysLog); err != nil {
		log.Printf("Failed to record system log: %v", err)
		return
	}
}

// RecordSecurity is a convenience method that satisfies the
// middleware.SecurityLogger interface.
// It records a security-category log entry with the given fields.
func (s *Service) RecordSecurity(
	ctx context.Context,
	tx datastore.DB,
	action, message string,
	userEmail, userID, ipAddress, userAgent structs.NullableString,
) {
	s.Record(ctx, tx, audit.LogEntry{
		Category:  audit.CategorySecurity,
		Action:    action,
		Message:   message,
		UserID:    userID,
		UserEmail: userEmail,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	})
}

// ListLogs returns a paginated list of system logs with filters
func (s *Service) ListLogs(
	ctx context.Context,
	req ListSystemLogsRequest,
) (*ListSystemLogsDTO, error) {
	req.SetDefaults("created_at")

	results, err := s.repo.List(
		ctx,
		req.GetOffset(), req.PageSize,
		req.Category, req.Action, req.UserEmail,
		req.TargetType, req.TargetEmail,
		req.Search, req.StartDate, req.EndDate, req.OrderBy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list system logs: %w", err)
	}

	dtos := s.mapLogsToDTOs(results)

	total, err := s.repo.GetTotalCount(
		ctx,
		req.Category, req.Action, req.UserEmail,
		req.TargetType, req.TargetEmail,
		req.Search, req.StartDate, req.EndDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to count system logs: %w", err)
	}

	return &ListSystemLogsDTO{
		Logs: dtos,
		Meta: structs.CalculateMetadata(total, req.Page, req.PageSize),
	}, nil
}

// GetStats returns log counts by category
func (s *Service) GetStats(
	ctx context.Context,
	startDate, endDate string,
) ([]LogStatsDTO, error) {
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
			UserID:    structs.FromSqlNull(l.UserID),
			UserEmail: structs.FromSqlNull(l.UserEmail),
			IPAddress: structs.FromSqlNull(l.IPAddress),
			UserAgent: structs.FromSqlNull(l.UserAgent),
			TraceID:   structs.FromSqlNull(l.TraceID),
			CreatedAt: l.CreatedAt,
		}

		if l.Metadata.Valid {
			dto.Metadata = json.RawMessage(l.Metadata.String)
		}

		dtos = append(dtos, dto)
	}

	return dtos
}
