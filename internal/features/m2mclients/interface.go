package m2mclients

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type ServiceInterface interface {
	CreateClient(
		ctx context.Context,
		userID string,
		req CreateM2MClientRequest,
	) (*CreateM2MClientResponse, error)
	Authenticate(
		ctx context.Context,
		clientID, clientSecret string,
	) (*M2MClient, error)
	IssueToken(
		ctx context.Context,
		client *M2MClient,
	) (*M2MTokenResponse, error)
	RefreshToken(
		ctx context.Context,
		refreshToken string,
	) (*M2MTokenResponse, error)
	ListClients(
		ctx context.Context,
		userID string,
		includeRevoked bool,
	) ([]M2MClientDTO, error)
	RevokeClient(ctx context.Context, id int) error
	RegenerateSecret(ctx context.Context, id int) (string, error)
	VerifyClient(ctx context.Context, id int) error
}

type RepositoryInterface interface {
	GetDB() *sqlx.DB
	Create(ctx context.Context, tx datastore.DB, client M2MClient) (int, error)
	GetByClientID(ctx context.Context, clientID string) (*M2MClient, error)
	GetActiveByUserID(ctx context.Context, userID string) (*M2MClient, error)
	List(
		ctx context.Context,
		userID string,
		includeRevoked bool,
	) ([]M2MClient, error)
	Revoke(ctx context.Context, tx datastore.DB, id int) error
	UpdateSecret(
		ctx context.Context,
		tx datastore.DB,
		id int,
		secretHash string,
	) error
	UpdateVerificationStatus(
		ctx context.Context,
		tx datastore.DB,
		id int,
		isVerified bool,
	) error
	TouchLastUsed(ctx context.Context, tx datastore.DB, id int) error
}
