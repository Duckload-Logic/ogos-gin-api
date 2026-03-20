package auth

type TTL int

const (
	AccessTokenTTL  TTL = 60 * 30
	RefreshTokenTTL TTL = 60 * 60 * 12
)

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

type IDPAuthorizeURLResponse struct {
	AuthorizationURL string `json:"authorizationUrl"`
}

type IDPLoginRequest struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token"` // optional, but we may need it
}

type IDPTokenExchangeRequest struct {
	Code string `json:"code" binding:"required"`
}

// IDPTokenResponse represents the response from IDP token endpoint
type IDPTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// IDPUserInfo represents user information from IDP userinfo endpoint
type IDPUserInfo struct {
	Email      string `json:"email"`
	GivenName  string `json:"given_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name,omitempty"`
}
