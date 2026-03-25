package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/clients/idp"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/tokens"
	"github.com/olazo-johnalbert/duckload-api/internal/database"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo      *users.Repository
	idpClient *idp.IDPClient
	redis     *database.RedisClient
}

func NewService(repo *users.Repository, redis *database.RedisClient) *Service {
	return &Service{
		repo:      repo,
		idpClient: idp.NewIDPClient(),
		redis:     redis,
	}
}

// tokenService is now called inline to ensure environment variables are loaded

// AuthenticateUser
func (s *Service) AuthenticateUser(
	ctx context.Context, email, password string,
) (string, string, string, error) {
	// Fetch user from database (Native only)
	user, err := s.repo.GetUserByEmail(ctx, email, "native")
	if err != nil {
		return "", "", "", errors.New("invalid credentials")
	}

	// Compare hashed password
	if !user.PasswordHash.Valid {
		return "", "", "", errors.New("invalid credentials")
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash.String),
		[]byte(password),
	)
	if err != nil {
		return "", "", "", errors.New("invalid credentials")
	}

	// Generate the token
	token, claims, err := tokens.NewService().GenerateToken(user.Email, user.ID, user.RoleID, "", "native", constants.AccessTokenMaxAge/60)
	if err != nil {
		return "", "", "", errors.New("failed to generate session")
	}

	// Generate refresh token
	refreshToken, _, err := tokens.NewService().GenerateToken(user.Email, user.ID, user.RoleID, "", "native", constants.RefreshTokenMaxAge/60)
	if err != nil {
		return "", "", "", errors.New("failed to generate refresh token")
	}

	// Store in Redis using the Token ID (jti)
	err = s.storeTokenInRedis(ctx, user.ID, claims.ID, "native", nil, constants.AccessTokenMaxAge/60)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to store token in redis: %v", err)
	}
	log.Println("Successfully stored token in redis")

	return user.ID, token, refreshToken, nil
}

func (s *Service) RefreshToken(
	ctx context.Context, refreshToken string, cfg *config.Config,
) (string, string, error) {
	claims, err := tokens.NewService().ValidateToken(refreshToken)
	if err != nil {
		return "", "", errors.New("Invalid refresh token")
	}

	// Check token type
	if claims.TokenType == "idp" {
		// Get IDP refresh token from Redis
		idpRefreshKey := fmt.Sprintf("idp_refresh:%s", claims.ID)
		idpRefreshToken, err := s.redis.Get(ctx, idpRefreshKey)
		if err != nil {
			return "", "", fmt.Errorf("[AuthService] {Get IDP Refresh Token}: token missing or expired in Redis")
		}

		// Call IDP refresh endpoint
		tokenResp, err := s.idpClient.RefreshToken(ctx, idpRefreshToken, cfg)
		if err != nil {
			return "", "", fmt.Errorf("[AuthService] {IDP Refresh}: %w", err)
		}

		// Generate NEW App Tokens
		newAppAccessToken, accessClaims, err := tokens.NewService().GenerateToken(
			claims.UserEmail,
			claims.UserID,
			claims.RoleID,
			"",
			"idp",
			constants.AccessTokenMaxAge/60,
		)
		if err != nil {
			return "", "", fmt.Errorf("[AuthService] {Generate App Access Token}: %w", err)
		}

		newAppRefreshToken, refreshClaims, err := tokens.NewService().GenerateToken(
			claims.UserEmail,
			claims.UserID,
			claims.RoleID,
			"",
			"idp",
			constants.RefreshTokenMaxAge/60,
		)
		if err != nil {
			return "", "", fmt.Errorf("[AuthService] {Generate App Refresh Token}: %w", err)
		}

		// Update Redis: App Access session
		idpAccess := tokenResp.AccessToken
		err = s.storeTokenInRedis(ctx, claims.UserID, accessClaims.ID, "idp", &idpAccess, constants.AccessTokenMaxAge/60)
		if err != nil {
			return "", "", fmt.Errorf("[AuthService] {Store Access in Redis}: %w", err)
		}

		// Update Redis: IDP Refresh linked to NEW App Refresh Token's ID
		newIdpRefreshKey := fmt.Sprintf("idp_refresh:%s", refreshClaims.ID)
		idpRefreshTokenToStore := tokenResp.RefreshToken
		if idpRefreshTokenToStore == "" {
			idpRefreshTokenToStore = idpRefreshToken // Reuse the existing one if IDP didn't rotate
		}
		err = s.redis.Set(ctx, newIdpRefreshKey, idpRefreshTokenToStore, time.Duration(constants.RefreshTokenMaxAge)*time.Second)
		if err != nil {
			return "", "", fmt.Errorf("[AuthService] {Store IDP Refresh in Redis}: %w", err)
		}

		// Optional: Clean up old refresh key
		_ = s.redis.Del(ctx, idpRefreshKey)

		return newAppAccessToken, newAppRefreshToken, nil
	}

	// Generate new token
	newToken, newClaims, err := tokens.NewService().GenerateToken(claims.UserEmail, claims.UserID, claims.RoleID, "", "native", constants.AccessTokenMaxAge/60)
	if err != nil {
		return "", "", errors.New("Failed to generate new token")
	}

	// Generate new refresh token
	newRefreshToken, _, err := tokens.NewService().GenerateToken(claims.UserEmail, claims.UserID, claims.RoleID, "", "native", constants.RefreshTokenMaxAge/60)
	if err != nil {
		return "", "", errors.New("Failed to generate new refresh token")
	}

	// Update Redis using new jti
	err = s.storeTokenInRedis(ctx, claims.UserID, newClaims.ID, "native", nil, constants.AccessTokenMaxAge/60)
	if err != nil {
		return "", "", fmt.Errorf("failed to update token in redis: %v", err)
	}

	return newToken, newRefreshToken, nil
}

