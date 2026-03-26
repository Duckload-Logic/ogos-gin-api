package tokens

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Service struct {
	secret []byte
}

func NewService() *Service {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Printf(
			"[WARNING] {NewService}: JWT_SECRET is empty. Signing will fail.",
		)
	}
	return &Service{secret: []byte(secret)}
}

func (s *Service) GenerateToken(
	userEmail string,
	userID string,
	roleID int,
	roleName string,
	tokenType string,
	expireMinutes int,
) (string, *Claims, error) {
	claims := &Claims{
		UserEmail: userEmail,
		UserID:    userID,
		RoleID:    roleID,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(
				time.Duration(expireMinutes) * time.Minute),
			),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer:   "pupt-ogos-api",
			ID:       uuid.New().String(),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	return signed, claims, err
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
