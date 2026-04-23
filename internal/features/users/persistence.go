package users

import (
	"database/sql"
)

type RoleDB struct {
	ID   int    `db:"id"   json:"id"`
	Name string `db:"name" json:"name"`
}

type ProfilePictureDB struct {
	FileID string `db:"file_id" json:"file_id"`
	UserID string `db:"user_id" json:"user_id"`
}

type UserDB struct {
	ID           string         `db:"id"            json:"id"`
	FirstName    string         `db:"first_name"    json:"firstName"`
	MiddleName   sql.NullString `db:"middle_name"   json:"middleName"`
	LastName     string         `db:"last_name"     json:"lastName"`
	SuffixName   sql.NullString `db:"suffix_name"   json:"suffixName"`
	Email        string         `db:"email"         json:"email"`
	PasswordHash sql.NullString `db:"password_hash" json:"passwordHash"`
	AuthType     string         `db:"auth_type"     json:"authType"`
	IsActive     int            `db:"is_active"     json:"isActive"`
	CreatedAt    sql.NullTime   `db:"created_at"    json:"createdAt"`
	UpdatedAt    sql.NullTime   `db:"updated_at"    json:"updatedAt"`
}

type UserRoleDB struct {
	UserID      string         `db:"user_id"`
	RoleID      int            `db:"role_id"`
	AssignedAt  sql.NullTime   `db:"assigned_at"`
	AssignedBy  sql.NullString `db:"assigned_by"`
	Reason      sql.NullString `db:"reason"`
	ReferenceID sql.NullString `db:"reference_id"`
}
