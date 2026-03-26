package datastore

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(cfg *config.Config) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPass,
		DB:       cfg.RedisDB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &RedisClient{Client: rdb}, nil
}

func (r *RedisClient) Set(
	ctx context.Context,
	key string,
	value interface{},
	expiration time.Duration,
) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *RedisClient) Del(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}
