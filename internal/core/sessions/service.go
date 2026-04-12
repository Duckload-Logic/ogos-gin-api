package sessions

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Service struct {
	redis *datastore.RedisClient
}

func NewService(redis *datastore.RedisClient) *Service {
	return &Service{redis: redis}
}

// StoreToken saves session data in Redis with a JTI-based key.
func (s *Service) StoreToken(
	ctx context.Context,
	jti JTIDTO,
	data map[string]string,
	expireSeconds int,
) error {
	key := jti.ToSessionKey()
	valJSON, _ := json.Marshal(data)

	err := s.redis.Set(
		ctx,
		key,
		string(valJSON),
		time.Duration(expireSeconds)*time.Second,
	)
	if err != nil {
		return fmt.Errorf("failed to store token in redis: %w", err)
	}

	return nil
}

// GetToken retrieves session data from Redis.
func (s *Service) GetToken(
	ctx context.Context,
	jti JTIDTO,
) (map[string]string, error) {
	key := jti.ToSessionKey()
	val, err := s.redis.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("session not found or expired: %w", err)
	}

	var data map[string]string
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, fmt.Errorf("failed to parse session data: %w", err)
	}

	return data, nil
}

// DeleteToken removes a session from Redis.
func (s *Service) DeleteToken(ctx context.Context, jti JTIDTO) error {
	key := jti.ToSessionKey()
	return s.redis.Del(ctx, key)
}
