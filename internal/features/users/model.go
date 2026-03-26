package users

import "database/sql"

type Role struct {
	ID   int    `db:"id"   json:"id"`
	Name string `db:"name" json:"name"`
}

type User struct {
	ID           string         `db:"id"            json:"id"`
	RoleID       int            `db:"role_id"       json:"role_id"`
	FirstName    string         `db:"first_name"    json:"first_name"`
	MiddleName   sql.NullString `db:"middle_name"   json:"middle_name"`
	LastName     string         `db:"last_name"     json:"last_name"`
	Email        string         `db:"email"         json:"email"`
	PasswordHash sql.NullString `db:"password_hash" json:"-"`
	AuthType     string         `db:"auth_type"     json:"authType"`
	IsActive     int            `db:"is_active"     json:"isActive"`
	CreatedAt    sql.NullTime   `db:"created_at"    json:"created_at"`
	UpdatedAt    sql.NullTime   `db:"updated_at"    json:"updated_at"`
}
