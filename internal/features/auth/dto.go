package auth

type TTL int

const (
	AccessTokenTTL  TTL = 60 * 30
	RefreshTokenTTL TTL = 60 * 60 * 12
)

type IDPUser struct {
	Email string   `json:"email"`
	Name  string   `json:"name"`
	Roles []string `json:"roles"` 
}

type ExchangeRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
}

type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type TokenDTO struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}