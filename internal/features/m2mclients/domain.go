package m2mclients

import (
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

// M2MClient represents a pure business entity for a Machine-to-Machine client.
type M2MClient struct {
	ID                int
	UserID            string
	ClientName        string
	ClientID          string
	ClientSecretHash  string
	ClientDescription string
	Scopes            structs.NullableString
	IsActive          bool
	IsVerified        bool
	LastUsedAt        structs.NullableTime
	ExpiresAt         structs.NullableTime
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
