package auth

import (
	"context"
	"errors"

	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo         *users.Repository
	tokenService *TokenService
}

func NewService(repo *users.Repository, tokenService *TokenService) *Service {
	return &Service{repo: repo, tokenService: tokenService}
}

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
	token, err := s.tokenService.GenerateToken(user.ID, user.RoleID, "access", 30)
	if err != nil {
		return "", "", errors.New("failed to generate session")
	}

	// Generate refresh token 
	refreshToken, err := s.tokenService.GenerateToken(user.ID, user.RoleID, "refresh", 1440)
	if err != nil {
		return "", "", errors.New("failed to generate refresh token")
	}

	return token, refreshToken, nil
}

func (s *Service) RefreshToken(
	ctx context.Context, refreshToken string,
) (string, string, error) {
	claims, err := s.tokenService.ValidateToken(refreshToken)

	if err != nil {
		return "", "", errors.New("Invalid refresh token")
	}

	userID := claims.UserID
	roleID := claims.RoleID

	// Generate new token 
	newToken, err := s.tokenService.GenerateToken(userID, roleID, "access", 30)
	if err != nil {
		return "", "", errors.New("Failed to generate new token")
	}

	// Generate new refresh token
	newRefreshToken, err := s.tokenService.GenerateToken(userID, roleID, "refresh", 1440)
	if err != nil {
		return "", "", errors.New("Failed to generate new refresh token")
	}

	return newToken, newRefreshToken, nil
}