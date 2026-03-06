package trails

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

// Record logs an audit trail entry. This is the primary method used by
// other features to record CREATE, UPDATE, and DELETE actions.
//
// Usage from other services:
//
//	auditService.Record(ctx, trails.AuditEntry{
//	    UserEmail:  userEmail,
//	    Action:     trails.ActionCreate,
//	    EntityType: "appointment",
//	    EntityID:   appointmentID,
//	    NewValues:  appointmentData,
//	    IPAddress:  ipAddress,
//	    UserAgent:  userAgent,
//	})
func (s *Service) Record(ctx context.Context, entry AuditEntry) {
	trail := &AuditTrail{
		UserEmail:  sql.NullString{String: entry.UserEmail, Valid: entry.UserEmail != ""},
		Action:     entry.Action,
		EntityType: entry.EntityType,
		EntityID:   entry.EntityID,
		OldValues:  toNullString(entry.OldValues),
		NewValues:  toNullString(entry.NewValues),
		IPAddress:  sql.NullString{String: entry.IPAddress, Valid: entry.IPAddress != ""},
		UserAgent:  sql.NullString{String: entry.UserAgent, Valid: entry.UserAgent != ""},
	}

	if err := s.repo.Record(ctx, trail); err != nil {
		// Log the error but don't fail the parent operation
		log.Printf("Failed to record audit trail: %v", err)
	}
}

// ListAuditTrails returns a paginated list of audit trails with filters
func (s *Service) ListAuditTrails(ctx context.Context, req ListAuditTrailsRequest) (*ListAuditTrailsDTO, error) {
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

	trails, err := s.repo.List(
		ctx,
		req.GetOffset(), req.PageSize,
		req.Action, req.EntityType, req.EntityID, req.UserEmail,
		req.Search, req.StartDate, req.EndDate, req.OrderBy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list audit trails: %w", err)
	}

	dtos := s.mapTrailsToDTOs(trails)

	total, err := s.repo.GetTotalCount(
		ctx,
		req.Action, req.EntityType, req.EntityID, req.UserEmail,
		req.Search, req.StartDate, req.EndDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to count audit trails: %w", err)
	}

	return &ListAuditTrailsDTO{
		AuditTrails: dtos,
		Total:       total,
		Page:        req.Page,
		PageSize:    req.PageSize,
		TotalPages:  (total + req.PageSize - 1) / req.PageSize,
	}, nil
}

// GetByEntity returns all audit trail entries for a specific entity
func (s *Service) GetByEntity(ctx context.Context, entityType string, entityID int) ([]AuditTrailDTO, error) {
	trails, err := s.repo.GetByEntityRef(ctx, entityType, entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit trails: %w", err)
	}

	return s.mapTrailsToDTOs(trails), nil
}

func (s *Service) mapTrailsToDTOs(trails []AuditTrailWithUserView) []AuditTrailDTO {
	dtos := make([]AuditTrailDTO, 0, len(trails))

	for _, t := range trails {
		dto := AuditTrailDTO{
			ID:         t.ID,
			Action:     t.Action,
			EntityType: t.EntityType,
			EntityID:   t.EntityID,
			IPAddress:  structs.FromSqlNull(t.IPAddress),
			UserAgent:  structs.FromSqlNull(t.UserAgent),
			CreatedAt:  t.CreatedAt,
		}

		if t.UserEmail.Valid {
			dto.UserEmail = t.UserEmail.String
		}

		// Build user display name
		if t.UserFirstName.Valid && t.UserLastName.Valid {
			name := t.UserFirstName.String
			if t.UserMiddleName.Valid {
				name += " " + t.UserMiddleName.String
			}
			name += " " + t.UserLastName.String
			dto.UserName = structs.NullableString{String: name, Valid: true}
		}

		if t.OldValues.Valid {
			dto.OldValues = json.RawMessage(t.OldValues.String)
		}
		if t.NewValues.Valid {
			dto.NewValues = json.RawMessage(t.NewValues.String)
		}

		dtos = append(dtos, dto)
	}

	return dtos
}
