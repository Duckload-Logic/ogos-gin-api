package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/tokens"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/identity/idp"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo      RepositoryInterface
	idpClient *idp.IDPClient
	redis     *datastore.RedisClient
}

// NewService creates a new auth service instance.
func NewService(
	repo RepositoryInterface,
	redis *datastore.RedisClient,
) *Service {
	return &Service{
		repo:      repo,
		idpClient: idp.NewIDPClient(),
		redis:     redis,
	}
}

// tokenService is now called inline to ensure environment variables are loaded

// AuthenticateUser handles native email/password authentication.
func (s *Service) AuthenticateUser(
	ctx context.Context, email, password string,
) (string, string, string, error) {
	// Fetch user from database (Native only)
	user, err := s.repo.GetUserByEmail(ctx, email, "native")
	if err != nil {
		return "", "", "", errors.New("invalid credentials")
	}

	log.Printf(`[AuthService:AuthenticateUser] {Get User}: %v`, user)

	// Compare hashed password
	if !user.PasswordHash.Valid {
		return "", "", "", errors.New("invalid credentials")
	}
	log.Printf(
		`[AuthService:AuthenticateUser] {Password Hash}: %v`,
		user.PasswordHash.String,
	)
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash.String),
		[]byte(password),
	)
	if err != nil {
		return "", "", "", errors.New("invalid credentials")
	}

	log.Printf(`[AuthService:AuthenticateUser] {Password Match}: %v`, err)

	// Generate the token
	token, claims, err := tokens.NewService().GenerateToken(
		user.Email,
		user.ID,
		user.RoleID,
		"",
		"native",
		constants.AccessTokenMaxAge,
	)
	if err != nil {
		return "", "", "", errors.New("failed to generate session")
	}

	// Generate refresh token
	refreshToken, _, err := tokens.NewService().GenerateToken(
		user.Email,
		user.ID,
		user.RoleID,
		"",
		"native",
		constants.RefreshTokenMaxAge,
	)
	if err != nil {
		return "", "", "", errors.New("failed to generate refresh token")
	}

	// Store in Redis using the Token ID (jti)
	err = s.storeTokenInRedis(
		ctx,
		user.ID,
		claims.ID,
		refreshToken,
		"native",
		nil,
		constants.AccessTokenMaxAge,
	)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to store token in redis: %v", err)
	}
	log.Printf(
		`[AuthService:AuthenticateUser] {Redis Store}:
		 Successfully stored token in redis`,
	)

	return user.ID, token, refreshToken, nil
}

