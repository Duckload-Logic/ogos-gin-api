package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/sessions"
	"github.com/olazo-johnalbert/duckload-api/internal/core/tokens"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/email"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/identity/idp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo           RepositoryInterface
	idpClient      *idp.IDPClient
	redis          *datastore.RedisClient
	sessionService *sessions.Service
	emailer        email.Emailer
}

func NewService(
	repo RepositoryInterface,
	redis *datastore.RedisClient,
	sessionService *sessions.Service,
	emailer email.Emailer,
) *Service {
	return &Service{
		repo:           repo,
		idpClient:      idp.NewIDPClient(),
		redis:          redis,
		sessionService: sessionService,
		emailer:        emailer,
	}
}

// RegisterUser handles native user registration.
func (s *Service) RegisterUser(
	ctx context.Context,
	req RegisterDTO,
) (string, error) {
	// Check if user already exists
	existingUser, _ := s.repo.GetUserByEmail(
		ctx,
		req.Email,
		string(constants.AuthTypeNative),
	)
	if existingUser != nil {
		return "", errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}

	transactionID := uuid.NewString()

	// Create user object
	user := users.User{
		ID:        uuid.NewString(),
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		MiddleName: sql.NullString{
			String: req.MiddleName,
			Valid:  req.MiddleName != "",
		},
		SuffixName: sql.NullString{
			String: req.SuffixName,
			Valid:  req.SuffixName != "",
		},
		PasswordHash: sql.NullString{
			String: string(hashedPassword),
			Valid:  true,
		},
		RoleID:   int(constants.DeveloperRoleID),
		AuthType: string(constants.AuthTypeNative),
		IsActive: 0,
	}

	verificationOTP, err := s.get6DigitOTP()
	if err != nil {
		return "", fmt.Errorf("failed to generate verification token: %v", err)
	}
	hashedOTP, err := bcrypt.GenerateFromPassword(
		[]byte(verificationOTP),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate verification token: %v", err)
	}

	// Use a map for storage to include the hidden PasswordHash
	storageData := map[string]interface{}{
		"user":         user,
		"passwordHash": user.PasswordHash,
	}

	userJSON, err := json.Marshal(storageData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal user: %v", err)
	}

	// Store in Redis using the UUID (transactionID)
	val := map[string]string{
		"registrationID":    transactionID,
		"verificationToken": string(hashedOTP),
		"user":              string(userJSON),
	}
	err = s.sessionService.StoreToken(
		ctx,
		sessions.NewJTI(transactionID),
		val,
		300,
	)
	if err != nil {
		return "", fmt.Errorf("failed to store token in redis: %v", err)
	}

	if err := s.validateEmailDomain(req.Email); err != nil {
		return "", err
	}

	err = s.sendVerificationEmail(ctx, verificationOTP, req.Email)
	if err != nil {
		return "", err
	}

	return transactionID, nil
}

func (s *Service) ResendVerification(
	ctx context.Context,
	registrationID string,
) error {
	val, err := s.sessionService.GetToken(ctx, sessions.NewJTI(registrationID))
	if err != nil {
		return fmt.Errorf("failed to get token from redis: %v", err)
	}

	verificationOTP, err := s.get6DigitOTP()
	if err != nil {
		return fmt.Errorf("failed to generate verification token: %v", err)
	}

	hashedOTP, err := bcrypt.GenerateFromPassword(
		[]byte(verificationOTP),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf("failed to generate verification token: %v", err)
	}

	val["verificationToken"] = string(hashedOTP)
	err = s.sessionService.StoreToken(
		ctx,
		sessions.NewJTI(registrationID),
		val,
		300,
	)
	if err != nil {
		return fmt.Errorf("failed to store token in redis: %v", err)
	}

	var storageData struct {
		User users.User `json:"user"`
	}
	err = json.Unmarshal([]byte(val["user"]), &storageData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user: %v", err)
	}
	user := storageData.User

	if err := s.validateEmailDomain(user.Email); err != nil {
		return err
	}

	err = s.sendVerificationEmail(ctx, verificationOTP, user.Email)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) get6DigitOTP() (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      constants.ClaimsIssuer,
		AccountName: constants.FromEmail(),
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate verification token: %v", err)
	}

	verificationOTP, err := totp.GenerateCode(key.Secret(), time.Now())
	if err != nil {
		return "", fmt.Errorf("failed to generate verification token: %v", err)
	}

	return verificationOTP, nil
}

func (s *Service) validateEmailDomain(email string) error {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return errors.New("invalid email format")
	}
	domain := parts[1]
	mx, err := net.LookupMX(domain)
	if err != nil || len(mx) == 0 {
		return fmt.Errorf("email domain %s is unreachable", domain)
	}
	return nil
}

