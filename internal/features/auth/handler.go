package auth

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/response"
	"github.com/olazo-johnalbert/duckload-api/internal/core/sessions"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/core/tokens"
	"github.com/olazo-johnalbert/duckload-api/internal/features/logs"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/identity/idp"
)

type Handler struct {
	service    ServiceInterface
	logService logs.ServiceInterface
	cfg        *config.Config
}

func NewHandler(
	s ServiceInterface,
	logService logs.ServiceInterface,
	cfg *config.Config,
) *Handler {
	return &Handler{service: s, logService: logService, cfg: cfg}
}

// PostLogin godoc
// @Summary      User login
// @Description  Authenticates a user and sets JWT cookies.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body      LoginDTO true "Login Credentials"
// @Success      200     {object}  map[string]interface{} "Returns user info
// (optional)"
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Router       /auth/login [post]
func (h *Handler) PostLogin(c *gin.Context) {
	var req LoginDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendFail(c, gin.H{"error": "Invalid request format"})
		return
	}

	ip := c.ClientIP()
	ua := c.Request.UserAgent()

	userID, token, refreshToken, err := h.service.AuthenticateUser(
		c,
		req.Email,
		req.Password,
		ip,
		ua,
	)
	if err != nil {
		// Clear cookies on failure to prevent ghost sessions
		h.clearAuthCookies(c)

		h.logService.Record(
			c.Request.Context(),
			h.logService.GetDB(),
			audit.LogEntry{
				Level:    audit.LevelError,
				Category: audit.CategorySecurity,
				Action:   audit.ActionLoginFailed,
				Message: fmt.Sprintf(
					"Failed login attempt for %s: %s",
					req.Email,
					err.Error(),
				),
				UserID:    structs.StringToNullableString(userID),
				UserEmail: structs.StringToNullableString(req.Email),
				IPAddress: structs.StringToNullableString(ip),
				UserAgent: structs.StringToNullableString(ua),
			},
		)
		log.Printf("[PostLogin] {AuthenticateUser}: %v", err)
		response.SendFail(
			c,
			gin.H{"error": err.Error()},
			http.StatusUnauthorized,
		)
		return
	}

	// Set cookies
	h.setAuthCookies(c, token, refreshToken)

	// Log success
	h.logService.Record(
		c.Request.Context(),
		h.logService.GetDB(),
		audit.LogEntry{
			Level:     audit.LevelInfo,
			Category:  audit.CategorySecurity,
			Action:    audit.ActionLoginSuccess,
			Message:   fmt.Sprintf("User %s logged in successfully", req.Email),
			UserID:    structs.StringToNullableString(userID),
			UserEmail: structs.StringToNullableString(req.Email),
			IPAddress: structs.StringToNullableString(ip),
			UserAgent: structs.StringToNullableString(ua),
		},
	)

	// Optionally return user info (but no tokens)
	response.SendSuccess(c, gin.H{"message": "Login successful"})
}

// PostRegister godoc
// @Summary      User registration
// @Description  Creates a new developer account.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body      RegisterDTO true "Registration Data"
// @Success      201     {object}  map[string]string "Success message"
// @Failure      400     {object}  map[string]string
// @Failure      409     {object}  map[string]string
// @Router       /auth/register [post]
func (h *Handler) PostRegister(c *gin.Context) {
	var req RegisterDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendFail(c, gin.H{"error": "Invalid request format"})
		return
	}

	registrationID, err := h.service.RegisterUser(c.Request.Context(), req)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "user already exists" {
			status = http.StatusConflict
		}
		response.SendFail(c, gin.H{"error": err.Error()}, status)
		return
	}

	response.SendSuccess(
		c,
		gin.H{"registrationId": registrationID},
		http.StatusCreated,
	)
}

