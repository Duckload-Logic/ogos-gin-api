package tokens

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    string `json:"userId"`
	UserEmail string `json:"userEmail"`
	RoleID    int    `json:"roleId"`
	TokenType string `json:"tokenType"` // "native" or "idp"
	jwt.RegisteredClaims
}
