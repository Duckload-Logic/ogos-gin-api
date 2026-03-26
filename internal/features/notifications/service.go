package notifications

import (
	"context"
	"fmt"
)

type Service struct {
	repo RepositoryInterface
}

func NewService(repo RepositoryInterface) *Service {
	return &Service{repo: repo}
}

// Send handles creating a new notification using the email string as identifier
func (s *Service) Send(
	ctx context.Context,
	userID string,
	title, message, notifType string,
) error {
	return s.repo.Create(ctx, s.repo.GetDB(), userID, title, message, notifType)
}

func (s *Service) GetUserNotifications(
	ctx context.Context,
	userID string,
) ([]NotificationDTO, error) {
	models, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		// Consistent error wrapping with %w
		return nil, fmt.Errorf(
			"failed to fetch notifications for user %s: %w",
			userID,
			err,
		)
	}

	var dtos []NotificationDTO
	for _, m := range models {
		dtos = append(dtos, NotificationDTO{
			ID:        uint(m.ID),
			UserID:    m.UserID,
			Title:     m.Title,
			Message:   m.Message,
			Type:      m.Type,
			IsRead:    m.IsRead,
			CreatedAt: m.CreatedAt,
		})
	}
	return dtos, nil
}

func (s *Service) MarkAsRead(ctx context.Context, id int) error {
	return s.repo.MarkAsRead(ctx, s.repo.GetDB(), id)
}
