package users

import (
	"context"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetUserByID(
	ctx context.Context, userID int,
) (*GetUserResponse, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.mapUserModelToResponse(user), nil
}

func (s *Service) GetUserByEmail(
	ctx context.Context, email string,
) (*GetUserResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return s.mapUserModelToResponse(user), nil
}

func (s *Service) mapUserModelToResponse(user *User) *GetUserResponse {
	role, err := s.repo.GetRoleByID(context.Background(), user.RoleID)
	if err != nil {
		return nil
	}

	return &GetUserResponse{
		Role:       *role,
		ID:         user.ID,
		FirstName:  user.FirstName,
		MiddleName: structs.FromSqlNull(user.MiddleName),
		LastName:   user.LastName,
		Email:      user.Email,
		CreatedAt:  user.CreatedAt.Time.String(),
		UpdatedAt:  user.UpdatedAt.Time.String(),
	}
}
