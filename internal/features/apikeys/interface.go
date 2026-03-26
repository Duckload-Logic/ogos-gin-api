package apikeys

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type ServiceInterface interface {
	GenerateKey(
		ctx context.Context,
		req CreateAPIKeyRequest,
	) (*CreateAPIKeyResponse, error)
	ValidateKey(ctx context.Context, plaintext string) (*APIKey, error)
	ValidateKeyFunc() func(
		ctx context.Context, plaintext string,
	) (int, string, error)
	ListKeys(ctx context.Context, includeRevoked bool) ([]APIKeyDTO, error)
	RevokeKey(ctx context.Context, id int) error
}

type RepositoryInterface interface {
	GetDB() *sqlx.DB
	Create(ctx context.Context, tx datastore.DB, key APIKey) (int, error)
	GetByHash(ctx context.Context, keyHash string) (*APIKey, error)
	List(ctx context.Context, includeRevoked bool) ([]APIKey, error)
	Revoke(ctx context.Context, tx datastore.DB, id int) error
	TouchLastUsed(ctx context.Context, tx datastore.DB, id int) error
}
