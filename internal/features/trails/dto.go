package trails

import (
	"encoding/json"
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/request"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

// Action constants for audit trail entries
const (
	ActionCreate = "CREATE"
	ActionUpdate = "UPDATE"
	ActionDelete = "DELETE"
)

// AuditEntry is used by other features to record an audit trail entry.
// This is the reusable struct that services pass when logging actions.
type AuditEntry struct {
	UserID     int
	Action     string
	EntityType string
	EntityID   int
	OldValues  interface{}
	NewValues  interface{}
	IPAddress  string
	UserAgent  string
}

// ListAuditTrailsRequest holds query parameters for listing audit trails
type ListAuditTrailsRequest struct {
	request.PaginationParams
	Search     string `form:"search,omitempty"`
	Action     string `form:"action,omitempty" binding:"omitempty,oneof=CREATE UPDATE DELETE"`
	EntityType string `form:"entity_type,omitempty"`
	EntityID   int    `form:"entity_id,omitempty"`
	UserID     int    `form:"user_id,omitempty"`
	StartDate  string `form:"start_date,omitempty"`
	EndDate    string `form:"end_date,omitempty"`
	OrderBy    string `form:"order_by,omitempty" binding:"omitempty,oneof=created_at"`
}

// ListAuditTrailsDTO is the paginated response for audit trails
type ListAuditTrailsDTO struct {
	AuditTrails []AuditTrailDTO `json:"auditTrails"`
	Total       int             `json:"total"`
	Page        int             `json:"page"`
	PageSize    int             `json:"pageSize"`
	TotalPages  int             `json:"totalPages"`
}

// AuditTrailDTO is the response DTO for a single audit trail entry
type AuditTrailDTO struct {
	ID         int                    `json:"id"`
	UserID     int                    `json:"userId,omitempty"`
	UserName   structs.NullableString `json:"userName,omitempty"`
	Action     string                 `json:"action"`
	EntityType string                 `json:"entityType"`
	EntityID   int                    `json:"entityId"`
	OldValues  json.RawMessage        `json:"oldValues,omitempty"`
	NewValues  json.RawMessage        `json:"newValues,omitempty"`
	IPAddress  structs.NullableString `json:"ipAddress,omitempty"`
	UserAgent  structs.NullableString `json:"userAgent,omitempty"`
	CreatedAt  time.Time              `json:"createdAt"`
}
