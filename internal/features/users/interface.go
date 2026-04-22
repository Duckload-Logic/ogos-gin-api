package users

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type ServiceInterface interface {
	GetUserByID(ctx context.Context, userID string) (*GetUserResponse, error)
	GetUserByEmail(
		ctx context.Context,
		email string,
		authType string,
	) (*GetUserResponse, error)
	GetUserIDsByRole(ctx context.Context, roleID int) ([]string, error)
	ListUsers(
		ctx context.Context,
		params ListUsersParams,
	) (*ListUsersResponse, error)
	GetRoleDistribution(ctx context.Context) ([]RoleDistributionDTO, error)
	PostProfilePicture(ctx context.Context, userID string, fileID string) error
	BlockUser(ctx context.Context, userID string) error
	UnblockUser(ctx context.Context, userID string) error
}

type RepositoryInterface interface {
	WithTransaction(ctx context.Context, fn func(datastore.DB) error) error
	GetDB() *sqlx.DB
	GetUserByID(ctx context.Context, userID string) (*User, error)
	GetRoleByID(ctx context.Context, roleID int) (*Role, error)
	GetUserByEmail(
		ctx context.Context,
		email string,
		authType string,
	) (*User, error)
	GetUserIDsByRole(ctx context.Context, roleID int) ([]string, error)
	ListUsers(ctx context.Context, params ListUsersParams) ([]User, int, error)
	GetRoleDistribution(ctx context.Context) ([]RoleDistributionDTO, error)
	CreateUser(ctx context.Context, tx datastore.DB, user User) error
	PostProfilePicture(
		ctx context.Context,
		tx datastore.DB,
		userID string,
		fileID string,
	) error
	BlockUser(ctx context.Context, tx datastore.DB, userID string) error
	UnblockUser(ctx context.Context, tx datastore.DB, userID string) error
	CheckUserWhitelist(ctx context.Context, email string) (int, error)
}
