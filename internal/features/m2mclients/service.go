package m2mclients

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/sessions"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
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

func (s *Service) CreateClient(
	ctx context.Context,
	userID string,
	req CreateM2MClientRequest,
) (*CreateM2MClientResponse, error) {
	existing, _ := s.repo.GetActiveByUserID(ctx, userID)
	if existing != nil {
		return nil, fmt.Errorf("developer already has an active m2m client")
	}

	rawSecret := make([]byte, 32)
	if _, err := rand.Read(rawSecret); err != nil {
		return nil, fmt.Errorf("failed to generate random secret: %w", err)
	}

	plaintextSecret := hex.EncodeToString(rawSecret)
	hash := sha256.Sum256([]byte(plaintextSecret))
	secretHash := hex.EncodeToString(hash[:])

	clientID := uuid.New().String()

	var scopesStr string
	if len(req.Scopes) > 0 {
		b, _ := json.Marshal(req.Scopes)
		scopesStr = string(b)
	}

	var expiresAt structs.NullableTime
	if req.ExpiresAt != nil {
		t, err := time.Parse(time.RFC3339, *req.ExpiresAt)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid expiresAt format (use RFC3339): %w",
				err,
			)
		}
		expiresAt = structs.TimeToNullableTime(t)
	}

	client := M2MClient{
		UserID:            userID,
		ClientName:        req.ClientName,
		ClientID:          clientID,
		ClientSecretHash:  secretHash,
		ClientDescription: req.ClientDescription,
		Scopes:            structs.StringToNullableString(scopesStr),
		IsVerified:        false,
		IsActive:          true,
		ExpiresAt:         expiresAt,
	}

	var id int
	err := s.repo.WithTransaction(ctx, func(tx datastore.DB) error {
		var err error
		id, err = s.repo.Create(ctx, tx, client)
		return err
	})

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

	client.ID = id
	dto := mapClientToDTO(client)

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

	if !client.IsVerified {
		dto.ClientID = "********"
	}

	return &CreateM2MClientResponse{
		M2MClientDTO: dto,
		ClientSecret: plaintextSecret,
	}, nil
}

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

	if !client.IsVerified {
		return nil, fmt.Errorf("client is pending superadmin verification")
	}

	if client.ExpiresAt.Valid && client.ExpiresAt.Time.Before(time.Now()) {
		return nil, fmt.Errorf("client has expired")
	}

	go func() {
		_ = s.repo.TouchLastUsed(
			context.Background(),
			nil,
			client.ID,
		)
	}()

	return client, nil
}

