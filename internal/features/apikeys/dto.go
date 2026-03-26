package apikeys

import "time"

// Requests

type CreateAPIKeyRequest struct {
	Name      string   `json:"name"                binding:"required"`
	Scopes    []string `json:"scopes,omitempty"`
	ExpiresAt *string  `json:"expiresAt,omitempty"` // RFC3339 format, optional
}

type ListAPIKeysRequest struct {
	IncludeRevoked bool `form:"include_revoked,omitempty"`
}

// Responses

type APIKeyDTO struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	KeyPrefix  string     `json:"keyPrefix"`
	Scopes     []string   `json:"scopes,omitempty"`
	IsActive   bool       `json:"isActive"`
	LastUsedAt *time.Time `json:"lastUsedAt,omitempty"`
	ExpiresAt  *time.Time `json:"expiresAt,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
}

// CreateAPIKeyResponse includes the full key — only shown once at creation time.
type CreateAPIKeyResponse struct {
	APIKeyDTO
	Key string `json:"key"`
}
