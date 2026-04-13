package auth

import (
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
)

type TTL int

const (
	AccessTokenTTL  TTL = 60 * 60      // 1 hour
	RefreshTokenTTL TTL = 60 * 60 * 12 // 12 hours
)

type MeResponse struct {
	ID         string     `json:"id"`
	Email      string     `json:"email"`
	FirstName  string     `json:"firstName"`
	LastName   string     `json:"lastName"`
	SuffixName string     `json:"suffixName,omitempty"`
	MiddleName string     `json:"middleName,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
	Role       users.Role `json:"role"`
	Type       string     `json:"type"` // "native" or "idp"
}

type IDPRefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LoginDTO struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterDTO struct {
	Email      string `json:"email"      binding:"required,email"`
	Password   string `json:"password"   binding:"required,min=8"`
	FirstName  string `json:"firstName"  binding:"required"`
	LastName   string `json:"lastName"   binding:"required"`
	SuffixName string `json:"suffixName"`
	MiddleName string `json:"middleName"`
}

type VerifyDTO struct {
	VerificationOTP string `json:"otp" binding:"required"`
}
