package datastore

import (
	"context"
	"time"
)

type RedisClientInterface interface {
	Set(
		ctx context.Context,
		key string,
		value interface{},
		expiration time.Duration,
	) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
	// SAdd, SRem, SMembers are also used in sessions
	SAdd(ctx context.Context, key string, members ...interface{}) error
	SRem(ctx context.Context, key string, members ...interface{}) error
	SMembers(ctx context.Context, key string) ([]string, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	Keys(ctx context.Context, pattern string) ([]string, error)
}
