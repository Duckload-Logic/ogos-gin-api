package m2mclients

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/sessions"
	"github.com/olazo-johnalbert/duckload-api/internal/core/tokens"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Service struct {
	repo           RepositoryInterface
	logService     audit.Logger
	notifService   audit.Notifier
	tokenService   *tokens.Service
	sessionService *sessions.Service
}

func NewService(
	repo RepositoryInterface,
	logService audit.Logger,
	notifService audit.Notifier,
	tokenService *tokens.Service,
	sessionService *sessions.Service,
) *Service {
	return &Service{
		repo:           repo,
		logService:     logService,
		notifService:   notifService,
		tokenService:   tokenService,
		sessionService: sessionService,
	}
}

// CreateClient creates a new M2M client, returns the client and plaintext
// secret.
func (s *Service) CreateClient(
	ctx context.Context,
	userID string,
	req CreateM2MClientRequest,
) (*CreateM2MClientResponse, error) {
	// Restriction: Only one active M2M client per developer account
	existing, _ := s.repo.GetActiveByUserID(ctx, userID)
	if existing != nil {
		return nil, fmt.Errorf("developer already has an active m2m client")
	}

	// Generate a random 32-byte secret
	rawSecret := make([]byte, 32)
	if _, err := rand.Read(rawSecret); err != nil {
		return nil, fmt.Errorf("failed to generate random secret: %w", err)
	}

	plaintextSecret := hex.EncodeToString(rawSecret)
	hash := sha256.Sum256([]byte(plaintextSecret))
	secretHash := hex.EncodeToString(hash[:])

	clientID := uuid.New().String()

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

	client := M2MClient{
		UserID:            userID,
		ClientName:        req.ClientName,
		ClientID:          clientID,
		ClientSecretHash:  secretHash,
		ClientDescription: req.ClientDescription,
		Scopes:            scopesJSON,
		IsVerified:        true, // Auto-verified for now
		IsActive:          true,
		ExpiresAt:         expiresAt,
	}

	id, err := s.repo.Create(ctx, s.repo.GetDB(), client)
	if err != nil {
		audit.Dispatch(ctx, s.logService, s.notifService, audit.DispatchParams{
			Log: &audit.LogParams{
				Level:    audit.LevelError,
				Category: audit.CategorySystem,
				Action:   audit.ActionM2MClientCreateFailed,
				Message: fmt.Sprintf(
					"Failed to create M2M client '%s'",
					req.ClientName,
				),
				Metadata: &audit.LogMetadata{
					NewValues: req,
				},
			},
		})
		return nil, err
	}

	dto := mapClientToDTO(client)
	dto.ID = id

	audit.Dispatch(ctx, s.logService, s.notifService, audit.DispatchParams{
		Log: &audit.LogParams{
			Level:    audit.LevelInfo,
			Category: audit.CategorySystem,
			Action:   audit.ActionM2MClientCreated,
			Message: fmt.Sprintf(
				"M2M client '%s' created (ID: %s)",
				req.ClientName,
				clientID,
			),
			Metadata: &audit.LogMetadata{
				EntityType: constants.M2MClientEntityType,
				EntityID:   clientID,
				NewValues:  req,
			},
		},
	})

	return &CreateM2MClientResponse{
		M2MClientDTO: dto,
		ClientSecret: plaintextSecret,
	}, nil
}

// Authenticate validates client credentials and returns the client record.
func (s *Service) Authenticate(
	ctx context.Context,
	clientID, clientSecret string,
) (*M2MClient, error) {
	client, err := s.repo.GetByClientID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("invalid client credentials")
	}

	hash := sha256.Sum256([]byte(clientSecret))
	secretHash := hex.EncodeToString(hash[:])

	if client.ClientSecretHash != secretHash {
		return nil, fmt.Errorf("invalid client credentials")
	}

	if !client.IsActive {
		return nil, fmt.Errorf("client has been revoked")
	}

	if client.ExpiresAt.Valid && client.ExpiresAt.Time.Before(time.Now()) {
		return nil, fmt.Errorf("client has expired")
	}

	// Update last_used_at
	go func() {
		_ = s.repo.TouchLastUsed(
			context.Background(),
			s.repo.GetDB(),
			client.ID,
		)
	}()

	return client, nil
}

