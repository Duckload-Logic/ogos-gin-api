package users

import "github.com/olazo-johnalbert/duckload-api/internal/core/structs"

type GetUserResponse struct {
	ID         int                    `json:"id"`
	Role       Role                   `json:"role"`
	FirstName  string                 `json:"firstName"`
	MiddleName structs.NullableString `json:"middleName,omitempty"`
	LastName   string                 `json:"lastName"`
	Email      string                 `json:"email"`
	CreatedAt  string                 `json:"createdAt,omitempty"`
	UpdatedAt  string                 `json:"updatedAt,omitempty"`
}

type CreateUserRequest struct {
	RoleID     int    `json:"roleId" binding:"required"`
	FirstName  string `json:"firstName" binding:"required"`
	MiddleName string `json:"middleName"`
	LastName   string `json:"lastName" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=8"`
}