// RefreshToken generates a new access token using a valid session handle.
func (s *Service) RefreshToken(
	ctx context.Context,
	accessTokenJTI string,
	cfg *config.Config,
) (string, string, error) {
	// 1. Get session from Redis
	key := fmt.Sprintf("%s%s", constants.RedisSessionKeyPrefix, accessTokenJTI)
	sessionData, err := s.redis.Get(ctx, key)
	if err != nil {
		return "", "", errors.New("Session expired or invalid")
	}

	var session map[string]string
	if err := json.Unmarshal([]byte(sessionData), &session); err != nil {
		return "", "", errors.New("Failed to parse session data")
	}

	refreshToken := session["appRefreshToken"]
	if refreshToken == "" {
		return "", "", errors.New("Refresh token missing from session")
	}

	// 2. Validate App Refresh Token
	claims, err := tokens.NewService().ValidateToken(refreshToken)
	if err != nil {
		return "", "", errors.New("Invalid refresh token")
	}

	// Check token type
	if claims.TokenType == "idp" {
		// Get IDP refresh token from Redis
		idpRefreshKey := fmt.Sprintf(
			"%s%s",
			constants.RedisIDPRefreshKeyPrefix,
			claims.ID,
		)
		idpRefreshToken, err := s.redis.Get(ctx, idpRefreshKey)
		if err != nil {
			return "", "", fmt.Errorf(
				"[AuthService] {Get IDP Refresh Token}: idp token missing",
			)
		}

		// Call IDP refresh endpoint
		tokenResp, err := s.idpClient.RefreshToken(ctx, idpRefreshToken, cfg)
		if err != nil {
			return "", "", fmt.Errorf("[AuthService] {IDP Refresh}: %w", err)
		}

		// Generate NEW App Tokens
		newAppAccessToken, accessClaims, err := tokens.NewService().
			GenerateToken(
				claims.UserEmail,
				claims.UserID,
				claims.RoleID,
				"",
				"idp",
				constants.AccessTokenMaxAge,
			)
		if err != nil {
			return "", "", err
		}

		newAppRefreshToken, refreshClaims, err := tokens.NewService().
			GenerateToken(
				claims.UserEmail,
				claims.UserID,
				claims.RoleID,
				"",
				"idp",
				constants.RefreshTokenMaxAge,
			)
		if err != nil {
			return "", "", err
		}

		// Update Redis: App Access session
		idpAccess := tokenResp.AccessToken
		err = s.storeTokenInRedis(
			ctx,
			claims.UserID,
			accessClaims.ID,
			newAppRefreshToken,
			"idp",
			&idpAccess,
			constants.AccessTokenMaxAge,
		)
		if err != nil {
			return "", "", err
		}

		// Update Redis: IDP Refresh linked to NEW App Refresh Token's ID
		newIdpRefreshKey := fmt.Sprintf(
			"%s%s",
			constants.RedisIDPRefreshKeyPrefix,
			refreshClaims.ID,
		)
		idpRefreshTokenToStore := tokenResp.RefreshToken
		if idpRefreshTokenToStore == "" {
			idpRefreshTokenToStore = idpRefreshToken
		}
		err = s.redis.Set(
			ctx,
			newIdpRefreshKey,
			idpRefreshTokenToStore,
			time.Duration(constants.RefreshTokenMaxAge)*time.Second,
		)
		if err != nil {
			return "", "", err
		}

		// Clean up OLD session and IDP refresh keys
		_ = s.redis.Del(ctx, key)
		_ = s.redis.Del(ctx, idpRefreshKey)

		return newAppAccessToken, newAppRefreshToken, nil
	}

	// Native flow
	newToken, newClaims, err := tokens.NewService().GenerateToken(
		claims.UserEmail,
		claims.UserID,
		claims.RoleID,
		"",
		"native",
		constants.AccessTokenMaxAge,
	)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, _, err := tokens.NewService().GenerateToken(
		claims.UserEmail,
		claims.UserID,
		claims.RoleID,
		"",
		"native",
		constants.RefreshTokenMaxAge,
	)
	if err != nil {
		return "", "", err
	}

	// Update Redis using new jti
	err = s.storeTokenInRedis(
		ctx,
		claims.UserID,
		newClaims.ID,
		newRefreshToken,
		"native",
		nil,
		constants.AccessTokenMaxAge,
	)
	if err != nil {
		return "", "", err
	}

	// Clean up OLD session key
	_ = s.redis.Del(ctx, key)

	return newToken, newRefreshToken, nil
}

// RefreshIDPToken handles token refresh for IDP-authenticated sessions.
func (s *Service) RefreshIDPToken(
	ctx context.Context, refreshToken string, cfg *config.Config,
) (string, string, error) {
	// Call IDP refresh endpoint
	tokenResp, err := s.idpClient.RefreshToken(ctx, refreshToken, cfg)
	if err != nil {
		return "", "", fmt.Errorf("[AuthService] {IDP Refresh}: %w", err)
	}

	return tokenResp.AccessToken, tokenResp.RefreshToken, nil
}

