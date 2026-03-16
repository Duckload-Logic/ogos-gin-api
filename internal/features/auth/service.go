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

type IDPUser struct {
	Email string   `json:"email"`
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
}

func NewService(repo *users.Repository) *Service {
	return &Service{repo: repo}
}

var tokenService = tokens.NewService()

const (
	accessTokenValidity  = 60      // 1 hour
	refreshTokenValidity = 60 * 12 // 12 hours
)

//Role Validation and DB Lookup
func (s *Service) SyncIDPUser(
	ctx context.Context, idpUser *IDPUser,
) (*users.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, idpUser.Email)
	if err != nil {
		return nil, errors.New("user not found in OGOS system")
	}

	//Gatekeeper role check
	isValid := false
	for _, role := range idpUser.Roles {
		if role == "Student" || role == "Counselor" {
			isValid = true
			break
		}
	}
	if !isValid {
		return nil, errors.New("insufficient permissions from IDP")
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

	at, rt, err := s.GenerateTokens(user)
	return user.ID, at, rt, err
}