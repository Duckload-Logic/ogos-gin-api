package users

import (
	"context"
)

type ServiceInterface interface {
	GetUserByID(ctx context.Context, userID string) (*GetUserResponse, error)
	GetUserByEmail(
		ctx context.Context,
		email string,
		authType string,
	) (*GetUserResponse, error)
}

type RepositoryInterface interface {
	GetUserByID(ctx context.Context, userID string) (*User, error)
	GetRoleByID(ctx context.Context, roleID int) (*Role, error)
	GetUserByEmail(
		ctx context.Context,
		email string,
		authType string,
	) (*User, error)
	CreateUser(ctx context.Context, user User) error
}
