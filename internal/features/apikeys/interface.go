package apikeys

import "context"

type ServiceInterface interface {
	GenerateKey(
		ctx context.Context,
		req CreateAPIKeyRequest,
	) (*CreateAPIKeyResponse, error)
	ValidateKey(ctx context.Context, plaintext string) (*APIKey, error)
	ValidateKeyFunc() func(ctx context.Context, plaintext string) (int, string, error)
	ListKeys(ctx context.Context, includeRevoked bool) ([]APIKeyDTO, error)
	RevokeKey(ctx context.Context, id int) error
}

type RepositoryInterface interface {
	Create(ctx context.Context, key APIKey) (int, error)
	GetByHash(ctx context.Context, keyHash string) (*APIKey, error)
	List(ctx context.Context, includeRevoked bool) ([]APIKey, error)
	Revoke(ctx context.Context, id int) error
	TouchLastUsed(ctx context.Context, id int) error
}