func (s *Service) IssueToken(
	ctx context.Context,
	client *M2MClient,
) (*M2MTokenResponse, error) {
	accessToken, claims, err := s.tokenService.GenerateToken(
		client.ClientName,
		"", // No userID for M2M
		[]int{},
		string(constants.AuthTypeM2M),
		constants.M2MAccessTokenMaxAge,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	claims.M2MClientID = client.ClientID

	accessSession := map[string]string{
		"clientId":   client.ClientID,
		"clientName": client.ClientName,
		"tokenType":  string(constants.AuthTypeM2M),
		"isVerified": fmt.Sprintf("%t", client.IsVerified),
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

	refreshToken, rClaims, err := s.tokenService.GenerateToken(
		client.ClientName,
		"",
		[]int{},
		string(constants.AuthTypeM2M),
		constants.M2MRefreshTokenMaxAge,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

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

func (s *Service) RefreshToken(
	ctx context.Context,
	refreshToken string,
) (*M2MTokenResponse, error) {
	claims, err := s.tokenService.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired refresh token")
	}

	session, err := s.sessionService.GetToken(ctx, sessions.NewJTI(claims.ID))
	if err != nil {
		return nil, fmt.Errorf("refresh session expired or revoked")
	}

	clientID := session["clientId"]
	if clientID == "" {
		return nil, fmt.Errorf("invalid session data")
	}

	client, err := s.repo.GetByClientID(ctx, clientID)
	if err != nil || !client.IsActive {
		return nil, fmt.Errorf("client is inactive or revoked")
	}

	newTokens, err := s.IssueToken(ctx, client)
	if err != nil {
		return nil, err
	}

	_ = s.sessionService.DeleteToken(ctx, sessions.NewJTI(claims.ID))

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

func (s *Service) ListClients(
	ctx context.Context,
	userID string,
	includeRevoked bool,
	roleIDs []int,
) ([]M2MClientDTO, error) {
	clients, err := s.repo.List(ctx, userID, includeRevoked)
	if err != nil {
		return nil, err
	}

	isSuperAdmin := false
	for _, rid := range roleIDs {
		if rid == int(constants.SuperAdminRoleID) {
			isSuperAdmin = true
			break
		}
	}

	dtos := make([]M2MClientDTO, len(clients))
	for i, c := range clients {
		dto := mapClientToDTO(c)

		if !isSuperAdmin && !c.IsVerified {
			dto.ClientID = "********"
			dto.Scopes = []string{"********"}
		}

		dtos[i] = dto
	}
	return dtos, nil
}

func (s *Service) RegenerateSecret(
	ctx context.Context,
	id int,
) (string, error) {
	rawSecret := make([]byte, 32)
	if _, err := rand.Read(rawSecret); err != nil {
		return "", fmt.Errorf("failed to generate random secret: %w", err)
	}
	plaintextSecret := hex.EncodeToString(rawSecret)

	sum := sha256.Sum256([]byte(plaintextSecret))
	secretHash := hex.EncodeToString(sum[:])

	err := s.repo.WithTransaction(
		ctx,
		func(tx datastore.DB) error {
			if err := s.repo.UpdateSecret(
				ctx, tx, id, secretHash,
			); err != nil {
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

func (s *Service) RevokeClient(ctx context.Context, id int) error {
	return s.repo.WithTransaction(
		ctx,
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

			actorID := audit.ExtractUserID(ctx)
			m2mClient, err := s.repo.GetByID(ctx, tx, id)
			if err != nil {
				return err
			}

			notifications := []audit.NotificationParams{
				{
					ReceiverID: structs.StringToNullableString(actorID),
					TargetID: structs.StringToNullableString(
						fmt.Sprintf("%d", id),
					),
					TargetType: structs.StringToNullableString(
						constants.M2MClientEntityType,
					),
					Title: "M2M Client Revoked",
					Message: fmt.Sprintf(
						"You have revoked M2M client %s",
						m2mClient.ClientName,
					),
					Type: constants.SystemEntityType,
				},
			}

			// Notify owner if different from actor
			if m2mClient.UserID != "" && m2mClient.UserID != actorID {
				notifications = append(notifications, audit.NotificationParams{
					ReceiverID: structs.StringToNullableString(
						m2mClient.UserID,
					),
					TargetID: structs.StringToNullableString(
						fmt.Sprintf("%d", id),
					),
					TargetType: structs.StringToNullableString(
						constants.M2MClientEntityType,
					),
					Title: "Your M2M Client has been Revoked",
					Message: fmt.Sprintf(
						"Your M2M client %s has been revoked",
						m2mClient.ClientName,
					),
					Type: constants.SystemEntityType,
				})
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
					Notifications: notifications,
				},
			)

			return nil
		},
	)
}

func (s *Service) VerifyClient(ctx context.Context, id int) error {
	return s.repo.WithTransaction(
		ctx,
		func(tx datastore.DB) error {
			err := s.repo.UpdateVerificationStatus(ctx, tx, id, true)
			if err != nil {
				return err
			}

			actorID := audit.ExtractUserID(ctx)
			m2mClient, err := s.repo.GetByID(ctx, tx, id)
			if err != nil {
				return err
			}

			notifications := []audit.NotificationParams{
				{
					ReceiverID: structs.StringToNullableString(actorID),
					TargetID: structs.StringToNullableString(
						fmt.Sprintf("%d", id),
					),
					TargetType: structs.StringToNullableString(
						constants.M2MClientEntityType,
					),
					Title: "M2M Client Verified",
					Message: fmt.Sprintf(
						"You have verified M2M client #%d", id,
					),
					Type: constants.SystemEntityType,
				},
			}

			// Notify owner if different from actor
			if m2mClient.UserID != "" && m2mClient.UserID != actorID {
				notifications = append(notifications, audit.NotificationParams{
					ReceiverID: structs.StringToNullableString(
						m2mClient.UserID,
					),
					TargetID: structs.StringToNullableString(
						fmt.Sprintf("%d", id),
					),
					TargetType: structs.StringToNullableString(
						constants.M2MClientEntityType,
					),
					Title: "Your M2M Client has been Verified",
					Message: fmt.Sprintf(
						"Your M2M client %s has been verified",
						structs.TruncateString(m2mClient.ClientName, 20),
					),
					Type: constants.SystemEntityType,
				})
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
						Action:   audit.ActionM2MClientVerified,
						Message: fmt.Sprintf(
							"M2M client #%d has been verified",
							id,
						),
						Metadata: &audit.LogMetadata{
							EntityType: constants.M2MClientEntityType,
							EntityID:   fmt.Sprintf("%d", id),
						},
					},
					Notifications: notifications,
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
