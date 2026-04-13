package users

import "github.com/olazo-johnalbert/duckload-api/internal/core/structs"

type GetUserResponse struct {
	ID         string                 `json:"id"`
	Role       Role                   `json:"role"`
	FirstName  string                 `json:"firstName"`
	MiddleName structs.NullableString `json:"middleName,omitempty"`
	LastName   string                 `json:"lastName"`
	SuffixName structs.NullableString `json:"suffixName,omitempty"`
	Email      string                 `json:"email,omitempty"`
	IsActive   bool                   `json:"isActive"`
	CreatedAt  string                 `json:"createdAt,omitempty"`
	UpdatedAt  string                 `json:"updatedAt,omitempty"`
}

type CreateUserRequest struct {
	RoleID     int    `json:"roleId"     binding:"required"`
	FirstName  string `json:"firstName"  binding:"required"`
	MiddleName string `json:"middleName"`
	LastName   string `json:"lastName"   binding:"required"`
	SuffixName string `json:"suffixName"`
	Email      string `json:"email"      binding:"required,email"`
	Password   string `json:"password"   binding:"required,min=8"`
}

type ListUsersParams struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	RoleID   int    `form:"role_id"`
	Search   string `form:"search"`
	Active   *bool  `form:"active"`
}

type ListUsersResponse struct {
	Users      []GetUserResponse `json:"users"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"pageSize"`
	TotalPages int               `json:"totalPages"`
}

type RoleDistributionDTO struct {
	RoleName string `json:"roleName" db:"role_name"`
	Count    int    `json:"count"    db:"count"`
}
