package users

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetUserByID(
	ctx context.Context, userID int,
) (*GetUserResponse, error) {
	user, err := s.repo.GetUser(ctx, &userID, nil)
	if err != nil {
		return nil, err
	}

	return s.mapUserModelToResponse(user), nil
}

func (s *Service) GetUserByEmail(
	ctx context.Context, email string,
) (*GetUserResponse, error) {
	user, err := s.repo.GetUser(ctx, nil, &email)
	if err != nil {
		return nil, err
	}

	return s.mapUserModelToResponse(user), nil
}

func (s *Service) mapUserModelToResponse(user *User) *GetUserResponse {
	return &GetUserResponse{
		ID:         user.ID,
		RoleID:     user.RoleID,
		FirstName:  user.FirstName,
		MiddleName: user.MiddleName.String,
		LastName:   user.LastName,
		Email:      user.Email,
		CreatedAt:  user.CreatedAt.Time.String(),
		UpdatedAt:  user.UpdatedAt.Time.String(),
	}
}