func (s *Service) sendVerificationEmail(
	ctx context.Context,
	verificationOTP, receiverEmail string,
) error {
	isSent, err := s.emailer.SendOTP(ctx, receiverEmail, verificationOTP)
	if err != nil {
		return fmt.Errorf("failed to send verification email: %v", err)
	}
	if !isSent {
		return errors.New("failed to send verification email")
	}

	return nil
}

func (s *Service) VerifyUser(
	ctx context.Context,
	registrationID string,
	verificationOTP string,
) (string, string, error) {
	// Get user from Redis
	val, err := s.sessionService.GetToken(ctx, sessions.NewJTI(registrationID))
	if err != nil {
		return "", "", fmt.Errorf("failed to get token from redis: %v", err)
	}

	// Verify OTP
	err = bcrypt.CompareHashAndPassword(
		[]byte(val["verificationToken"]),
		[]byte(verificationOTP),
	)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	var storageData struct {
		User         users.User     `json:"user"`
		PasswordHash sql.NullString `json:"passwordHash"`
	}

	err = json.Unmarshal([]byte(val["user"]), &storageData)
	if err != nil {
		return "", "", fmt.Errorf("failed to unmarshal user data: %v", err)
	}

	user := storageData.User
	user.PasswordHash = storageData.PasswordHash

	user.IsActive = 1

	// Run in transaction
	err = datastore.RunInTransaction(
		ctx,
		s.repo.(*users.Repository).GetDB(),
		func(tx datastore.DB) error {
			return s.repo.CreateUser(ctx, tx, user)
		},
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to create user: %v", err)
	}

	return user.ID, user.Email, nil
}

// AuthenticateUser handles native email/password authentication.
func (s *Service) AuthenticateUser(
	ctx context.Context, email, password, ipAddress, userAgent string,
) (string, string, string, error) {
	// Fetch user from database (Native only)
	user, err := s.repo.GetUserByEmail(
		ctx,
		email,
		string(constants.AuthTypeNative),
	)
	if err != nil {
		return "", "", "", errors.New("Invalid credentials")
	}

	if user.IsActive == 0 {
		return "", "", "", errors.New("User is not active")
	}

	// Compare hashed password
	if !user.PasswordHash.Valid {
		return "", "", "", errors.New("Invalid credentials")
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash.String),
		[]byte(password),
	)
	if err != nil {
		return "", "", "", errors.New("Invalid credentials")
	}

	// Generate the token
	token, claims, err := tokens.NewService().GenerateToken(
		user.Email,
		user.ID,
		user.RoleID,
		"",
		string(constants.AuthTypeNative),
		constants.AccessTokenMaxAge,
	)
	if err != nil {
		return "", "", "", errors.New("Failed to generate session")
	}

	// Generate refresh token
	refreshToken, _, err := tokens.NewService().GenerateToken(
		user.Email,
		user.ID,
		user.RoleID,
		"",
		string(constants.AuthTypeNative),
		constants.RefreshTokenMaxAge,
	)
	if err != nil {
		return "", "", "", errors.New("Failed to generate refresh token")
	}

	// Store in Redis using the Token ID (jti)
	val := map[string]string{
		"userID":          user.ID,
		"tokenType":       string(constants.AuthTypeNative),
		"appRefreshToken": refreshToken,
		"ipAddress":       ipAddress,
		"userAgent":       userAgent,
	}
	err = s.sessionService.StoreUserToken(
		ctx,
		user.ID,
		sessions.NewJTI(claims.ID),
		val,
		constants.RefreshTokenMaxAge,
	)
	if err != nil {
		return "", "", "", fmt.Errorf("Failed to store token in redis: %v", err)
	}

	return user.ID, token, refreshToken, nil
}

