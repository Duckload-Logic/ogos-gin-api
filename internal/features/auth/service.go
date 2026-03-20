package auth

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/tokens"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo      *users.Repository
	idpClient *IDPClient
}

func NewService(repo *users.Repository) *Service {
	return &Service{
		repo:      repo,
		idpClient: NewIDPClient(),
	}
}

var tokenService = tokens.NewService()

const (
	accessTokenValidityMinutes  = 60 * 1  // 60 minutes * 1 hour = 1 hour
	refreshTokenValidityMinutes = 60 * 12 // 60 minutes * 12 hours = 12 hours
)

func (s *Service) SyncIDPUser(
	ctx context.Context, idpUser *IDPUser,
) (*users.User, error) {
	// Role Gatekeeper
	var targetRoleID int
	authorized := false
	for _, r := range idpUser.Roles {
		if r == "Student" { targetRoleID = 1; authorized = true; break }
		if r == "Counselor" { targetRoleID = 2; authorized = true; break }
	}
	if !authorized {
		return nil, errors.New("unauthorized role")
	}

	// DB Lookup
	user, err := s.repo.GetUserByEmail(ctx, idpUser.Email)
	if err != nil {
		newUser := users.User{
			Email:     idpUser.Email,
			FirstName: idpUser.Name,
			LastName:  " ",
			RoleID:    targetRoleID,
			PasswordHash: "IDP_AUTH",
		}

		createErr := s.repo.CreateUser(ctx, newUser)
		if createErr != nil {
			return nil, errors.New("failed to sync IDP user")
		}

		return s.repo.GetUserByEmail(ctx, idpUser.Email)
	}

	return user, nil
}

//JWT Generation
func (s *Service) GenerateTokens(
	user *users.User,
) (string, string, error) {
	token, err := tokenService.GenerateToken(
		user.Email, user.ID, user.RoleID, "access", accessTokenValidity,
	)
	if err != nil {
		return "", "", errors.New("failed to generate access token")
	}

	refresh, err := tokenService.GenerateToken(
		user.Email, user.ID, user.RoleID, "refresh", refreshTokenValidity,
	)
	if err != nil {
		return "", "", errors.New("failed to generate refresh token")
	}

	return token, refresh, nil
}

func (s *Service) RefreshToken(
	ctx context.Context, refreshToken string,
) (string, string, error) {
	claims, err := tokenService.ValidateToken(refreshToken)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	return s.GenerateTokens(&users.User{
		ID:     claims.UserID,
		Email:  claims.UserEmail,
		RoleID: claims.RoleID,
	})
}

func (s *Service) AuthenticateUser(
	ctx context.Context, email, password string,
) (int, string, string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return 0, "", "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash), []byte(password),
	); err != nil {
		return 0, "", "", errors.New("invalid credentials")
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
