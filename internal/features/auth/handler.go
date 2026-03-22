package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/clients/idp"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/features/logs"
)

type Handler struct {
	service    *Service
	logService *logs.Service
	cfg        *config.Config
}

func NewHandler(s *Service, logService *logs.Service, cfg *config.Config) *Handler {
	return &Handler{service: s, logService: logService, cfg: cfg}
}

// HandleLogin godoc
// @Summary      User login
// @Description  Authenticates a user and sets JWT cookies.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body      LoginDTO true "Login Credentials"
// @Success      200     {object}  map[string]interface{} "Returns user info (optional)"
// @Failure      400     {object}  map[string]string
// @Failure      401     {object}  map[string]string
// @Router       /auth/login [post]
func (h *Handler) HandleLogin(c *gin.Context) {
	var req LoginDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[HandleLogin] {Bind JSON}: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	ip := c.ClientIP()
	ua := c.Request.UserAgent()

	userID, token, refreshToken, err := h.service.AuthenticateUser(c, req.Email, req.Password)
	if err != nil {
		h.logService.Record(c.Request.Context(), logs.LogEntry{
			Category:  logs.CategorySecurity,
			Action:    logs.ActionLoginFailed,
			Message:   fmt.Sprintf("Failed login attempt for %s: %s", req.Email, err.Error()),
			UserID:    userID,
			UserEmail: req.Email,
			IPAddress: ip,
			UserAgent: ua,
		})
		log.Printf("[HandleLogin] {AuthenticateUser}: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
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
	h.logService.Record(c.Request.Context(), logs.LogEntry{
		Category:  logs.CategorySecurity,
		Action:    logs.ActionLoginSuccess,
		Message:   fmt.Sprintf("User %s logged in successfully", req.Email),
		UserID:    userID,
		UserEmail: req.Email,
		IPAddress: ip,
		UserAgent: ua,
	})

	// Optionally return user info (but no tokens)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// HandleRefreshToken godoc
// @Summary      Refresh JWT token
// @Description  Refreshes the JWT token using the refresh token cookie.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]string "New access token (optional)"
// @Failure      401 {object} map[string]string
// @Router       /auth/refresh [post]
func (h *Handler) HandleRefreshToken(c *gin.Context) {
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
		log.Printf("[HandleRefreshToken] {Check Token}: Refresh token missing")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication session missing or expired"})
		return
	}

	ip := c.ClientIP()
	ua := c.Request.UserAgent()

	newToken, newRefreshToken, err := h.service.RefreshToken(c, refreshToken, h.cfg)
	if err != nil {
		h.logService.Record(c.Request.Context(), logs.LogEntry{
			Category:  logs.CategorySecurity,
			Action:    logs.ActionInvalidToken,
			Message:   "Token refresh failed: invalid or expired refresh token",
			IPAddress: ip,
			UserAgent: ua,
		})
		log.Printf("[HandleRefreshToken] {RefreshToken}: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Set new cookies
	if h.cfg.IsProduction {
		c.SetSameSite(http.SameSiteNoneMode)
	} else {
		c.SetSameSite(http.SameSiteLaxMode) // or omit – default is Lax
	}
	c.SetCookie("access_token", newToken, int(AccessTokenTTL), "/", "", h.cfg.IsProduction, true)
	c.SetCookie("refresh_token", newRefreshToken, int(RefreshTokenTTL), "/", "", h.cfg.IsProduction, true)

	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed"})
}

// HandleGetMe godoc
// @Summary      Get current user info
// @Description  Retrieves information about the currently authenticated user (native or IDP).
// @Tags         Auth
// @Produce      json
// @Success      200 {object} MeResponse
// @Failure      401 {object} map[string]string
// @Router       /auth/me [get]
func (h *Handler) HandleGetMe(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	accessToken, _ := c.Cookie("access_token")
	if accessToken == "" {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			accessToken = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if accessToken == "" {
		log.Printf("[HandleGetMe] {Check Token}: Access token missing")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token missing"})
		return
	}

	tokenType, _ := c.Get("tokenType")
	tType := "native"
	if tt, ok := tokenType.(string); ok {
		tType = tt
	}

	if tType == "native" {
		resp, err := h.service.GetMe(c.Request.Context(), userID, tType)
		if err != nil {
			log.Printf("[HandleGetMe] {GetMe}: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
			return
		}

		c.JSON(http.StatusOK, resp)
	} else {
		// Retrieve IDP Access Token from context (set by AuthMiddleware from Redis)
		idpAccessToken, ok := c.Get("idpAccessToken")
		if !ok || idpAccessToken == "" {
			log.Printf("[HandleGetMe] {Check IDP Token}: IDP access token missing from session")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "IDP session expired"})
			return
		}

		resp, err := h.service.idpClient.GetUserInfo(c.Request.Context(), idpAccessToken.(string), h.cfg)
		if err != nil {
			log.Printf("[HandleGetMe] {Get IDP UserInfo}: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to get user info from IDP"})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

// HandleLogout godoc
// @Summary      User logout
// @Description  Invalidates the user's tokens by clearing cookies.
// @Tags         Auth
// @Success      200 {object} map[string]string
// @Router       /auth/logout [post]
func (h *Handler) HandleLogout(c *gin.Context) {
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
		_ = h.service.Logout(c.Request.Context(), tokenString, tokenType.(string))
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
		h.logService.Record(c.Request.Context(), logs.LogEntry{
			Category:  logs.CategorySecurity,
			Action:    logs.ActionLogout,
			Message:   fmt.Sprintf("User %s logged out", userEmail),
			UserID:    userID.(string),
			UserEmail: userEmail.(string),
			IPAddress: c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// IDP integration handlers

// GetAuthorizeURL godoc
// @Summary      Get IDP authorization URL
// @Description  Generates OAuth 2.0 authorization URL with PKCE
// @Tags         Auth
// @Produce      json
// @Success      200 {object} idp.IDPAuthorizeURLResponse
// @Failure      500 {object} map[string]string
// @Router       /auth/idp/authorize-url [get]
func (h *Handler) GetAuthorizeURL(c *gin.Context) {
	// Generate authorization URL with state and PKCE parameters
	authURL, err := h.service.GetAuthorizeURL(h.cfg)
	if err != nil {
		h.logService.Record(c.Request.Context(), logs.LogEntry{
			Category:  logs.CategorySecurity,
			Action:    logs.ActionLoginFailed,
			Message:   fmt.Sprintf("[GetAuthorizeURL] {Generate URL}: %s", err.Error()),
			IPAddress: c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authorization URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"authorization_url": authURL})
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
		h.logService.Record(c.Request.Context(), logs.LogEntry{
			Category:  logs.CategorySecurity,
			Action:    logs.ActionLoginFailed,
			Message:   fmt.Sprintf("[PostIDPToken] {Bind JSON}: %s", err.Error()),
			IPAddress: c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Perform token exchange
	accessToken, refreshToken, userID, userEmail, roleName, err := h.service.PostIDPTokenExchange(c.Request.Context(), req.Code, h.cfg)
	if err != nil {
		h.logService.Record(c.Request.Context(), logs.LogEntry{
			Category:  logs.CategorySecurity,
			Action:    logs.ActionLoginFailed,
			Message:   fmt.Sprintf("[PostIDPToken] {Service Call}: %s", err.Error()),
			IPAddress: c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
		})
		log.Printf("[PostIDPToken] {Service Call}: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
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
	h.logService.Record(c.Request.Context(), logs.LogEntry{
		Category:  logs.CategorySecurity,
		Action:    logs.ActionLoginSuccess,
		Message:   fmt.Sprintf("User %s logged in successfully via IDP", userEmail),
		UserID:    userID,
		UserEmail: userEmail,
		IPAddress: c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	})

	// Return Message and Role info for immediate redirect
	c.JSON(http.StatusOK, gin.H{
		"userID":    userID,
		"userEmail": userEmail,
		"role":      roleName,
		"message":   "Login successful",
		"roles":     []string{roleName},
	})
}

// HandleValidateIDPSession godoc
// @Summary      Validate IDP session
// @Description  Validates the current IDP session using the idp_session cookie.
// @Tags         Auth
// @Produce      json
// @Success      200 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /auth/idp/session [get]
func (h *Handler) HandleValidateIDPSession(c *gin.Context) {
	sessionID, err := c.Cookie("idp_session")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "IDP session cookie missing"})
		return
	}

	resp, err := h.service.ValidateIDPSession(c, sessionID, h.cfg)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid IDP session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": resp.Message})
}

// containsStr checks if a string contains a substring
func containsStr(s, substr string) bool {

	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
