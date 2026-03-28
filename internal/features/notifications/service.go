package notifications

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Send handles creating a new notification using the email string as identifier
func (s *Service) Send(
	ctx context.Context,
	notif audit.NotificationEntry,
) error {
	return s.repo.Create(ctx, s.repo.GetDB(), &NotificationModel{
		ID: uuid.New().String(),

		ReceiverID: structs.ToSqlNull(notif.ReceiverID),
		ActorID:    structs.ToSqlNull(notif.ActorID),
		TargetID:   structs.ToSqlNull(notif.TargetID),
		TargetType: structs.ToSqlNull(notif.TargetType),

		Title:   notif.Title,
		Message: notif.Message,
		Type:    notif.Type,
	})
}

func (s *Service) GetUserNotifications(
	ctx context.Context,
	userID string,
) ([]audit.NotificationEntry, error) {
	models, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		// Consistent error wrapping with %w
		return nil, fmt.Errorf(
			"failed to fetch notifications for user %s: %w",
			userID,
			err,
		)
	}

	var dtos []audit.NotificationEntry
	for _, m := range models {
		dtos = append(dtos, audit.NotificationEntry{
			ID:         m.ID,
			ReceiverID: structs.NullableString(m.ReceiverID),
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
	return s.repo.MarkAsRead(ctx, s.repo.GetDB(), id)
}