// GetMe retrieves the currently authenticated user's profile information.
func (s *Service) GetMe(
	ctx context.Context,
	userID, tokenType string,
) (*MeResponse, error) {
	// only fetch user info for native tokens
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	role, err := s.repo.GetRoleByID(ctx, user.RoleID)
	if err != nil {
		// Fallback if role not found
		return &MeResponse{
			ID:         user.ID,
			Email:      user.Email,
			FirstName:  user.FirstName,
			LastName:   user.LastName,
			MiddleName: user.MiddleName.String,
			Roles:      []string{"user"},
			Type:       tokenType,
		}, nil
	}

	return &MeResponse{
		ID:         user.ID,
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: user.MiddleName.String,
		Roles:      []string{role.Name},
		Type:       tokenType,
	}, nil
}

// ParseIDPRoles sanitizes role strings received from the IDP.
func (s *Service) ParseIDPRoles(roles []string) []string {
	parsedRoles := make([]string, 0, len(roles))
	for _, role := range roles {
		// Split by ":" and take the last part
		parts := strings.Split(role, ":")
		if len(parts) > 1 {
			parsedRoles = append(parsedRoles, parts[len(parts)-1])
		} else {
			parsedRoles = append(parsedRoles, role)
		}
	}
	return parsedRoles
}

func (s *Service) storeTokenInRedis(
	ctx context.Context,
	userID, jti, appRefreshToken, tokenType string,
	idpAccessToken *string,
	expiryMinutes int,
) error {
	key := fmt.Sprintf("%s%s", constants.RedisSessionKeyPrefix, jti)
	val := map[string]string{
		"userID":          userID,
		"tokenType":       tokenType,
		"appRefreshToken": appRefreshToken,
	}
	if idpAccessToken != nil {
		val["idpAccessToken"] = *idpAccessToken
	}
	valJSON, _ := json.Marshal(val)

	err := s.redis.Set(
		ctx,
		key,
		valJSON,
		time.Duration(expiryMinutes)*time.Minute,
	)
	if err != nil {
		return fmt.Errorf("failed to store token in redis: %v", err)
	}

	return nil
}

// Logout invalidates the user's session in Redis and optionally the IDP.
func (s *Service) Logout(
	ctx context.Context,
	token string,
	tokenType string,
	cfg *config.Config,
) error {
	// 1. Identify the session using the Access Token JTI
	claims, err := tokens.NewService().ParseTokenUnverified(token)
	if err != nil {
		log.Printf("[AuthService:Logout] {Parse Error}: %v", err)
		return nil // Move on since we can't identify the session
	}
	accessJTI := claims.ID

	// Fetch the session data to find linked refresh tokens
	sessionKey := fmt.Sprintf(
		"%s%s",
		constants.RedisSessionKeyPrefix,
		accessJTI,
	)
	sessionData, _ := s.redis.Get(ctx, sessionKey)

	// Delete any linked IDP refresh tokens
	if sessionData != "" {
		var session map[string]string
		if err := json.Unmarshal([]byte(sessionData), &session); err == nil {
			if appRefreshToken := session["appRefreshToken"]; appRefreshToken != "" {
				// Get refresh token claims to identify IDP refresh key
				rClaims, err := tokens.NewService().
					ParseTokenUnverified(appRefreshToken)
				if err == nil {
					idpKey := fmt.Sprintf(
						"%s%s",
						constants.RedisIDPRefreshKeyPrefix,
						rClaims.ID,
					)
					_ = s.redis.Del(ctx, idpKey)
				}
			}
		}
	}

	// 4. Delete the primary session key
	if err := s.redis.Del(ctx, sessionKey); err != nil {
		log.Printf("[AuthService:Logout] {Redis Error}: %v", err)
	}

	// 5. Return the IDP logout URL if appropriate
	if tokenType == "idp" {
		log.Println("[AuthService:Logout] {IDP Logout}: Call IDP logout")
		_, err = s.idpClient.Logout(ctx, cfg)
		if err != nil {
			log.Printf("[AuthService:Logout] {IDP Logout Error}: %v", err)
			return err
		}

		log.Println("[AuthService:Logout] {IDP Logout}: IDP logout successful")
		return nil
	}

	return nil
}

// IDP integration methods

