package users

import "database/sql"

type Role struct {
	ID   int    `db:"id"   json:"id"`
	Name string `db:"name" json:"name"`
}

type User struct {
	ID           string         `db:"id"            json:"id"`
	RoleID       int            `db:"role_id"       json:"role_id"`
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