func (h *Handler) PostResendVerification(c *gin.Context) {
	registrationID := c.Query("registration_id")
	if registrationID == "" {
		response.SendFail(c, gin.H{"error": "Registration ID is required"})
		return
	}

	err := h.service.ResendVerification(c.Request.Context(), registrationID)
	if err != nil {
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	response.SendSuccess(
		c,
		gin.H{"message": "Verification email sent successfully"},
	)
}

func (h *Handler) PostVerify(c *gin.Context) {
	registrationID := c.Query("registration_id")
	if registrationID == "" {
		response.SendFail(c, gin.H{"error": "Registration ID is required"})
		return
	}

	var req VerifyDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendFail(c, gin.H{"error": "Invalid request format"})
		return
	}

	userID, userEmail, err := h.service.VerifyUser(
		c.Request.Context(),
		registrationID,
		req.VerificationOTP,
	)
	if err != nil {
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	// Log success
	h.logService.Record(
		c.Request.Context(),
		h.logService.GetDB(),
		audit.LogEntry{
			Level:    audit.LevelInfo,
			Category: audit.CategorySecurity,
			Action:   audit.ActionUserCreated,
			Message: fmt.Sprintf(
				"Developer %s registered successfully",
				userEmail,
			),
			UserID:    structs.StringToNullableString(userID),
			UserEmail: structs.StringToNullableString(userEmail),
			IPAddress: structs.StringToNullableString(c.ClientIP()),
			UserAgent: structs.StringToNullableString(c.Request.UserAgent()),
		},
	)

	response.SendSuccess(c, gin.H{"message": "User verified successfully"})
}

// PostRefreshToken godoc
// @Summary      Refresh JWT token
// @Description  Refreshes the JWT token using the refresh token cookie.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]string "New access token (optional)"
// @Failure      401 {object} map[string]string
// @Router       /auth/refresh [post]
func (h *Handler) PostRefreshToken(c *gin.Context) {
	// Get Access Token (which contains the session Handle/JTI)
	accessToken, err := c.Cookie(constants.AccessTokenCookieName)
	if err != nil {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			accessToken = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if accessToken == "" {
		response.SendFail(
			c,
			gin.H{"error": "Access token missing"},
			http.StatusUnauthorized,
		)
		return
	}

	// Parse token UNVERIFIED to get the JTI (it might be expired)
	tokenService := tokens.NewService()
	claims, err := tokenService.ParseTokenUnverified(accessToken)
	if err != nil {
		response.SendFail(
			c,
			gin.H{"error": "Invalid session handle"},
			http.StatusUnauthorized,
		)
		return
	}

	ip := c.ClientIP()
	ua := c.Request.UserAgent()

	// Refresh using the JTI
	newAccessToken, newRefreshToken, err := h.service.RefreshToken(
		c,
		sessions.NewJTI(claims.ID),
		h.cfg,
		ip,
		ua,
	)
	if err != nil {
		// Remove previous cookies
		h.clearAuthCookies(c)

		h.logService.Record(
			c.Request.Context(),
			h.logService.GetDB(),
			audit.LogEntry{
				Level:     audit.LevelError,
				Category:  audit.CategorySecurity,
				Action:    audit.ActionInvalidToken,
				UserID:    structs.StringToNullableString(claims.UserID),
				UserEmail: structs.StringToNullableString(claims.UserEmail),
				Message:   "Token refresh failed: " + err.Error(),
				IPAddress: structs.StringToNullableString(ip),
				UserAgent: structs.StringToNullableString(ua),
			},
		)
		log.Printf("[PostRefreshToken] {RefreshToken}: %v", err)
		response.SendFail(
			c,
			gin.H{"error": "Session expired or invalid"},
			http.StatusUnauthorized,
		)
		return
	}

	h.setAuthCookies(c, newAccessToken, newRefreshToken)

	response.SendSuccess(c, gin.H{"message": "Session refreshed"})
}

// GetMe godoc
// @Summary      Get current user info
// @Description  Retrieves information about the currently authenticated user
// (native or IDP).
// @Tags         Auth
// @Produce      json
// @Success      200 {object} MeResponse
// @Failure      401 {object} map[string]string
// @Router       /auth/me [get]
func (h *Handler) GetMe(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	accessToken, _ := c.Cookie("access_token")
	if accessToken == "" {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			accessToken = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if accessToken == "" {
		log.Printf("[GetMe] {Check Token}: Access token missing")
		response.SendFail(
			c,
			gin.H{"error": "Access token missing"},
			http.StatusUnauthorized,
		)
		return
	}

	tokenType, _ := c.Get("tokenType")
	tType := "native"
	if tt, ok := tokenType.(string); ok {
		tType = tt
	}

	resp, err := h.service.GetMe(c.Request.Context(), userID, tType)
	if err != nil {
		log.Printf("[GetMe] {GetMe}: %v", err)
		response.SendError(
			c,
			"Failed to get user info",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, resp)
}

// GetLogout godoc
// @Summary      User logout
// @Description  Invalidates the user's tokens by clearing cookies.
// @Tags         Auth
// @Success      200 {object} map[string]string
// @Router       /auth/logout [get]
func (h *Handler) GetLogout(c *gin.Context) {
	// Extract token to invalidate in Redis
	var tokenString string
	cookie, err := c.Cookie("access_token")
	if err == nil && cookie != "" {
		tokenString = cookie
	}
	if tokenString == "" {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	tokenType, _ := c.Get("tokenType")
	tType := "native"
	if tt, ok := tokenType.(string); ok {
		tType = tt
	}

	var logoutUrl string
	if tokenString != "" {
		logoutUrl, _ = h.service.Logout(
			c.Request.Context(),
			tokenString,
			tType,
			h.cfg,
		)
	}

	// Always clear cookies before any redirect or success response
	h.clearAuthCookies(c)

	// Log event
	userID, _ := c.Get("userID")
	userEmail, _ := c.Get("userEmail")
	if userID != nil && userEmail != nil {
		h.logService.Record(
			c.Request.Context(),
			h.logService.GetDB(),
			audit.LogEntry{
				Level:    audit.LevelInfo,
				Category: audit.CategorySecurity,
				Action:   audit.ActionLogout,
				Message:  fmt.Sprintf("User %s logged out", userEmail),

				UserID:    structs.StringToNullableString(userID.(string)),
				UserEmail: structs.StringToNullableString(userEmail.(string)),
				IPAddress: structs.StringToNullableString(c.ClientIP()),
				UserAgent: structs.StringToNullableString(
					c.Request.UserAgent(),
				),
			},
		)
	}

	// Determine redirection target
	redirectTarget := "/"
	if logoutUrl != "" {
		redirectTarget = logoutUrl
	}

	// Handle dynamic redirection for native/local sessions
	if tType != string(constants.AuthTypeIDP) {
		candidate := c.Query("redirect_uri")

		if candidate != "" {
			parsedURL, err := url.Parse(candidate)
			if err == nil {
				origin := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
				if h.isAllowedOrigin(origin) {
					redirectTarget = origin + "/"
				}
			}
		}
	}

	c.Redirect(http.StatusFound, redirectTarget)
}

// isAllowedOrigin checks if the given origin is permitted for redirects.
func (h *Handler) isAllowedOrigin(origin string) bool {
	if h.cfg.IsProduction {
		// Support subdomains of dllbsit2027.com
		return strings.HasSuffix(origin, ".dllbsit2027.com")
	}

	// Support localhost development
	return strings.HasPrefix(origin, "http://localhost")
}

// IDP integration handlers

// GetAuthorizeURL godoc
// @Summary      Get IDP authorization URL
// @Description  Redirects to OAuth 2.0 authorization page on the IDP.
// @Tags         Auth
// @Produce      json
// @Success      302 {string} string "Redirect to IDP login page"
// @Failure      500 {object} map[string]string
// @Router       /auth/idp/authorize [get]
func (h *Handler) GetAuthorizeURL(c *gin.Context) {
	// Generate authorization URL with state and PKCE parameters
	authURL, err := h.service.GetAuthorizeURL(h.cfg)
	if err != nil {
		h.logService.Record(
			c.Request.Context(),
			h.logService.GetDB(),
			audit.LogEntry{
				Level:    audit.LevelError,
				Category: audit.CategorySecurity,
				Action:   audit.ActionLoginFailed,
				Message: fmt.Sprintf(
					"[GetAuthorizeURL] {Generate URL}: %s",
					err.Error(),
				),
				IPAddress: structs.StringToNullableString(c.ClientIP()),
				UserAgent: structs.StringToNullableString(
					c.Request.UserAgent(),
				),
			},
		)
		response.SendError(
			c,
			"Failed to generate authorization URL",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	c.Redirect(http.StatusFound, authURL)
}

// PostIDPToken godoc
// @Summary      Exchange IDP authorization code for tokens
// @Description  Completes OAuth 2.0 flow and provisions user
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body idp.IDPTokenExchangeRequest true "Code & State"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /auth/idp/token [post]
func (h *Handler) PostIDPToken(c *gin.Context) {
	var req idp.IDPTokenExchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logService.Record(
			c.Request.Context(),
			h.logService.GetDB(),
			audit.LogEntry{
				Category: audit.CategorySecurity,
				Action:   audit.ActionLoginFailed,
				Message: fmt.Sprintf(
					"[PostIDPToken] {Bind JSON}: %s",
					err.Error(),
				),
				IPAddress: structs.StringToNullableString(c.ClientIP()),
				UserAgent: structs.StringToNullableString(
					c.Request.UserAgent(),
				),
			},
		)
		response.SendFail(c, gin.H{"error": "Invalid request body"})
		return
	}

	// Perform token exchange
	accessToken, refreshToken, userID, userEmail, roleName,
		err := h.service.PostIDPTokenExchange(
		c.Request.Context(),
		req.Code,
		h.cfg,
		c.ClientIP(),
		c.Request.UserAgent(),
	)
	if err != nil {
		// Clear cookies on failure to prevent ghost sessions
		h.clearAuthCookies(c)

		h.logService.Record(
			c.Request.Context(),
			h.logService.GetDB(),
			audit.LogEntry{
				Category: audit.CategorySecurity,
				Action:   audit.ActionLoginFailed,
				Message: fmt.Sprintf(
					"[PostIDPTokenExchange] {Service Call}: %s",
					err.Error(),
				),
				IPAddress: structs.StringToNullableString(c.ClientIP()),
				UserAgent: structs.StringToNullableString(
					c.Request.UserAgent(),
				),
			},
		)
		log.Printf("[PostIDPTokenExchange] {Service Call}: %v", err)
		response.SendFail(
			c,
			gin.H{"error": err.Error()},
			http.StatusUnauthorized,
		)
		return
	}

	// Set cookies for frontend
	h.setAuthCookies(c, accessToken, refreshToken)

	// Log success
	h.logService.Record(
		c.Request.Context(),
		h.logService.GetDB(),
		audit.LogEntry{
			Level:    audit.LevelInfo,
			Category: audit.CategorySecurity,
			Action:   audit.ActionLoginSuccess,
			Message: fmt.Sprintf(
				"User %s logged in successfully via IDP",
				userEmail,
			),
			UserID:    structs.StringToNullableString(userID),
			UserEmail: structs.StringToNullableString(userEmail),
			IPAddress: structs.StringToNullableString(c.ClientIP()),
			UserAgent: structs.StringToNullableString(c.Request.UserAgent()),
		},
	)

	// Return Message and Role info for immediate redirect
	response.SendSuccess(c, gin.H{
		"userID":    userID,
		"userEmail": userEmail,
		"role":      roleName,
		"message":   "Login successful",
	})
}

func (h *Handler) setAuthCookies(c *gin.Context, accessToken, refreshToken string) {
	if h.cfg.IsProduction {
		c.SetSameSite(http.SameSiteNoneMode)
	} else {
		c.SetSameSite(http.SameSiteLaxMode)
	}

	if accessToken != "" {
		c.SetCookie(
			constants.AccessTokenCookieName,
			accessToken,
			int(constants.RefreshTokenMaxAge),
			constants.CookiePathRoot,
			"",
			h.cfg.IsProduction,
			true,
		)
	}

	if refreshToken != "" {
		c.SetCookie(
			constants.RefreshTokenCookieName,
			refreshToken,
			int(constants.RefreshTokenMaxAge),
			constants.CookiePathRoot,
			"",
			h.cfg.IsProduction,
			true,
		)
	}
}

func (h *Handler) clearAuthCookies(c *gin.Context) {
	if h.cfg.IsProduction {
		c.SetSameSite(http.SameSiteNoneMode)
	} else {
		c.SetSameSite(http.SameSiteLaxMode)
	}

	c.SetCookie(
		constants.AccessTokenCookieName,
		"",
		-1,
		constants.CookiePathRoot,
		"",
		h.cfg.IsProduction,
		true,
	)
	c.SetCookie(
		constants.RefreshTokenCookieName,
		"",
		-1,
		constants.CookiePathRoot,
		"",
		h.cfg.IsProduction,
		true,
	)
}
