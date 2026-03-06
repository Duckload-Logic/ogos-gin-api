package tokens

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserEmail string
	RoleID    int
	jwt.RegisteredClaims
}
