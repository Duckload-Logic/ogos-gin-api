package constants

type ErrorMessage string

const (
	ErrNotFound            ErrorMessage = "Resource not found"
	ErrInternalServerError ErrorMessage = "Something went wrong :<"
	ErrInvalidRequest      ErrorMessage = "Invalid request"

	// IDP

	// ErrInvalidState indicates the state parameter is invalid
	// or expired
	ErrInvalidState = "Invalid state parameter"

	// ErrTokenExchangeFailed indicates token exchange with IDP failed
	ErrTokenExchangeFailed = "Token exchange failed"

	// ErrUserInfoFailed indicates user info retrieval from IDP failed
	ErrUserInfoFailed = "Failed to retrieve user info"

	// ErrUserProvisioningFailed indicates user provisioning failed
	ErrUserProvisioningFailed = "User provisioning failed"

	// ErrInvalidVerifier indicates code_verifier is invalid
	ErrInvalidVerifier = "invalid code_verifier"

	// ErrInvalidChallenge indicates code_challenge is invalid
	ErrInvalidChallenge = "invalid code_challenge"

	// ErrMissingCode indicates authorization code is missing
	ErrMissingCode = "authorization code is required"

	// ErrMissingState indicates state parameter is missing
	ErrMissingState = "state parameter is required"

	// ErrInvalidRedirectURI indicates redirect URI is invalid
	ErrInvalidRedirectURI = "invalid redirect URI"

	// ErrIDPError indicates IDP returned an error
	ErrIDPError = "IDP returned an error"

	// ErrInvalidResponse indicates IDP response is invalid
	ErrInvalidResponse = "invalid response from IDP"
)
