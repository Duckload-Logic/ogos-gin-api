package users

import (
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type RoleAssignment struct {
	UserID      string
	RoleID      int
	AssignedAt  structs.NullableTime
	AssignedBy  structs.NullableString
	Reason      structs.NullableString
	ReferenceID structs.NullableString
}

type ProfilePicture struct {
	FileID string
	UserID string
}

type User struct {
	ID           string
	Roles        []Role
	FirstName    string
	MiddleName   structs.NullableString
	LastName     string
	SuffixName   structs.NullableString
	Email        string
	PasswordHash structs.NullableString
	AuthType     string
	IsActive     bool
	CreatedAt    structs.NullableTime
	UpdatedAt    structs.NullableTime
}
