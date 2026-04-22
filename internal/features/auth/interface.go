package auth

import (
	"context"

	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/sessions"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/identity/idp"
)

type ServiceInterface interface {
	AuthenticateUser(
		ctx context.Context,
		email, password, ipAddress, userAgent string,
	) (string, string, string, error)
	RegisterUser(
		ctx context.Context,
		req RegisterDTO,
	) (string, error)
	ResendVerification(ctx context.Context, registrationID string) error
	VerifyUser(
		ctx context.Context,
		registrationID string,
		verificationOTP string,
	) (string, string, error)
	RefreshToken(
		ctx context.Context,
		accessTokenJTI sessions.JTIDTO,
		cfg *config.Config,
		ipAddress, userAgent string,
	) (string, string, error)
	RefreshIDPToken(
		ctx context.Context,
		refreshToken string,
		cfg *config.Config,
	) (string, string, error)
	GetMe(ctx context.Context, userID, tokenType string) (*MeResponse, error)
	Logout(
		ctx context.Context,
		token string,
		tokenType string,
		cfg *config.Config,
	) (string, error)
	GetAuthorizeURL(cfg *config.Config) (string, error)
	PostIDPTokenExchange(
		ctx context.Context,
		code string,
		cfg *config.Config,
		ipAddress, userAgent string,
	) (string, string, string, string, string, error)
	GetIDPUserInfo(
		ctx context.Context,
		accessToken string,
		cfg *config.Config,
	) (*idp.IDPUserInfo, error)
	BlockUser(ctx context.Context, userID string) error
	UnblockUser(ctx context.Context, userID string) error
}

type RepositoryInterface interface {
	GetUserByEmail(
		ctx context.Context,
		email string,
		authType string,
	) (*users.User, error)
	GetUserByID(ctx context.Context, userID string) (*users.User, error)
	GetRoleByID(ctx context.Context, roleID int) (*users.Role, error)
	CreateUser(ctx context.Context, tx datastore.DB, user users.User) error
	BlockUser(ctx context.Context, tx datastore.DB, userID string) error
	UnblockUser(ctx context.Context, tx datastore.DB, userID string) error
	CheckUserWhitelist(ctx context.Context, email string) (int, error)
	WithTransaction(ctx context.Context, fn func(datastore.DB) error) error
}
