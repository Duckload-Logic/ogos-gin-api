package logs

import (
	"encoding/json"
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

// Log categories
const (
	CategoryAudit    = "AUDIT"
	CategorySystem   = "SYSTEM"
	CategorySecurity = "SECURITY"
	CategoryConsent  = "CONSENT"
)

// Audit log actions — track data changes and admin operations
const (
	ActionUserCreated          = "USER_CREATED"
	ActionUserUpdated          = "USER_UPDATED"
	ActionUserDeleted          = "USER_DELETED"
	ActionRoleChanged          = "ROLE_CHANGED"
	ActionAppointmentCreated   = "APPOINTMENT_CREATED"
	ActionAppointmentUpdated   = "APPOINTMENT_UPDATED"
	ActionSlipCreated          = "SLIP_CREATED"
	ActionSlipStatusUpdated    = "SLIP_STATUS_UPDATED"
	ActionStudentRecordCreated = "STUDENT_RECORD_CREATED"
	ActionStudentRecordUpdated = "STUDENT_RECORD_UPDATED"
)

// System log actions — track system-level events
const (
	ActionAPIKeyCreated  = "API_KEY_CREATED"
	ActionAPIKeyRevoked  = "API_KEY_REVOKED"
	ActionSettingChanged = "SETTING_CHANGED"
)

// Security log actions — track authentication and access events
const (
	ActionLoginSuccess      = "LOGIN_SUCCESS"
	ActionLoginFailed       = "LOGIN_FAILED"
	ActionLogout            = "LOGOUT"
	ActionTokenRefreshed    = "TOKEN_REFRESHED"
	ActionAccessDenied      = "ACCESS_DENIED"
	ActionRateLimitExceeded = "RATE_LIMIT_EXCEEDED"
	ActionInvalidToken      = "INVALID_TOKEN"
	ActionAPIKeyUsed        = "API_KEY_USED"
	ActionAPIKeyInvalid     = "API_KEY_INVALID"
)

// Consent log actions — track user agreements and legal updates
const (
	// User-side events
	ActionTermsAccepted   = "TERMS_ACCEPTED"
	ActionPrivacyAccepted = "PRIVACY_ACCEPTED"

	// Admin-side/System events
	ActionPolicyUpdated  = "POLICY_UPDATED"  // When we upload a new version to Azure
	ActionConsentRevoked = "CONSENT_REVOKED" // If a user deletes their account or opts out
)

// LogEntry is the input struct used by other services to record a log.
// It now uses UserID and TargetID (integers) instead of email strings.
type LogEntry struct {
	Category  string
	Action    string
	Message   string
	UserID    string // "" means no user (system event)
	UserEmail string
	IPAddress string
	UserAgent string
	Metadata  interface{}
}

// ListSystemLogsRequest holds query parameters for listing system logs
type ListSystemLogsRequest struct {
	structs.PaginationRequest
	Category  string `form:"category,omitempty"   binding:"omitempty,oneof=AUDIT SYSTEM SECURITY"`
	Action    string `form:"action,omitempty"`
	UserEmail string `form:"user_email,omitempty"`
	StartDate string `form:"start_date,omitempty"`
	EndDate   string `form:"end_date,omitempty"`
}

// ListSystemLogsDTO is the paginated response for system logs
type ListSystemLogsDTO struct {
	Logs []SystemLogDTO             `json:"logs"`
	Meta structs.PaginationMetadata `json:"meta"`
}

// SystemLogDTO is the response DTO for a single system log entry
type SystemLogDTO struct {
	ID        int                    `json:"id"`
	Category  string                 `json:"category"`
	Action    string                 `json:"action"`
	Message   string                 `json:"message"`
	UserID    structs.NullableString `json:"userId,omitempty"`
	UserEmail structs.NullableString `json:"userEmail,omitempty"`
	IPAddress structs.NullableString `json:"ipAddress,omitempty"`
	UserAgent structs.NullableString `json:"userAgent,omitempty"`
	Metadata  json.RawMessage        `json:"metadata,omitempty"`
	CreatedAt time.Time              `json:"createdAt"`
}

// LogStatsDTO returns summary counts by category
type LogStatsDTO struct {
	Category string `json:"category" db:"category"`
	Count    int    `json:"count"    db:"count"`
}
