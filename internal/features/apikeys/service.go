package apikeys

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/features/logs"
)

type Service struct {
	repo       *Repository
	logService *logs.Service
}

func NewService(repo *Repository, logService ...*logs.Service) *Service {
	s := &Service{repo: repo}
	if len(logService) > 0 {
		s.logService = logService[0]
	}
	return s
}

// GenerateKey creates a new API key, stores its hash, and returns the plaintext key (shown once).
func (s *Service) GenerateKey(ctx context.Context, req CreateAPIKeyRequest) (*CreateAPIKeyResponse, error) {
	// Generate a random 32-byte key
	rawKey := make([]byte, 32)
	if _, err := rand.Read(rawKey); err != nil {
		return nil, fmt.Errorf("failed to generate random key: %w", err)
	}

	plaintext := "ogos_" + hex.EncodeToString(rawKey) // 69-char key with prefix
	prefix := plaintext[:8]

	hash := sha256.Sum256([]byte(plaintext))
	keyHash := hex.EncodeToString(hash[:])

	var scopesJSON sql.NullString
	if len(req.Scopes) > 0 {
		b, _ := json.Marshal(req.Scopes)
		scopesJSON = sql.NullString{String: string(b), Valid: true}
	}

	var expiresAt sql.NullTime
	if req.ExpiresAt != nil {
		t, err := time.Parse(time.RFC3339, *req.ExpiresAt)
		if err != nil {
			return nil, fmt.Errorf("invalid expiresAt format (use RFC3339): %w", err)
		}
		expiresAt = sql.NullTime{Time: t, Valid: true}
	}

	apiKey := APIKey{
		Name:      req.Name,
		KeyHash:   keyHash,
		KeyPrefix: prefix,
		Scopes:    scopesJSON,
		IsActive:  true,
		ExpiresAt: expiresAt,
	}

	id, err := s.repo.Create(ctx, apiKey)
	if err != nil {
		return nil, err
	}

	dto := mapKeyToDTO(apiKey)
	dto.ID = id

	// Record system log
	if s.logService != nil {
		s.logService.Record(ctx, logs.LogEntry{
			Category: logs.CategorySystem,
			Action:   logs.ActionAPIKeyCreated,
			Message:  fmt.Sprintf("API key '%s' created (prefix: %s)", req.Name, prefix),
			Metadata: map[string]interface{}{
				"keyName":   req.Name,
				"keyPrefix": prefix,
				"scopes":    req.Scopes,
			},
		})
	}

	return &CreateAPIKeyResponse{
		APIKeyDTO: dto,
		Key:       plaintext,
	}, nil
}

// ValidateKey checks if a plaintext key is valid, active, and not expired.
// Returns the APIKey record if valid, or an error.
func (s *Service) ValidateKey(ctx context.Context, plaintext string) (*APIKey, error) {
	hash := sha256.Sum256([]byte(plaintext))
	keyHash := hex.EncodeToString(hash[:])

	key, err := s.repo.GetByHash(ctx, keyHash)
	if err != nil {
		return nil, fmt.Errorf("invalid api key")
	}

	if !key.IsActive {
		return nil, fmt.Errorf("api key has been revoked")
	}

	if key.ExpiresAt.Valid && key.ExpiresAt.Time.Before(time.Now()) {
		return nil, fmt.Errorf("api key has expired")
	}

	// Fire-and-forget: update last_used_at
	go func() {
		_ = s.repo.TouchLastUsed(context.Background(), key.ID)
	}()

	return key, nil
}

// ValidateKeyFunc returns a closure compatible with middleware.APIKeyValidator.
func (s *Service) ValidateKeyFunc() func(ctx context.Context, plaintext string) (int, string, error) {
	return func(ctx context.Context, plaintext string) (int, string, error) {
		key, err := s.ValidateKey(ctx, plaintext)
		if err != nil {
			return 0, "", err
		}
		return key.ID, key.Name, nil
	}
}

// ListKeys returns all API keys (optionally including revoked ones).
func (s *Service) ListKeys(ctx context.Context, includeRevoked bool) ([]APIKeyDTO, error) {
	keys, err := s.repo.List(ctx, includeRevoked)
	if err != nil {
		return nil, err
	}

	dtos := make([]APIKeyDTO, len(keys))
	for i, k := range keys {
		dtos[i] = mapKeyToDTO(k)
	}
	return dtos, nil
}

// RevokeKey deactivates an API key.
func (s *Service) RevokeKey(ctx context.Context, id int) error {
	err := s.repo.Revoke(ctx, id)
	if err != nil {
		return err
	}

	// Record system log
	if s.logService != nil {
		s.logService.Record(ctx, logs.LogEntry{
			Category: logs.CategorySystem,
			Action:   logs.ActionAPIKeyRevoked,
			Message:  fmt.Sprintf("API key #%d has been revoked", id),
			Metadata: map[string]interface{}{
				"keyId": id,
			},
		})
	}

	return nil
}

func mapKeyToDTO(key APIKey) APIKeyDTO {
	dto := APIKeyDTO{
		ID:        key.ID,
		Name:      key.Name,
		KeyPrefix: key.KeyPrefix,
		IsActive:  key.IsActive,
		CreatedAt: key.CreatedAt,
	}

	if key.Scopes.Valid {
		var scopes []string
		_ = json.Unmarshal([]byte(key.Scopes.String), &scopes)
		dto.Scopes = scopes
	}

	if key.LastUsedAt.Valid {
		dto.LastUsedAt = &key.LastUsedAt.Time
	}
	if key.ExpiresAt.Valid {
		dto.ExpiresAt = &key.ExpiresAt.Time
	}

	return dto
}