// IssueToken generates a new JWT and stores it in Redis as a session.
func (s *Service) IssueToken(
	ctx context.Context,
	client *M2MClient,
) (*M2MTokenResponse, error) {
	accessToken, claims, err := s.tokenService.GenerateToken(
		client.ClientName,
		"", // No userID for M2M
		0,  // No roleID for M2M
		"m2m_client",
		string(constants.AuthTypeM2M),
		constants.M2MAccessTokenMaxAge,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Set M2M marker in claims
	claims.M2MClientID = client.ClientID

	// Access Token session in Redis
	accessSession := map[string]string{
		"clientId":   client.ClientID,
		"clientName": client.ClientName,
		"tokenType":  string(constants.AuthTypeM2M),
	}
	err = s.sessionService.StoreToken(
		ctx,
		sessions.NewJTI(claims.ID),
		accessSession,
		constants.M2MAccessTokenMaxAge,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to store m2m access session: %w", err)
	}

	// Refresh Token (24 hours)
	refreshToken, rClaims, err := s.tokenService.GenerateToken(
		client.ClientName,
		"",
		0,
		"m2m_client",
		string(constants.AuthTypeM2M),
		constants.M2MRefreshTokenMaxAge,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Refresh Token session in Redis
	refreshSession := map[string]string{
		"clientId":   client.ClientID,
		"clientName": client.ClientName,
		"tokenType":  "m2m_refresh",
	}
	err = s.sessionService.StoreToken(
		ctx,
		sessions.NewJTI(rClaims.ID),
		refreshSession,
		constants.M2MRefreshTokenMaxAge,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to store m2m refresh session: %w", err)
	}

	return &M2MTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}, nil
}

// RefreshToken validates a refresh token and issues a new pair.
func (s *Service) RefreshToken(
	ctx context.Context,
	refreshToken string,
) (*M2MTokenResponse, error) {
	// Validate the refresh token
	claims, err := s.tokenService.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired refresh token")
	}

	// Check Redis session for this refresh token
	session, err := s.sessionService.GetToken(ctx, sessions.NewJTI(claims.ID))
	if err != nil {
		return nil, fmt.Errorf("refresh session expired or revoked")
	}

	clientID := session["clientId"]
	if clientID == "" {
		return nil, fmt.Errorf("invalid session data")
	}

	// Fetch client from DB to ensure they are still active
	client, err := s.repo.GetByClientID(ctx, clientID)
	if err != nil || !client.IsActive {
		return nil, fmt.Errorf("client is inactive or revoked")
	}

	// Issue new token pair
	newTokens, err := s.IssueToken(ctx, client)
	if err != nil {
		return nil, err
	}

	// Revoke the OLD refresh token session
	_ = s.sessionService.DeleteToken(ctx, sessions.NewJTI(claims.ID))

	// Log the refresh event
	audit.Dispatch(ctx, s.logService, s.notifService, audit.DispatchParams{
		Log: &audit.LogParams{
			Level:    audit.LevelInfo,
			Category: audit.CategorySecurity,
			Action:   audit.ActionM2MTokenRefreshed,
			Message: fmt.Sprintf(
				"M2M token refreshed for client: %s",
				client.ClientName,
			),
			Metadata: &audit.LogMetadata{
				EntityType: constants.M2MClientEntityType,
				EntityID:   client.ClientID,
			},
		},
	})

	return newTokens, nil
}

