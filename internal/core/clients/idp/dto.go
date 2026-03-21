package idp

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
	ID         string   `json:"id"`
	Email      string   `json:"email"`
	FirstName  string   `json:"first_name"`
	LastName   string   `json:"last_name"`
	MiddleName string   `json:"middle_name,omitempty"`
	Roles      []string `json:"roles,omitempty"`
}
