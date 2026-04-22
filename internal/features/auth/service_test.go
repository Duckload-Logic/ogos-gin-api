package auth

import (
	"context"
	"database/sql"
	"testing"

	"github.com/olazo-johnalbert/duckload-api/internal/core/sessions"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/email"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/identity/idp"
	"go.uber.org/mock/gomock"
)

func TestService_BlockUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepositoryInterface(ctrl)
	mockRedis := datastore.NewMockRedisClientInterface(ctrl)
	mockEmailer := email.NewMockEmailer(ctrl)
	mockIDP := idp.NewMockIDPClientInterface(ctrl)

	sessionSvc := sessions.NewService(mockRedis)
	svc := &Service{
		repo:           mockRepo,
		redis:          mockRedis,
		sessionService: sessionSvc,
		emailer:        mockEmailer,
		idpClient:      mockIDP,
	}

	ctx := context.Background()
	userID := "test-user-id"

	t.Run("success", func(t *testing.T) {
		mockRepo.EXPECT().
			WithTransaction(ctx, gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(datastore.DB) error) error {
				return fn(nil) // Passing nil as tx for simplicity in mock
			})

		mockRepo.EXPECT().
			BlockUser(ctx, nil, userID).
			Return(nil)

		err := svc.BlockUser(ctx, userID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestService_VerifyUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepositoryInterface(ctrl)
	mockRedis := datastore.NewMockRedisClientInterface(ctrl)
	sessionSvc := sessions.NewService(mockRedis)
	svc := &Service{
		repo:           mockRepo,
		sessionService: sessionSvc,
	}

	ctx := context.Background()
	regID := "test-reg-id"
	otp := "123456"

	t.Run("invalid OTP", func(t *testing.T) {
		// Mock Redis getting the registration data
		// StoreToken was used with bcrypt, so we need a hashed value
		// For this test, we'll just return malformed or mismatched hash
		mockRedis.EXPECT().
			Get(ctx, sessions.NewJTI(regID).ToSessionKey()).
			Return(`{"verificationToken":"invalid-hash"}`, nil)

		_, _, err := svc.VerifyUser(ctx, regID, otp)
		if err == nil {
			t.Error("expected error for invalid OTP")
		}
	})
}

func TestService_GetMe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepositoryInterface(ctrl)
	svc := &Service{repo: mockRepo}

	ctx := context.Background()
	userID := "user-123"

	t.Run("success", func(t *testing.T) {
		user := &users.User{
			ID:     userID,
			Email:  "test@example.com",
			RoleID: 1,
		}
		role := &users.Role{
			ID:   1,
			Name: "Developer",
		}

		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(user, nil)
		mockRepo.EXPECT().GetRoleByID(ctx, user.RoleID).Return(role, nil)

		resp, err := svc.GetMe(ctx, userID, "native")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Email != user.Email || resp.Role.Name != role.Name {
			t.Errorf("resp mismatch: %+v", resp)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(nil, sql.ErrNoRows)

		_, err := svc.GetMe(ctx, userID, "native")
		if err == nil {
			t.Error("expected error for missing user")
		}
	})
}
