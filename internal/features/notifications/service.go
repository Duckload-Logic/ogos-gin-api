package notifications

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
)

type Service struct {
	repo RepositoryInterface
}

func NewService(repo RepositoryInterface) *Service {
	return &Service{repo: repo}
}

func (s *Service) Send(
	ctx context.Context,
	notif audit.NotificationEntry,
) error {
	err := s.repo.Create(ctx, nil, Notification{
		ID:         uuid.New().String(),
		ReceiverID: notif.ReceiverID,
		ActorID:    notif.ActorID,
		TargetID:   notif.TargetID,
		TargetType: notif.TargetType,
		Title:      notif.Title,
		Message:    notif.Message,
		Type:       notif.Type,
	})
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	return nil
}

func (s *Service) GetUserNotifications(
	ctx context.Context,
	userID string,
) ([]audit.NotificationEntry, error) {
	models, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to fetch notifications for user %s: %w",
			userID,
			err,
		)
	}

	dtos := make([]audit.NotificationEntry, 0, len(models))
	for _, m := range models {
		dtos = append(dtos, audit.NotificationEntry{
			ID:         m.ID,
			ReceiverID: m.ReceiverID,
			Title:      m.Title,
			Message:    m.Message,
			Type:       m.Type,
			IsRead:     m.IsRead,
			CreatedAt:  m.CreatedAt,
		})
	}
	return dtos, nil
}

func (s *Service) MarkAsRead(ctx context.Context, id string) error {
	return s.repo.MarkAsRead(ctx, nil, id)
}
