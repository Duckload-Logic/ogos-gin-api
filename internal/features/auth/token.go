package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

type TokenService struct {
	secret []byte
}

func NewTokenService() *TokenService {
	return &TokenService{secret: []byte(os.Getenv("JWT_SECRET"))}
}

func (s *TokenService) GenerateToken(
	userID, roleID int, roleName string, expireMinutes int,
) (string, error) {
	// Create the JWT claims
	claims := &middleware.Claims{
		UserID: userID,
		RoleID: roleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(
				time.Duration(expireMinutes) * time.Minute),
			),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer:   "pupt-ogos-api",
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *TokenService) ValidateToken(tokenString string) (
	*middleware.Claims, error,
) {
	// Parse the token
	token, err := jwt.ParseWithClaims(
		tokenString,
		&middleware.Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return s.secret, nil
		},
	)
	if err != nil {
		return nil, err
	}

	// Validate the claims
	claims, ok := token.Claims.(*middleware.Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil

}
