package logs

import (
	"encoding/json"
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

// ListSystemLogsRequest holds query parameters for listing system logs
type ListSystemLogsRequest struct {
	structs.PaginationRequest
	Level       string `form:"level,omitempty"        binding:"omitempty,oneof=INFO WARNING ERROR CRITICAL"`
	Category    string `form:"category,omitempty"     binding:"omitempty,oneof=AUDIT SYSTEM SECURITY"`
	Action      string `form:"action,omitempty"`
	UserEmail   string `form:"user_email,omitempty"`
	TargetType  string `form:"target_type,omitempty"`
	TargetEmail string `form:"target_email,omitempty"`
	StartDate   string `form:"start_date,omitempty"`
	EndDate     string `form:"end_date,omitempty"`
}

// ListSystemLogsDTO is the paginated response for system logs
type ListSystemLogsDTO struct {
	Logs []SystemLogDTO             `json:"logs"`
	Meta structs.PaginationMetadata `json:"meta"`
}

// SystemLogDTO is the response DTO for a single system log entry
type SystemLogDTO struct {
	ID       int    `json:"id"`
	Level    string `json:"level"`
	Category string `json:"category"`
	Action   string `json:"action"`
	Message  string `json:"message"`

	UserID     structs.NullableString `json:"userId,omitempty"`
	TargetID   structs.NullableString `json:"targetId,omitempty"`
	TargetType structs.NullableString `json:"targetType,omitempty"`

	UserEmail   structs.NullableString `json:"userEmail,omitempty"`
	TargetEmail structs.NullableString `json:"targetEmail,omitempty"`

	IPAddress structs.NullableString `json:"ipAddress,omitempty"`
	UserAgent structs.NullableString `json:"userAgent,omitempty"`
	Metadata  json.RawMessage        `json:"metadata,omitempty"`
	TraceID   structs.NullableString `json:"traceId,omitempty"`
	CreatedAt time.Time              `json:"createdAt"`
}

// LogStatsDTO returns summary counts by category
type LogStatsDTO struct {
	Category string `json:"category" db:"category"`
	Count    int    `json:"count"    db:"count"`
}

type LogActivityDTO struct {
	Time     string `json:"time"     db:"time"`
	Requests int    `json:"requests" db:"requests"`
	Errors   int    `json:"errors"   db:"errors"`
}
