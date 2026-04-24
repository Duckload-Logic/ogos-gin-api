package m2mclients

import (
	"time"
)

type M2MClient struct {
	ID           string    `db:"id"            json:"id"`
	UserID       string    `db:"user_id"       json:"userId"`
	ClientName   string    `db:"client_name"   json:"clientName"`
	ClientID     string    `db:"client_id"     json:"clientId"`
	ClientSecret string    `db:"client_secret" json:"-"`
	IsActive     bool      `db:"is_active"     json:"isActive"`
	CreatedAt    time.Time `db:"created_at"    json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at"    json:"updatedAt"`
}
