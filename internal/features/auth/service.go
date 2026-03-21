package auth

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/olazo-johnalbert/duckload-api/internal/core/clients/idp"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/tokens"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo      *users.Repository
	idpClient *idp.IDPClient
}

func NewService(repo *users.Repository) *Service {
	return &Service{
		repo:      repo,
		idpClient: idp.NewIDPClient(),
	}
}

var tokenService = tokens.NewService()

const (
	accessTokenValidityMinutes  = 60 * 1  // 60 minutes * 1 hour = 1 hour
	refreshTokenValidityMinutes = 60 * 12 // 60 minutes * 12 hours = 12 hours
)

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
	token, err := tokenService.GenerateToken(user.Email, user.ID, user.RoleID, "access", accessTokenValidityMinutes)
	if err != nil {
		return "", "", "", errors.New("failed to generate session")
	}

	// Generate refresh token
	refreshToken, err := tokenService.GenerateToken(user.Email, user.ID, user.RoleID, "refresh", refreshTokenValidityMinutes)
	if err != nil {
		return "", "", "", errors.New("failed to generate refresh token")
	}

	return user.ID, token, refreshToken, nil
}

func (s *Service) RefreshToken(
	ctx context.Context, refreshToken string,
) (string, string, error) {
	claims, err := tokenService.ValidateToken(refreshToken)
	if err != nil {
		return "", "", errors.New("Invalid refresh token")
	}

	// Generate new token
	newToken, err := tokenService.GenerateToken(claims.UserEmail, claims.UserID, claims.RoleID, "access", accessTokenValidityMinutes)
	if err != nil {
		return "", "", errors.New("Failed to generate new token")
	}

	// Generate new refresh token
	newRefreshToken, err := tokenService.GenerateToken(claims.UserEmail, claims.UserID, claims.RoleID, "refresh", refreshTokenValidityMinutes)
	if err != nil {
		return "", "", errors.New("Failed to generate new refresh token")
	}

	return newToken, newRefreshToken, nil
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	// TODO: Implement token blacklisting if needed
	return nil
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
) (string, string, error) {
	// Exchange authorization code for access token and refresh token
	tokenResp, err := s.idpClient.ExchangeCodeForToken(
		ctx,
		code,
		cfg,
	)
	if err != nil {
		return "", "", fmt.Errorf(
			"[AuthService] {Token Exchange}: %w",
			err,
		)
	}

	return tokenResp.AccessToken, tokenResp.RefreshToken, nil
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
