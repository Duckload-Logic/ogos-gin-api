package auth

import "time"

type TTL int

const (
	AccessTokenTTL  TTL = 60 * 60      // 1 hour
	RefreshTokenTTL TTL = 60 * 60 * 12 // 12 hours
)

type MeResponse struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	MiddleName string    `json:"middleName,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	Roles      []string  `json:"roles"`
	Type       string    `json:"type"` // "native" or "idp"
}

type IDPRefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LoginDTO struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type TokenDTO struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