// GetAuthorizeURL generates the complete OAuth 2.0 authorization URL
// with PKCE parameters. This method creates a state parameter for CSRF
// protection, generates PKCE verifier and challenge, stores the state
// with metadata, and builds the authorization URL.
//
// Parameters:
//   - cfg: Application configuration containing IDP endpoints
//
// Returns the authorization URL and state parameter, or an error if
// generation fails.
func (s *Service) GetAuthorizeURL(
	cfg *config.Config,
) (string, error) {
	// Build authorization URL with all required parameters
	params := url.Values{}
	params.Set("client_id", cfg.IDPClientID)

	authURL := fmt.Sprintf(
		"%s?%s",
		cfg.IDPLoginURL,
		params.Encode(),
	)

	return authURL, nil
}

// PostIDPTokenExchange orchestrates the complete IDP login flow:
// validates state, exchanges code for token, fetches user info,
// provisions user, and generates application JWT tokens.
//
// Parameters:
//   - ctx: Context for database and HTTP operations
//   - code: Authorization code from IDP callback
//   - state: State parameter from IDP callback
//   - cfg: Application configuration
//
// Returns user ID and JWT tokens, or an error if any step fails.
func (s *Service) PostIDPTokenExchange(
	ctx context.Context,
	code string,
	cfg *config.Config,
) (string, string, string, string, string, error) {
	// Exchange authorization code for IDP tokens
	tokenResp, err := s.idpClient.ExchangeCodeForToken(ctx, code, cfg)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf(
			"[AuthService] {Token Exchange}: %w",
			err,
		)
	}

	// Fetch User Info from IDP
	userInfo, err := s.GetIDPUserInfo(ctx, tokenResp.AccessToken, cfg)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf(
			"[AuthService] {Fetch User Info}: %w",
			err,
		)
	}

	// Parse Tokens
	idpAccessToken := tokenResp.AccessToken
	idpRefreshToken := tokenResp.RefreshToken

	// Dynamic Role Mapping from IDP Tags
	// Logic: Identify tag:admin, tag:student, tag:superadmin at runtime
	appRoleID := s.mapIDPRolesToInternalID(userInfo.Roles)

	// Upsert IDP user into native database
	// This ensures non-existing users are added and existing ones
	// are synchronized
	err = datastore.RunInTransaction(
		ctx,
		s.repo.(*users.Repository).GetDB(),
		func(tx datastore.DB) error {
			return s.repo.CreateUser(ctx, tx, users.User{
				ID:           userInfo.ID,
				RoleID:       appRoleID,
				FirstName:    userInfo.FirstName,
				LastName:     userInfo.LastName,
				Email:        userInfo.Email,
				AuthType:     "idp",
				PasswordHash: sql.NullString{Valid: false},
				IsActive:     1,
			})
		},
	)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf(
			"[AuthService] {Provision IDP User}: %w",
			err,
		)
	}

	appUserID := userInfo.ID

	// Map Role ID to Name for frontend redirection
	role, err := s.repo.GetRoleByID(ctx, appRoleID)
	roleName := "student"
	if err == nil && role != nil {
		roleName = role.Name
	}

	// Generate internal App Tokens using the actual app IDs
	appAccessToken, accessClaims, err := tokens.NewService().GenerateToken(
		userInfo.Email,
		appUserID,
		appRoleID,
		"",
		"idp",
		constants.AccessTokenMaxAge/60,
	)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf(
			"[AuthService] {Generate App Access Token}: %w",
			err,
		)
	}

	appRefreshToken, refreshClaims, err := tokens.NewService().GenerateToken(
		userInfo.Email,
		appUserID,
		appRoleID,
		"",
		"idp",
		constants.RefreshTokenMaxAge/60,
	)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf(
			"[AuthService] {Generate App Refresh Token}: %w",
			err,
		)
	}

	// Store App Access Token in Redis using its ID (jti)
	if err := s.storeTokenInRedis(
		ctx,
		appUserID,
		accessClaims.ID,
		appRefreshToken,
		"idp",
		&idpAccessToken,
		constants.AccessTokenMaxAge/60,
	); err != nil {
		return "", "", "", "", "", fmt.Errorf(
			"[AuthService] {Store Access in Redis}: %w",
			err,
		)
	}

	// Store IDP Refresh Token in Redis associated with the App
	// Refresh Token's ID (jti)
	idpRefreshKey := fmt.Sprintf(
		"%s%s",
		constants.RedisIDPRefreshKeyPrefix,
		refreshClaims.ID,
	)
	idpRefreshTokenToStore := idpRefreshToken
	if idpRefreshTokenToStore == "" {
		// This shouldn't happen on initial login, but good for
		// robustness
	}
	err = s.redis.Set(
		ctx,
		idpRefreshKey,
		idpRefreshTokenToStore,
		time.Duration(constants.RefreshTokenMaxAge)*time.Second,
	)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf(
			"[AuthService] {Store IDP Refresh in Redis}: %w",
			err,
		)
	}

	return appAccessToken, appRefreshToken, appUserID, userInfo.Email,
		roleName, nil
}

