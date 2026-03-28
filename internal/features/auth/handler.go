package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/response"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
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
	)
	if err != nil {
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
	if h.cfg.IsProduction {
		c.SetSameSite(http.SameSiteNoneMode)
	} else {
		c.SetSameSite(http.SameSiteLaxMode) // or omit – default is Lax
	}
	// Access token: short-lived, HTTP-only, Secure in production
	c.SetCookie(
		"access_token", token, int(AccessTokenTTL), "/",
		"", h.cfg.IsProduction, true) // 1 hour
	// Refresh token: longer-lived, HTTP-only
	c.SetCookie(
		"refresh_token", refreshToken, int(RefreshTokenTTL), "/",
		"", h.cfg.IsProduction, true) // 12 hours

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
	// Try to get refresh token from cookie first
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		// Fallback: read from request body
		var req RefreshTokenDTO
		if err := c.ShouldBindJSON(&req); err == nil {
			refreshToken = req.RefreshToken
		}
	}

	if refreshToken == "" {
		log.Printf("[PostRefreshToken] {Check Token}: Refresh token missing")
		response.SendFail(
			c,
			gin.H{"error": "Authentication session missing or expired"},
			http.StatusUnauthorized,
		)
		return
	}

	ip := c.ClientIP()
	ua := c.Request.UserAgent()

	newToken, newRefreshToken, err := h.service.RefreshToken(
		c,
		refreshToken,
		h.cfg,
	)
	if err != nil {
		h.logService.Record(
			c.Request.Context(),
			h.logService.GetDB(),
			audit.LogEntry{
				Level:     audit.LevelError,
				Category:  audit.CategorySecurity,
				Action:    audit.ActionInvalidToken,
				Message:   "Token refresh failed: invalid or expired refresh token",
				IPAddress: structs.StringToNullableString(ip),
				UserAgent: structs.StringToNullableString(ua),
			},
		)
		log.Printf("[PostRefreshToken] {RefreshToken}: %v", err)
		response.SendFail(
			c,
			gin.H{"error": err.Error()},
			http.StatusUnauthorized,
		)
		return
	}

	// Set new cookies
	if h.cfg.IsProduction {
		c.SetSameSite(http.SameSiteNoneMode)
	} else {
		c.SetSameSite(http.SameSiteLaxMode) // or omit – default is Lax
	}
	c.SetCookie(
		"access_token",
		newToken,
		int(AccessTokenTTL),
		"/",
		"",
		h.cfg.IsProduction,
		true,
	)
	c.SetCookie(
		"refresh_token",
		newRefreshToken,
		int(RefreshTokenTTL),
		"/",
		"",
		h.cfg.IsProduction,
		true,
	)

	response.SendSuccess(c, gin.H{"message": "Token refreshed"})
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

// PostLogout godoc
// @Summary      User logout
// @Description  Invalidates the user's tokens by clearing cookies.
// @Tags         Auth
// @Success      200 {object} map[string]string
// @Router       /auth/logout [post]
func (h *Handler) PostLogout(c *gin.Context) {
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

	if tokenString != "" {
		_ = h.service.Logout(
			c.Request.Context(),
			tokenString,
			tokenType.(string),
			h.cfg,
		)
	}

	// Clear cookies
	if h.cfg.IsProduction {
		c.SetSameSite(http.SameSiteNoneMode)
	} else {
		c.SetSameSite(http.SameSiteLaxMode)
	}

	c.SetCookie("access_token", "", -1, "/", "", h.cfg.IsProduction, true)
	c.SetCookie("refresh_token", "", -1, "/", "", h.cfg.IsProduction, true)

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

	response.SendSuccess(c, gin.H{"message": "Logout successful"})
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
// @Param        request body IDPTokenExchangeRequest true "Code & State"
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
	)
	if err != nil {
		h.logService.Record(
			c.Request.Context(),
			h.logService.GetDB(),
			audit.LogEntry{
				Category: audit.CategorySecurity,
				Action:   audit.ActionLoginFailed,
				Message: fmt.Sprintf(
					"[PostIDPToken] {Service Call}: %s",
					err.Error(),
				),
				IPAddress: structs.StringToNullableString(c.ClientIP()),
				UserAgent: structs.StringToNullableString(
					c.Request.UserAgent(),
				),
			},
		)
		log.Printf("[PostIDPToken] {Service Call}: %v", err)
		response.SendFail(
			c,
			gin.H{"error": err.Error()},
			http.StatusUnauthorized,
		)
		return
	}

	// Set cookies for frontend
	if h.cfg.IsProduction {
		c.SetSameSite(http.SameSiteNoneMode)
	} else {
		c.SetSameSite(http.SameSiteLaxMode)
	}

	// Set access token cookie
	c.SetCookie(
		constants.AccessTokenCookieName,
		accessToken,
		constants.AccessTokenMaxAge,
		constants.CookiePathRoot,
		"",
		h.cfg.IsProduction,
		true,
	)

	// Set refresh token cookie
	c.SetCookie(
		constants.RefreshTokenCookieName,
		refreshToken,
		constants.RefreshTokenMaxAge,
		constants.CookiePathRoot,
		"",
		h.cfg.IsProduction,
		true,
	)

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
		"roles":     []string{roleName},
	})
}

// GetIDPValidateSession godoc
// @Summary      Validate IDP session
// @Description  Validates the current IDP session using the idp_session cookie.
// @Tags         Auth
// @Produce      json
// @Success      200 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /auth/idp/session [get]
func (h *Handler) GetIDPValidateSession(c *gin.Context) {
	sessionID, err := c.Cookie("idp_session")
	if err != nil {
		response.SendFail(
			c,
			gin.H{"error": "IDP session cookie missing"},
			http.StatusUnauthorized,
		)
		return
	}

	resp, err := h.service.ValidateIDPSession(c, sessionID, h.cfg)
	if err != nil {
		response.SendFail(
			c,
			gin.H{"error": "Invalid IDP session"},
			http.StatusUnauthorized,
		)
		return
	}

	response.SendSuccess(c, gin.H{"message": resp.Message})
}