// RefreshToken generates a new access token using a valid session handle.
func (s *Service) RefreshToken(
	ctx context.Context,
	accessTokenJTI sessions.JTIDTO,
	cfg *config.Config,
	ipAddress, userAgent string,
) (string, string, error) {
	// Get session from Redis
	session, err := s.sessionService.GetToken(ctx, accessTokenJTI)
	if err != nil {
		return "", "", errors.New("Session expired or invalid")
	}

	refreshToken := session["appRefreshToken"]
	if refreshToken == "" {
		return "", "", errors.New("Refresh token missing from session")
	}

	// Validate App Refresh Token
	claims, err := tokens.NewService().ValidateToken(refreshToken)
	if err != nil {
		return "", "", errors.New("Invalid refresh token")
	}

	// Check token type
	if claims.TokenType == string(constants.AuthTypeIDP) {
		// Get IDP refresh token from Redis
		idpRefreshKey := sessions.NewJTI(claims.ID).ToIDPRefreshKey()
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
				string(constants.AuthTypeIDP),
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
				string(constants.AuthTypeIDP),
				constants.RefreshTokenMaxAge,
			)
		if err != nil {
			return "", "", err
		}

		// Update Redis: App Access session
		val := map[string]string{
			"userID":          claims.UserID,
			"tokenType":       string(constants.AuthTypeIDP),
			"appRefreshToken": newAppRefreshToken,
			"idpAccessToken":  tokenResp.AccessToken,
		}
		err = s.sessionService.StoreToken(
			ctx,
			sessions.NewJTI(accessClaims.ID),
			val,
			constants.RefreshTokenMaxAge,
		)
		if err != nil {
			return "", "", err
		}

		// Update Redis: IDP Refresh linked to NEW App Refresh Token's ID
		newIdpRefreshKey := sessions.NewJTI(refreshClaims.ID).ToIDPRefreshKey()
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
		_ = s.sessionService.DeleteUserToken(ctx, claims.UserID, accessTokenJTI)
		_ = s.redis.Del(ctx, idpRefreshKey)

		return newAppAccessToken, newAppRefreshToken, nil
	}

	// Native flow
	newToken, newClaims, err := tokens.NewService().GenerateToken(
		claims.UserEmail,
		claims.UserID,
		claims.RoleID,
		"",
		string(constants.AuthTypeNative),
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
		string(constants.AuthTypeNative),
		constants.RefreshTokenMaxAge,
	)
	if err != nil {
		return "", "", err
	}

	// Update Redis using new jti
	val := map[string]string{
		"userID":          claims.UserID,
		"tokenType":       string(constants.AuthTypeNative),
		"appRefreshToken": newRefreshToken,
		"ipAddress":       ipAddress,
		"userAgent":       userAgent,
	}
	err = s.sessionService.StoreUserToken(
		ctx,
		claims.UserID,
		sessions.NewJTI(newClaims.ID),
		val,
		constants.RefreshTokenMaxAge,
	)
	if err != nil {
		return "", "", err
	}

	// Clean up OLD session key
	_ = s.sessionService.DeleteUserToken(ctx, claims.UserID, accessTokenJTI)

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
			SuffixName: user.SuffixName.String,
			MiddleName: user.MiddleName.String,
			CreatedAt:  user.CreatedAt.Time,
			Role:       users.Role{},
			Type:       tokenType,
		}, nil
	}

	return &MeResponse{
		ID:         user.ID,
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		SuffixName: user.SuffixName.String,
		MiddleName: user.MiddleName.String,
		CreatedAt:  user.CreatedAt.Time,
		Role:       *role,
		Type:       tokenType,
	}, nil
}

