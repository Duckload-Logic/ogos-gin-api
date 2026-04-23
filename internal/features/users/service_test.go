package users

import (
	"context"
	"testing"

	"github.com/olazo-johnalbert/duckload-api/internal/core/sessions"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
	"go.uber.org/mock/gomock"
)

func TestService_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepositoryInterface(ctrl)
	mockRedis := datastore.NewMockRedisClientInterface(ctrl)
	sessionSvc := sessions.NewService(mockRedis)
	svc := NewService(mockRepo, sessionSvc)

	ctx := context.Background()
	userID := "user-123"

	t.Run("success", func(t *testing.T) {
		user := &User{
			ID:    userID,
			Email: "test@example.com",
			Roles: []Role{{ID: 1, Name: "Admin"}},
		}

		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(user, nil)

		resp, err := svc.GetUserByID(ctx, userID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp.Email != user.Email {
			t.Errorf("got %s, want %s", resp.Email, user.Email)
		}
	})
}

func TestService_BlockUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepositoryInterface(ctrl)
	mockRedis := datastore.NewMockRedisClientInterface(ctrl)
	sessionSvc := sessions.NewService(mockRedis)
	svc := NewService(mockRepo, sessionSvc)

	ctx := context.Background()
	userID := "user-123"

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().
			WithTransaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(datastore.DB) error) error {
				return fn(nil)
			})

		mockRepo.EXPECT().
			BlockUser(ctx, nil, userID).
			Return(nil)
		
		mockRedis.EXPECT().
			SMembers(ctx, "user:sessions:"+userID).
			Return(nil, nil)
		
		mockRedis.EXPECT().
			Del(ctx, "user:sessions:"+userID).
			Return(nil)

		err := svc.BlockUser(ctx, userID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