func (s *Service) RefreshIDPToken(
	ctx context.Context, refreshToken string, cfg *config.Config,
) (string, string, error) {
	// Call IDP refresh endpoint
	tokenResp, err := s.idpClient.RefreshToken(ctx, refreshToken, cfg)
	if err != nil {
		return "", "", fmt.Errorf("[AuthService] {IDP Refresh}: %w", err)
	}

	return tokenResp.AccessToken, tokenResp.RefreshToken, nil
}

func (s *Service) GetMe(ctx context.Context, userID, tokenType string) (*MeResponse, error) {
	// only fetch user info for native tokens
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	role, err := s.repo.GetRoleByID(ctx, user.RoleID)
	if err != nil {
		// Fallback if role not found
		return &MeResponse{
			ID:         user.ID,
			Email:      user.Email,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			MiddleName: user.MiddleName.String,
			Roles:      []string{"user"},
			Type:       tokenType,
		}, nil
	}

	return &MeResponse{
		ID:         user.ID,
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: user.MiddleName.String,
		Roles:      []string{role.Name},
		Type:       tokenType,
	}, nil
}

func (s *Service) ParseIDPRoles(roles []string) []string {
	parsedRoles := make([]string, 0, len(roles))
	for _, role := range roles {
		// Split by ":" and take the last part
		parts := strings.Split(role, ":")
		if len(parts) > 1 {
			parsedRoles = append(parsedRoles, parts[len(parts)-1])
		} else {
			parsedRoles = append(parsedRoles, role)
		}
	}
	return parsedRoles
}

func (s *Service) storeTokenInRedis(ctx context.Context, userID, jti, tokenType string, idpAccessToken *string, expiryMinutes int) error {
	key := fmt.Sprintf("session:%s", jti)
	val := map[string]string{
		"userID":    userID,
		"tokenType": tokenType,
	}
	if idpAccessToken != nil {
		val["idpAccessToken"] = *idpAccessToken
	}
	valJSON, _ := json.Marshal(val)

	err := s.redis.Set(ctx, key, valJSON, time.Duration(expiryMinutes)*time.Minute)
	if err != nil {
		return fmt.Errorf("failed to store token in redis: %v", err)
	}

	return nil
}

func (s *Service) Logout(ctx context.Context, token string, tokenType string, cfg *config.Config) error {
	if tokenType != "" && tokenType == "idp" {
		_, _ = s.idpClient.Logout(ctx, cfg)
	}

	key := fmt.Sprintf("session:%s", token)
	return s.redis.Del(ctx, key)
}

// IDP integration methods

// GetAuthorizeURL generates the complete OAuth 2.0 authorization URL
// with PKCE parameters. This method creates a state parameter for CSRF
// protection, generates PKCE verifier and challenge, stores the state
// with metadata, and builds the authorization URL.
//
// Parameters:
//   - cfg: Application configuration containing IDP endpoints
//
// Returns the authorization URL and state parameter, or an error if
// generation fails.
func (s *Service) GetAuthorizeURL(
	cfg *config.Config,
) (string, error) {
	// Build authorization URL with all required parameters
	params := url.Values{}
	params.Set("client_id", cfg.IDPClientID)

	authURL := fmt.Sprintf(
		"%s?%s",
		cfg.IDPLoginURL,
		params.Encode(),
	)

	return authURL, nil
}