// ListClients returns all M2M clients for a specific user.
func (s *Service) ListClients(
	ctx context.Context,
	userID string,
	includeRevoked bool,
) ([]M2MClientDTO, error) {
	clients, err := s.repo.List(ctx, userID, includeRevoked)
	if err != nil {
		return nil, err
	}

	dtos := make([]M2MClientDTO, len(clients))
	for i, c := range clients {
		dtos[i] = mapClientToDTO(c)
	}
	return dtos, nil
}

// RegenerateSecret generates a new secret for an existing M2M client.
func (s *Service) RegenerateSecret(
	ctx context.Context,
	id int,
) (string, error) {
	// Generate a new random 32-byte secret
	rawSecret := make([]byte, 32)
	if _, err := rand.Read(rawSecret); err != nil {
		return "", fmt.Errorf("failed to generate random secret: %w", err)
	}
	plaintextSecret := base64.URLEncoding.EncodeToString(rawSecret)

	// Hash the new secret
	sum := sha256.Sum256([]byte(plaintextSecret))
	secretHash := hex.EncodeToString(sum[:])

	err := datastore.RunInTransaction(
		ctx,
		s.repo.GetDB(),
		func(tx datastore.DB) error {
			if err := s.repo.UpdateSecret(ctx, tx, id, secretHash); err != nil {
				return err
			}

			audit.Dispatch(
				ctx,
				s.logService,
				s.notifService,
				audit.DispatchParams{
					Tx: tx,
					Log: &audit.LogParams{
						Level:    audit.LevelInfo,
						Category: audit.CategorySystem,
						Action:   audit.ActionM2MClientSecretRotated,
						Message: fmt.Sprintf(
							"M2M client #%d secret has been rotated",
							id,
						),
						Metadata: &audit.LogMetadata{
							EntityType: constants.M2MClientEntityType,
							EntityID:   fmt.Sprintf("%d", id),
						},
					},
				},
			)

			return nil
		},
	)
	if err != nil {
		return "", err
	}

	return plaintextSecret, nil
}

// RevokeClient deactivates an M2M client.
func (s *Service) RevokeClient(ctx context.Context, id int) error {
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
							Level:    audit.LevelError,
							Category: audit.CategorySystem,
							Action:   audit.ActionM2MClientRevokeFailed,
							Message: fmt.Sprintf(
								"Failed to revoke M2M client #%d",
								id,
							),
							Metadata: &audit.LogMetadata{
								EntityType: constants.M2MClientEntityType,
								EntityID:   fmt.Sprintf("%d", id),
							},
						},
					},
				)
				return err
			}

			audit.Dispatch(
				ctx,
				s.logService,
				s.notifService,
				audit.DispatchParams{
					Tx: tx,
					Log: &audit.LogParams{
						Level:    audit.LevelInfo,
						Category: audit.CategorySystem,
						Action:   audit.ActionM2MClientRevoked,
						Message: fmt.Sprintf(
							"M2M client #%d has been revoked",
							id,
						),
						Metadata: &audit.LogMetadata{
							EntityType: constants.M2MClientEntityType,
							EntityID:   fmt.Sprintf("%d", id),
						},
					},
				},
			)

			return nil
		},
	)
}

func mapClientToDTO(client M2MClient) M2MClientDTO {
	dto := M2MClientDTO{
		ID:                client.ID,
		ClientName:        client.ClientName,
		ClientID:          client.ClientID,
		ClientDescription: client.ClientDescription,
		IsActive:          client.IsActive,
		IsVerified:        client.IsVerified,
		CreatedAt:         client.CreatedAt,
	}

	if client.Scopes.Valid {
		var scopes []string
		_ = json.Unmarshal([]byte(client.Scopes.String), &scopes)
		dto.Scopes = scopes
	}

	if client.LastUsedAt.Valid {
		dto.LastUsedAt = &client.LastUsedAt.Time
	}
	if client.ExpiresAt.Valid {
		dto.ExpiresAt = &client.ExpiresAt.Time
	}

	return dto
}
