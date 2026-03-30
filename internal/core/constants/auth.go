package constants

import "time"

// OAuth 2.0 parameter constants for IDP integration
const (
	// ResponseTypeCode is the OAuth 2.0 response_type parameter value
	// for authorization code flow
	ResponseTypeCode = "code"
)

// Timeout constants for IDP requests
const (
	// IDPRequestTimeout is the timeout duration for IDP HTTP requests
	// (10 seconds)
	IDPRequestTimeout = 30 * time.Second
)

// Cookie configuration constants
const (
	// AccessTokenCookieName is the name of the access token cookie
	AccessTokenCookieName = "access_token"

	// RefreshTokenCookieName is the name of the refresh token cookie
	RefreshTokenCookieName = "refresh_token"

	// AccessTokenMaxAge is the maximum age in seconds for access token
	// cookie (30 minutes = 1800 seconds)
	AccessTokenMaxAge = 1800

	// RefreshTokenMaxAge is the maximum age in seconds for refresh token
	// cookie (12 hours = 43200 seconds)
	RefreshTokenMaxAge = 43200

	// CookiePathRoot sets cookies to be accessible from root path
	CookiePathRoot = "/"
)

// Logging constants for consistent log messages
const (
	// LogCategorySecurity is the log category for security-related events
	LogCategorySecurity = "Security"

	// LogActionLoginSuccess is the log action for successful login
	LogActionLoginSuccess = "LoginSuccess"

	// LogActionLoginFailed is the log action for failed login
	LogActionLoginFailed = "LoginFailed"
)

// Redis Key Constants
const (
	// RedisSessionKeyPrefix is the prefix for session keys (session:jti)
	RedisSessionKeyPrefix = "session:"

	// RedisIDPRefreshKeyPrefix is the prefix for IDP refresh tokens
	// (idp_refresh:jti)
	RedisIDPRefreshKeyPrefix = "idp_refresh:"
)
