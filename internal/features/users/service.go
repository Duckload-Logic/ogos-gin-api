package users

import (
	"context"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Service struct {
	repo RepositoryInterface
}

// NewService creates a new users service.
func NewService(repo RepositoryInterface) *Service {
	return &Service{repo: repo}
}

// GetUserByID retrieves a user by their ID.
func (s *Service) GetUserByID(
	ctx context.Context,
	userID string,
) (*GetUserResponse, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.mapUserModelToResponse(user), nil
}

// GetUserByEmail retrieves a user by their email and auth type.
func (s *Service) GetUserByEmail(
	ctx context.Context,
	email string,
	authType string,
) (*GetUserResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, email, authType)
	if err != nil {
		return nil, err
	}

	return s.mapUserModelToResponse(user), nil
}

func (s *Service) GetUserIDsByRole(
	ctx context.Context,
	roleID int,
) ([]string, error) {
	return s.repo.GetUserIDsByRole(ctx, roleID)
}

func (s *Service) ListUsers(
	ctx context.Context,
	params ListUsersParams,
) (*ListUsersResponse, error) {
	users, total, err := s.repo.ListUsers(ctx, params)
	if err != nil {
		return nil, err
	}

	dtos := make([]GetUserResponse, 0, len(users))
	for _, user := range users {
		dto := s.mapUserModelToResponse(&user)
		if dto != nil {
			dtos = append(dtos, *dto)
		}
	}

	totalPages := (total + params.PageSize - 1) / params.PageSize

	return &ListUsersResponse{
		Users: dtos,
		Meta: structs.PaginationMetadata{
			Total:      total,
			Page:       params.Page,
			PageSize:   params.PageSize,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *Service) GetRoleDistribution(
	ctx context.Context,
) ([]RoleDistributionDTO, error) {
	return s.repo.GetRoleDistribution(ctx)
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
		SuffixName: structs.FromSqlNull(user.SuffixName),
		Email:      user.Email,
		IsActive:   user.IsActive == 1,
		CreatedAt:  user.CreatedAt.Time.String(),
		UpdatedAt:  user.UpdatedAt.Time.String(),
	}
}

func (s *Service) BlockUser(ctx context.Context, userID string) error {
	return datastore.RunInTransaction(
		ctx,
		s.repo.GetDB(),
		func(tx datastore.DB) error {
			return s.repo.BlockUser(ctx, tx, userID)
		},
	)
}

func (s *Service) UnblockUser(ctx context.Context, userID string) error {
	return datastore.RunInTransaction(
		ctx,
		s.repo.GetDB(),
		func(tx datastore.DB) error {
			return s.repo.UnblockUser(ctx, tx, userID)
		},
	)
}