// GetIDPUserInfo fetches user information from the IDP userinfo endpoint
// using the provided access token. This is typically called after a
// successful token exchange to retrieve user details for provisioning.
//
// Parameters:
//   - ctx: Context for the HTTP request
//   - accessToken: Access token obtained from IDP token exchange
//   - cfg: Application configuration containing IDP endpoints
//
// Returns the IDP user information or an error if retrieval fails.
func (s *Service) GetIDPUserInfo(
	ctx context.Context,
	accessToken string,
	cfg *config.Config,
) (*idp.IDPUserInfo, error) {
	userInfo, err := s.idpClient.GetUserInfo(ctx, accessToken, cfg)
	if err != nil {
		return nil, fmt.Errorf(
			"[AuthService] {Get IDP User Info}: %w",
			err,
		)
	}
	return userInfo, nil
}

// mapIDPRolesToInternalID translates IDP tags to internal role IDs.
// Tags format: tag:student, tag:admin, tag:superadmin
func (s *Service) mapIDPRolesToInternalID(roles []string) int {
	// Priority order: superadmin > admin > student
	hasAdmin := false
	hasSuper := false
	hasStudent := false

	for _, r := range roles {
		if r == "" {
			continue
		}

		// Safely split and check for parts to avoid panics on missing colon
		parts := strings.Split(r, ":")
		rolePart := ""
		if len(parts) > 1 {
			rolePart = strings.ToLower(parts[1])
		} else {
			rolePart = strings.ToLower(parts[0])
		}

		switch rolePart {
		case "superadmin":
			hasSuper = true
		case "admin":
			hasAdmin = true
		case "student":
			hasStudent = true
		}
	}

	if hasSuper {
		return int(constants.SuperAdminRoleID)
	}
	if hasAdmin {
		return int(constants.CounselorRoleID)
	}
	if hasStudent {
		return int(constants.StudentRoleID)
	}

	return int(constants.StudentRoleID) // Default to Student
}

// ValidateIDPSession checks if the provided session ID is valid on the IDP.
func (s *Service) ValidateIDPSession(
	ctx context.Context,
	sessionID string,
	cfg *config.Config,
) (*idp.IDPSessionResponse, error) {
	return s.idpClient.ValidateSession(ctx, sessionID, cfg)
}

func (s *Service) BlockUser(ctx context.Context, userID string) error {
	return datastore.RunInTransaction(
		ctx,
		s.repo.(*users.Repository).GetDB(),
		func(tx datastore.DB) error {
			return s.repo.BlockUser(ctx, tx, userID)
		},
	)
}

func (s *Service) UnblockUser(ctx context.Context, userID string) error {
	return datastore.RunInTransaction(
		ctx,
		s.repo.(*users.Repository).GetDB(),
		func(tx datastore.DB) error {
			return s.repo.UnblockUser(ctx, tx, userID)
		},
	)
}
