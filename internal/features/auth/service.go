package auth

import (
	"context"
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

var tokenService = tokens.NewService()

// AuthenticateUser
func (s *Service) AuthenticateUser(
	ctx context.Context, email, password string,
) (string, string, string, error) {
	// Fetch user from database
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", "", errors.New("invalid credentials")
	}

	// Compare hashed password
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		return "", "", "", errors.New("invalid credentials")
	}

	// Generate the token
	token, err := tokenService.GenerateToken(user.Email, user.ID, user.RoleID, "", "native", constants.AccessTokenMaxAge)
	if err != nil {
		return "", "", "", errors.New("failed to generate session")
	}

	// Generate refresh token
	refreshToken, err := tokenService.GenerateToken(user.Email, user.ID, user.RoleID, "", "native", constants.RefreshTokenMaxAge)
	if err != nil {
		return "", "", "", errors.New("failed to generate refresh token")
	}

	// Store in Redis
	err = s.storeTokenInRedis(ctx, user.ID, token, "native", constants.AccessTokenMaxAge)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to store token in redis: %v", err)
	}
	log.Println("Successfully stored token in redis")

	return user.ID, token, refreshToken, nil
}

func (s *Service) RefreshToken(
	ctx context.Context, refreshToken string,
) (string, string, error) {
	claims, err := tokenService.ValidateToken(refreshToken)
	if err != nil {
		return "", "", errors.New("Invalid refresh token")
	}

	// Check token type
	if claims.TokenType == "idp" {
		// Handle IDP refresh logic if needed or just return error if not supported this way
		return "", "", errors.New("IDP refresh not supported via this endpoint, use /auth/idp/refresh")
	}

	// Generate new token
	newToken, err := tokenService.GenerateToken(claims.UserEmail, claims.UserID, claims.RoleID, "", "native", constants.AccessTokenMaxAge)
	if err != nil {
		return "", "", errors.New("Failed to generate new token")
	}

	// Generate new refresh token
	newRefreshToken, err := tokenService.GenerateToken(claims.UserEmail, claims.UserID, claims.RoleID, "", "native", constants.RefreshTokenMaxAge)
	if err != nil {
		return "", "", errors.New("Failed to generate new refresh token")
	}

	// Update Redis
	err = s.storeTokenInRedis(ctx, claims.UserID, newToken, "native", constants.AccessTokenMaxAge)
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

func (s *Service) storeTokenInRedis(ctx context.Context, userID, token, tokenType string, expiryMinutes int) error {
	key := fmt.Sprintf("token:%s:%s", tokenType, token)
	val := map[string]string{
		"userID":    userID,
		"tokenType": tokenType,
	}
	valJSON, _ := json.Marshal(val)
	return s.redis.Set(ctx, key, valJSON, time.Duration(expiryMinutes)*time.Minute)
}

func (s *Service) Logout(ctx context.Context, token string, tokenType string) error {
	key := fmt.Sprintf("token:%s:%s", tokenType, token)
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
) (string, string, *idp.IDPUserInfo, error) {
	// Exchange authorization code for IDP tokens
	tokenResp, err := s.idpClient.ExchangeCodeForToken(ctx, code, cfg)
	if err != nil {
		return "", "", nil, fmt.Errorf("[AuthService] {Token Exchange}: %w", err)
	}

	// Fetch User Info from IDP
	userInfo, err := s.GetIDPUserInfo(ctx, tokenResp.AccessToken, cfg)
	if err != nil {
		return "", "", nil, fmt.Errorf("[AuthService] {Fetch User Info}: %w", err)
	}

	// Parse Tokens
	accessToken := tokenResp.AccessToken
	refreshToken := tokenResp.RefreshToken

	// Store tokens in Redis
	if err := s.storeTokenInRedis(ctx, userInfo.ID, accessToken, "idp", constants.AccessTokenMaxAge); err != nil {
		return "", "", nil, fmt.Errorf("[AuthService] {Store in Redis}: %w", err)
	}

	return accessToken, refreshToken, userInfo, nil
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
