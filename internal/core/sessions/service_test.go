package sessions

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
	"go.uber.org/mock/gomock"
)

func TestService_StoreToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := datastore.NewMockRedisClientInterface(ctrl)
	svc := NewService(mockRedis)

	ctx := context.Background()
	jti := NewJTI("test-jti")
	data := map[string]string{"foo": "bar"}
	expire := 60

	t.Run("success", func(t *testing.T) {
		valJSON, _ := json.Marshal(data)
		mockRedis.EXPECT().
			Set(ctx, jti.ToSessionKey(), string(valJSON), time.Duration(expire)*time.Second).
			Return(nil)

		err := svc.StoreToken(ctx, jti, data, expire)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestService_GetToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := datastore.NewMockRedisClientInterface(ctrl)
	svc := NewService(mockRedis)

	ctx := context.Background()
	jti := NewJTI("test-jti")
	data := map[string]string{"foo": "bar"}
	valJSON, _ := json.Marshal(data)

	t.Run("success", func(t *testing.T) {
		mockRedis.EXPECT().
			Get(ctx, jti.ToSessionKey()).
			Return(string(valJSON), nil)

		got, err := svc.GetToken(ctx, jti)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if got["foo"] != "bar" {
			t.Errorf("got %v, want foo=bar", got)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mockRedis.EXPECT().
			Get(ctx, jti.ToSessionKey()).
			Return("", context.DeadlineExceeded)

		_, err := svc.GetToken(ctx, jti)
		if err == nil {
			t.Error("expected error for missing session")
		}
	})
}

func TestService_DeleteToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := datastore.NewMockRedisClientInterface(ctrl)
	svc := NewService(mockRedis)

	ctx := context.Background()
	jti := NewJTI("test-jti")

	t.Run("success", func(t *testing.T) {
		mockRedis.EXPECT().
			Del(ctx, jti.ToSessionKey()).
			Return(nil)

		err := svc.DeleteToken(ctx, jti)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