// Logout invalidates the user's session in Redis and optionally the IDP.
func (s *Service) Logout(
	ctx context.Context,
	token string,
	tokenType string,
	cfg *config.Config,
) (string, error) {
	// Identify the session using the Access Token JTI
	claims, err := tokens.NewService().ParseTokenUnverified(token)
	if err != nil {
		log.Printf("[AuthService:Logout] {Parse Error}: %v", err)
		return "", nil // Move on since we can't identify the session
	}
	accessJTI := claims.ID

	// Fetch the session data to find linked refresh tokens
	sessionData, _ := s.sessionService.GetToken(ctx, sessions.NewJTI(accessJTI))
	var idpToken string
	if sessionData != nil {
		idpToken = sessionData["idpAccessToken"]
	}

	// Delete any linked IDP refresh tokens
	if sessionData != nil {
		if appRefreshToken := sessionData["appRefreshToken"]; appRefreshToken != "" {
			// Get refresh token claims to identify IDP refresh key
			rClaims, err := tokens.NewService().
				ParseTokenUnverified(appRefreshToken)
			if err == nil {
				idpKey := sessions.NewJTI(rClaims.ID).ToIDPRefreshKey()
				_ = s.redis.Del(ctx, idpKey)
			}
		}
	}

	// Delete the primary session key
	if userID := claims.UserID; userID != "" {
		if err := s.sessionService.DeleteUserToken(ctx, userID, sessions.NewJTI(accessJTI)); err != nil {
			log.Printf("[AuthService:Logout] {Redis Error}: %v", err)
		}
	} else {
		// Fallback for M2M or cases where userID missing
		if err := s.sessionService.DeleteToken(ctx, sessions.NewJTI(accessJTI)); err != nil {
			log.Printf("[AuthService:Logout] {Redis Error}: %v", err)
		}
	}

	userInfo, _ := s.idpClient.GetUserInfo(ctx, idpToken, cfg)

	// Generate IDP logout URL if IDP token
	if tokenType == string(constants.AuthTypeIDP) {
		// Even if server-side logout fails, we want the browser to redirect
		if idpToken != "" {
			_, _ = s.idpClient.Logout(ctx, cfg, idpToken, userInfo.ID)
		}

		logoutURL := fmt.Sprintf(
			"%s/auth/logout?client_id=%s",
			cfg.IDPBaseUrl,
			cfg.IDPClientID,
		)
		return logoutURL, nil
	}

	return "", nil
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
		fmt.Sprintf("%s/auth/authorize", cfg.IDPBaseUrl),
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
	ipAddress, userAgent string,
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

	// Perform JIT Provisioning & Whitelist Gate
	// User Existence Check (Anchor Lookup using Email)
	localUser, err := s.repo.GetUserByEmail(
		ctx,
		userInfo.Email,
		string(constants.AuthTypeIDP),
	)
	if err == sql.ErrNoRows {
		// JIT Provisioning (First Login Only)
		err = datastore.RunInTransaction(
			ctx,
			s.repo.(*users.Repository).GetDB(),
			func(tx datastore.DB) error {
				// Condition A: Whitelist Verification
				whitelistRoleID, err := s.repo.CheckUserWhitelist(
					ctx,
					userInfo.Email,
				)

				var assignedRoleID int
				if err == nil {
					assignedRoleID = whitelistRoleID
				} else if err == sql.ErrNoRows {
					// Default assigned role to Student
					assignedRoleID = int(constants.StudentRoleID)
				} else {
					return err
				}

				localUser = &users.User{
					ID:        uuid.NewString(),
					Email:     userInfo.Email,
					RoleID:    assignedRoleID,
					FirstName: userInfo.FirstName,
					LastName:  userInfo.LastName,
					SuffixName: sql.NullString{
						String: userInfo.SuffixName,
						Valid:  userInfo.SuffixName != "",
					},
					AuthType:     string(constants.AuthTypeIDP),
					PasswordHash: sql.NullString{Valid: false},
					IsActive:     1,
				}

				// Insert provisioned user utilizing existing repository method
				return s.repo.CreateUser(ctx, tx, *localUser)
			},
		)
		if err != nil {
			return "", "", "", "", "", fmt.Errorf(
				"[AuthService] {Provision IDP User}: %w",
				err,
			)
		}
	} else if err != nil {
		// Database failure during Anchor Lookup
		return "", "", "", "", "", fmt.Errorf(
			"[AuthService] {Anchor Check}: %w",
			err,
		)
	}

	appUserID := localUser.ID

	// Map Role ID to Name for frontend redirection
	role, err := s.repo.GetRoleByID(ctx, localUser.RoleID)
	roleName := "student"
	if err == nil && role != nil {
		roleName = role.Name
	}

	// Generate internal App Tokens using the actual app IDs
	appAccessToken, accessClaims, err := tokens.NewService().GenerateToken(
		userInfo.Email,
		appUserID,
		localUser.RoleID,
		"",
		string(constants.AuthTypeIDP),
		constants.AccessTokenMaxAge,
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
		localUser.RoleID,
		"",
		string(constants.AuthTypeIDP),
		constants.RefreshTokenMaxAge,
	)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf(
			"[AuthService] {Generate App Refresh Token}: %w",
			err,
		)
	}

	// Update Redis: App Access session
	val := map[string]string{
		"userID":          appUserID,
		"tokenType":       string(constants.AuthTypeIDP),
		"appRefreshToken": appRefreshToken,
		"idpAccessToken":  idpAccessToken,
		"ipAddress":       ipAddress,
		"userAgent":       userAgent,
	}
	if err := s.sessionService.StoreUserToken(
		ctx,
		appUserID,
		sessions.NewJTI(accessClaims.ID),
		val,
		constants.RefreshTokenMaxAge,
	); err != nil {
		return "", "", "", "", "", fmt.Errorf(
			"[AuthService] {Store Access in Redis}: %w",
			err,
		)
	}

	// Store IDP Refresh Token in Redis associated with the App
	// Refresh Token's ID (jti)
	idpRefreshKey := sessions.NewJTI(refreshClaims.ID).ToIDPRefreshKey()
	err = s.redis.Set(
		ctx,
		idpRefreshKey,
		idpRefreshToken,
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

// ValidateIDPSession checks if the provided session ID is valid on the IDP.
func (s *Service) ValidateIDPSession(
	ctx context.Context,
	sessionID string,
	cfg *config.Config,
) (*idp.IDPSessionResponse, error) {
	return s.idpClient.ValidateSession(ctx, sessionID, cfg)
}

func (s *Service) GetIDPLogoutURL(cfg *config.Config) string {
	return fmt.Sprintf(
		"%s/auth/logout?client_id=%s",
		cfg.IDPBaseUrl,
		cfg.IDPClientID,
	)
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
