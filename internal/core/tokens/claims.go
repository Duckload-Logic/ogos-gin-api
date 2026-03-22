package tokens

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    string `json:"user_id"`
	UserEmail string `json:"user_email"`
	RoleID    int    `json:"role_id"`
	TokenType string `json:"token_type"` // "native" or "idp"
	jwt.RegisteredClaims
}
