package users

import (
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

type Role struct {
	ID   int
	Name string
}

type ProfilePicture struct {
	FileID string
	UserID string
}

type User struct {
	ID           string
	RoleID       int
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
