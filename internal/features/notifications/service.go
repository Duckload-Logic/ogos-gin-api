package notifications

import (
	"context"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Send(ctx context.Context, userID int, title, message, notifType string) error {
	return s.repo.Create(ctx, userID, title, message, notifType)
}

func (s *Service) GetUserNotifications(ctx context.Context, userID string) ([]NotificationDTO, error) {
	models, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var dtos []NotificationDTO
	for _, m := range models {
		dtos = append(dtos, NotificationDTO{
			ID:        uint(m.ID),
			Title:     m.Title,
			Message:   m.Message,
			Type:      m.Type,
			IsRead:    m.IsRead,
			CreatedAt: m.CreatedAt,
		})
	}
	return dtos, nil
}