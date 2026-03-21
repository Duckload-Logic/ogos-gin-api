package tokens

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWTSecretEnv = os.Getenv("JWT_SECRET")

type Service struct {
	secret []byte
}

func NewService() *Service {
	return &Service{secret: []byte(JWTSecretEnv)}
}

func (s *Service) GenerateToken(
	userEmail string, userID string, roleID int, roleName string, expireMinutes int,
) (string, error) {
	claims := &Claims{
		UserEmail: userEmail,
		UserID:    userID,
		RoleID:    roleID,
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

func (s *Service) ValidateToken(tokenString string) (
	*Claims, error,
) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return s.secret, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

func (s *Service) InvalidateToken(tokenString string) error {
	// TODO: Implement token invalidation (e.g., using a blacklist or token store)
	return nil
}
