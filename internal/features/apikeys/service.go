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

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Service struct {
	repo         RepositoryInterface
	logService   audit.Logger
	notifService audit.Notifier
}

func NewService(
	repo RepositoryInterface,
	logService audit.Logger,
	notifService audit.Notifier,
) *Service {
	return &Service{
		repo:         repo,
		logService:   logService,
		notifService: notifService,
	}
}

// GenerateKey creates a new API key, stores its hash, and returns the plaintext
// key (shown once).
func (s *Service) GenerateKey(
	ctx context.Context,
	req CreateAPIKeyRequest,
) (*CreateAPIKeyResponse, error) {
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
			return nil, fmt.Errorf(
				"invalid expiresAt format (use RFC3339): %w",
				err,
			)
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

	id, err := s.repo.Create(ctx, s.repo.GetDB(), apiKey)
	if err != nil {
		audit.Dispatch(ctx, s.logService, s.notifService, audit.DispatchParams{
			Log: &audit.LogParams{
				Category: audit.CategorySystem,
				Action:   audit.ActionAPIKeyCreateFailed,
				Message: fmt.Sprintf(
					"Failed to create API key '%s'",
					req.Name,
				),
				Metadata: &audit.LogMetadata{
					NewValues: req,
				},
			},
		})
		return nil, err
	}

	dto := mapKeyToDTO(apiKey)
	dto.ID = id

	audit.Dispatch(ctx, s.logService, s.notifService, audit.DispatchParams{
		Log: &audit.LogParams{
			Category: audit.CategorySystem,
			Action:   audit.ActionAPIKeyCreated,
			Message: fmt.Sprintf(
				"API key '%s' created (prefix: %s)",
				req.Name,
				prefix,
			),
			Metadata: &audit.LogMetadata{
				EntityType: constants.APIKeyEntityType,
				EntityID:   prefix,
				NewValues:  req,
			},
		},
		Notification: &audit.NotificationParams{
			Type:    constants.APIKeyEntityType,
			Title:   "API Key created",
			Message: fmt.Sprintf("API key '%s' has been created", req.Name),
		},
	})

	return &CreateAPIKeyResponse{
		APIKeyDTO: dto,
		Key:       plaintext,
	}, nil
}

// ValidateKey checks if a plaintext key is valid, active, and not expired.
// Returns the APIKey record if valid, or an error.
func (s *Service) ValidateKey(
	ctx context.Context,
	plaintext string,
) (*APIKey, error) {
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
		_ = s.repo.TouchLastUsed(context.Background(), s.repo.GetDB(), key.ID)
	}()

	return key, nil
}

// ValidateKeyFunc returns a closure compatible with middleware.APIKeyValidator.
func (s *Service) ValidateKeyFunc() func(
	ctx context.Context, plaintext string,
) (int, string, error) {
	return func(ctx context.Context, plaintext string) (int, string, error) {
		key, err := s.ValidateKey(ctx, plaintext)
		if err != nil {
			return 0, "", err
		}
		return key.ID, key.Name, nil
	}
}

// ListKeys returns all API keys (optionally including revoked ones).
func (s *Service) ListKeys(
	ctx context.Context,
	includeRevoked bool,
) ([]APIKeyDTO, error) {
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
	return datastore.RunInTransaction(
		ctx,
		s.repo.GetDB(),
		func(tx datastore.DB) error {
			err := s.repo.Revoke(ctx, tx, id)
			if err != nil {
				audit.Dispatch(
					ctx,
					s.logService,
					s.notifService,
					audit.DispatchParams{
						Tx: tx,
						Log: &audit.LogParams{
							Category: audit.CategorySystem,
							Action:   audit.ActionAPIKeyRevokeFailed,
							Message: fmt.Sprintf(
								"Failed to revoke API key #%d",
								id,
							),
							Metadata: &audit.LogMetadata{
								EntityType: constants.APIKeyEntityType,
								EntityID:   fmt.Sprintf("%d", id),
							},
						},
					},
				)
				return err
			}

			// Record system log
			audit.Dispatch(
				ctx,
				s.logService,
				s.notifService,
				audit.DispatchParams{
					Tx: tx,
					Log: &audit.LogParams{
						Category: audit.CategorySystem,
						Action:   audit.ActionAPIKeyRevoked,
						Message: fmt.Sprintf(
							"API key #%d has been revoked",
							id,
						),
						Metadata: &audit.LogMetadata{
							EntityType: constants.APIKeyEntityType,
							EntityID:   fmt.Sprintf("%d", id),
						},
					},
				},
			)

			return nil
		},
	)
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
