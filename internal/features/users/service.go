package users

import (
	"context"
	"fmt"

	"github.com/olazo-johnalbert/duckload-api/internal/core/sessions"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Service struct {
	repo           RepositoryInterface
	sessionService *sessions.Service
}

// NewService creates a new users service.
func NewService(
	repo RepositoryInterface,
	sessionService *sessions.Service,
) ServiceInterface {
	return &Service{
		repo:           repo,
		sessionService: sessionService,
	}
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
	return &GetUserResponse{
		ID:         user.ID,
		Roles:      user.Roles,
		FirstName:  user.FirstName,
		MiddleName: user.MiddleName,
		LastName:   user.LastName,
		SuffixName: user.SuffixName,
		Email:      user.Email,
		IsActive:   user.IsActive,
		CreatedAt:  user.CreatedAt.Time.String(),
		UpdatedAt:  user.UpdatedAt.Time.String(),
	}
}

func (s *Service) PostProfilePicture(
	ctx context.Context,
	userID string,
	fileID string,
) error {
	return s.repo.WithTransaction(
		ctx,
		func(tx datastore.DB) error {
			return s.repo.PostProfilePicture(ctx, tx, userID, fileID)
		},
	)
}

func (s *Service) BlockUser(ctx context.Context, userID string) error {
	err := s.repo.WithTransaction(
		ctx,
		func(tx datastore.DB) error {
			return s.repo.BlockUser(ctx, tx, userID)
		},
	)
	if err != nil {
		return err
	}

	// Revoke sessions
	_ = s.sessionService.RevokeAllUserSessions(ctx, userID)

	return nil
}

func (s *Service) UnblockUser(ctx context.Context, userID string) error {
	return s.repo.WithTransaction(
		ctx,
		func(tx datastore.DB) error {
			return s.repo.UnblockUser(ctx, tx, userID)
		},
	)
}

func (s *Service) UpdateUserRoles(
	ctx context.Context,
	req UpdateRolesRequest,
	adminID string,
) error {
	return s.repo.WithTransaction(ctx, func(tx datastore.DB) error {
		// Remove current roles
		if err := s.repo.RemoveRoles(ctx, tx, req.UserID); err != nil {
			return fmt.Errorf("failed to remove old roles: %w", err)
		}

		// Assign new roles
		for _, roleID := range req.RoleIDs {
			assignment := RoleAssignment{
				UserID:      req.UserID,
				RoleID:      roleID,
				AssignedBy:  structs.StringToNullableString(adminID),
				Reason:      structs.StringToNullableString(req.Reason),
				ReferenceID: structs.StringToNullableString(req.ReferenceID),
			}
			if err := s.repo.AssignRole(ctx, tx, assignment); err != nil {
				return fmt.Errorf("failed to assign role %d: %w", roleID, err)
			}
		}

		return nil
	})
}
