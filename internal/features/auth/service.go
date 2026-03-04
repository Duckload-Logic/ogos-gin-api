package auth

import (
	"context"
	"errors"

	"github.com/olazo-johnalbert/duckload-api/internal/core/tokens"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *users.Repository
}

func NewService(repo *users.Repository) *Service {
	return &Service{repo: repo}
}

var tokenService = tokens.NewService()

const accessTokenValidityMinutes = 60 * 1   // 60 minutes * 1 hour = 1 hour
const refreshTokenValidityMinutes = 60 * 12 // 60 minutes * 12 hours = 12 hours

// AuthenticateUser
func (s *Service) AuthenticateUser(
	ctx context.Context, email, password string,
) (string, string, error) {
	// Fetch user from database
	user, err := s.repo.GetUser(ctx, nil, &email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	// Compare hashed password
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	// Generate the token
	token, err := tokenService.GenerateToken(user.ID, user.RoleID, "access", accessTokenValidityMinutes)
	if err != nil {
		return "", "", errors.New("failed to generate session")
	}

	// Generate refresh token
	refreshToken, err := tokenService.GenerateToken(user.ID, user.RoleID, "refresh", refreshTokenValidityMinutes)
	if err != nil {
		return "", "", errors.New("failed to generate refresh token")
	}

	return token, refreshToken, nil
}

func (s *Service) RefreshToken(
	ctx context.Context, refreshToken string,
) (string, string, error) {
	claims, err := tokenService.ValidateToken(refreshToken)

	if err != nil {
		return "", "", errors.New("Invalid refresh token")
	}

	userID := claims.UserID
	roleID := claims.RoleID

	// Generate new token
	newToken, err := tokenService.GenerateToken(userID, roleID, "access", accessTokenValidityMinutes)
	if err != nil {
		return "", "", errors.New("Failed to generate new token")
	}

	// Generate new refresh token
	newRefreshToken, err := tokenService.GenerateToken(userID, roleID, "refresh", refreshTokenValidityMinutes)
	if err != nil {
		return "", "", errors.New("Failed to generate new refresh token")
	}

	return newToken, newRefreshToken, nil
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	// TODO: Implement token blacklisting if needed
	return nil
}
