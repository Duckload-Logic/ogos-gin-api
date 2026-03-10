package users

import "database/sql"

type Role struct {
	ID   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type User struct {
	ID           int            `db:"id" json:"id"`
	RoleID       int            `db:"role_id" json:"role_id"`
	FirstName    string         `db:"first_name" json:"first_name"`
	MiddleName   sql.NullString `db:"middle_name" json:"middle_name"`
	LastName     string         `db:"last_name" json:"last_name"`
	Email        string         `db:"email" json:"email"`
	PasswordHash string         `db:"password_hash" json:"-"`
	CreatedAt    sql.NullTime   `db:"created_at" json:"created_at"`
	UpdatedAt    sql.NullTime   `db:"updated_at" json:"updated_at"`
}