// PostIDPTokenExchange orchestrates the complete IDP login flow:
// validates state, exchanges code for token, fetches user info,
// provisions user, and generates application JWT tokens.
//
// Parameters:
//   - ctx: Context for database and HTTP operations
//   - code: Authorization code from IDP callback
//   - state: State parameter from IDP callback
//   - cfg: Application configuration
//
// Returns user ID and JWT tokens, or an error if any step fails.
func (s *Service) PostIDPTokenExchange(
	ctx context.Context,
	code string,
	cfg *config.Config,
) (string, string, string, string, string, error) {
	// Exchange authorization code for IDP tokens
	tokenResp, err := s.idpClient.ExchangeCodeForToken(ctx, code, cfg)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("[AuthService] {Token Exchange}: %w", err)
	}

	// Fetch User Info from IDP
	userInfo, err := s.GetIDPUserInfo(ctx, tokenResp.AccessToken, cfg)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("[AuthService] {Fetch User Info}: %w", err)
	}

	// Parse Tokens
	idpAccessToken := tokenResp.AccessToken
	idpRefreshToken := tokenResp.RefreshToken

	// Dynamic Role Mapping from IDP Tags
	// Logic: Identify tag:admin, tag:student, tag:superadmin at runtime
	appRoleID := s.mapIDPRolesToInternalID(userInfo.Roles)

	// Upsert IDP user into native database
	// This ensures non-existing users are added and existing ones are synchronized
	err = s.repo.CreateUser(ctx, users.User{
		ID:           userInfo.ID,
		RoleID:       appRoleID,
		FirstName:    userInfo.FirstName,
		LastName:     userInfo.LastName,
		Email:        userInfo.Email,
		AuthType:     "idp",
		PasswordHash: sql.NullString{Valid: false},
		IsActive:     1,
	})
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("[AuthService] {Provision IDP User}: %w", err)
	}

	appUserID := userInfo.ID

	// Map Role ID to Name for frontend redirection
	role, err := s.repo.GetRoleByID(ctx, appRoleID)
	roleName := "student"
	if err == nil && role != nil {
		roleName = role.Name
	}

	// Generate internal App Tokens using the actual app IDs
	appAccessToken, accessClaims, err := tokens.NewService().GenerateToken(
		userInfo.Email,
		appUserID,
		appRoleID,
		"",
		"idp",
		constants.AccessTokenMaxAge/60, // convert seconds to minutes for GenerateToken
	)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("[AuthService] {Generate App Access Token}: %w", err)
	}

	appRefreshToken, refreshClaims, err := tokens.NewService().GenerateToken(
		userInfo.Email,
		appUserID,
		appRoleID,
		"",
		"idp",
		constants.RefreshTokenMaxAge/60,
	)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("[AuthService] {Generate App Refresh Token}: %w", err)
	}

	// Store App Access Token in Redis using its ID (jti)
	if err := s.storeTokenInRedis(ctx, appUserID, accessClaims.ID, "idp", &idpAccessToken, constants.AccessTokenMaxAge/60); err != nil {
		return "", "", "", "", "", fmt.Errorf("[AuthService] {Store Access in Redis}: %w", err)
	}

	// Store IDP Refresh Token in Redis associated with the App Refresh Token's ID (jti)
	idpRefreshKey := fmt.Sprintf("idp_refresh:%s", refreshClaims.ID)
	idpRefreshTokenToStore := idpRefreshToken
	if idpRefreshTokenToStore == "" {
		// This shouldn't happen on initial login, but good for robustness
	}
	err = s.redis.Set(ctx, idpRefreshKey, idpRefreshTokenToStore, time.Duration(constants.RefreshTokenMaxAge)*time.Second)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("[AuthService] {Store IDP Refresh in Redis}: %w", err)
	}

	return appAccessToken, appRefreshToken, appUserID, userInfo.Email, roleName, nil
}

// GetIDPUserInfo fetches user information from the IDP userinfo endpoint
// using the provided access token. This is typically called after a
// successful token exchange to retrieve user details for provisioning.
//
// Parameters:
//   - ctx: Context for the HTTP request
//   - accessToken: Access token obtained from IDP token exchange
//   - cfg: Application configuration containing IDP endpoints
//
// Returns the IDP user information or an error if retrieval fails.
func (s *Service) GetIDPUserInfo(
	ctx context.Context,
	accessToken string,
	cfg *config.Config,
) (*idp.IDPUserInfo, error) {
	userInfo, err := s.idpClient.GetUserInfo(ctx, accessToken, cfg)
	if err != nil {
		return nil, fmt.Errorf(
			"[AuthService] {Get IDP User Info}: %w",
			err,
		)
	}
	return userInfo, nil
}

// mapIDPRolesToInternalID translates IDP tags to internal role IDs.
// Tags format: tag:student, tag:admin, tag:superadmin
func (s *Service) mapIDPRolesToInternalID(roles []string) int {
	// Priority order: superadmin > admin > student
	hasAdmin := false
	hasSuper := false
	hasStudent := false

	for _, r := range roles {
		if r == "" {
			continue
		}

		// Safely split and check for parts to avoid panics on missing colon
		parts := strings.Split(r, ":")
		rolePart := ""
		if len(parts) > 1 {
			rolePart = strings.ToLower(parts[1])
		} else {
			rolePart = strings.ToLower(parts[0])
		}

		switch rolePart {
		case "superadmin":
			hasSuper = true
		case "admin":
			hasAdmin = true
		case "student":
			hasStudent = true
		}
	}

	if hasSuper {
		return 3
	}
	if hasAdmin {
		return 2
	}
	if hasStudent {
		return 1
	}

	return 1 // Default to Student
}

// ValidateIDPSession checks if the provided session ID is valid on the IDP.
func (s *Service) ValidateIDPSession(
	ctx context.Context,
	sessionID string,
	cfg *config.Config,
) (*idp.IDPSessionResponse, error) {
	return s.idpClient.ValidateSession(ctx, sessionID, cfg)
}
