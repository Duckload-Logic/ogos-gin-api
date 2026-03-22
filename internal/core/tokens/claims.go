package tokens

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    string
	UserEmail string
	RoleID    int
	TokenType string // "native" or "idp"
	jwt.RegisteredClaims
}
