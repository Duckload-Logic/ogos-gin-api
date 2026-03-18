package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
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
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}
		refreshToken = req.RefreshToken
	}

	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token missing"})
		return
	}

	ip := c.ClientIP()
	ua := c.Request.UserAgent()

	newToken, newRefreshToken, err := h.service.RefreshToken(c, refreshToken)
	if err != nil {
		h.logService.Record(c.Request.Context(), logs.LogEntry{
			Category:  logs.CategorySecurity,
			Action:    logs.ActionInvalidToken,
			Message:   "Token refresh failed: invalid or expired refresh token",
			IPAddress: ip,
			UserAgent: ua,
		})
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

// HandleLogout godoc
// @Summary      User logout
// @Description  Invalidates the user's tokens by clearing cookies.
// @Tags         Auth
// @Success      200 {object} map[string]string
// @Router       /auth/logout [post]
func (h *Handler) HandleLogout(c *gin.Context) {
	// Clear cookies
	if h.cfg.IsProduction {
		c.SetSameSite(http.SameSiteNoneMode)
	} else {
		c.SetSameSite(http.SameSiteLaxMode) // or omit – default is Lax
	}

	c.SetCookie("access_token", "", -1, "/", "", h.cfg.IsProduction, true)
	c.SetCookie("refresh_token", "", -1, "/", "", h.cfg.IsProduction, true)

	// Log event (if user info available)
	userID, _ := c.Get("userID")
	userEmail, _ := c.Get("userEmail")
	if userID != nil {
		h.logService.Record(c.Request.Context(), logs.LogEntry{
			Category:  logs.CategorySecurity,
			Action:    logs.ActionLogout,
			Message:   fmt.Sprintf("User %s logged out", userEmail),
			UserID:    userID.(int),
			UserEmail: userEmail.(string),
			IPAddress: c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func (h *Handler) GetAuthRedirect(c *gin.Context) {
	// HARDCODED BASE URL FOR TESTING
	const testBaseURL = "https://identity-provider.isaxbsit2027.com"

	// Build the full URL using the hardcoded base
	authURL := fmt.Sprintf(
		"%s/auth/authorize?client_id=%s",
		testBaseURL,
		h.cfg.IDPClientID,
	)

	log.Printf("[TEST] Redirecting to: %s", authURL)

	// Using StatusSeeOther (303) is sometimes more reliable for browsers
	c.Redirect(http.StatusSeeOther, authURL)
}

func (h *Handler) GetAuthCallback(c *gin.Context) {
	//Capture code
	code := c.Query("code")
	if code == "" {
		log.Printf("[GetAuthCallback] {Query Parameter}: missing code")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Auth code missing"})
		return
	}

	//Token Exchange
	token, err := h.exchangeCode(code)
	if err != nil {
		log.Printf("[GetAuthCallback] {Token Exchange}: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token fail"})
		return
	}

	//Fetch Identity
	user, err := h.fetchIDPUser(token)
	if err != nil {
		log.Printf("[GetAuthCallback] {API Request}: fetch /me failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ID fail"})
		return
	}

	//Finalize Session & Role Check
	h.finalizeIDPLogin(c, user)
}

// exchangeCode
func (h *Handler) exchangeCode(code string) (string, error) {
	payload := map[string]string{
		"client_id":     h.cfg.IDPClientID,
		"client_secret": h.cfg.IDPClientSecret,
		"code":          code,
		"redirect_uri":  h.cfg.IDPRedirectURI,
		"grant_type":    "authorization_code",
	}
	jsonData, _ := json.Marshal(payload)

	// DIRECT URL: Bypassing the config bug
	targetURL := "https://identity-provider.isaxbsit2027.com/api/v1/auth/token"

	resp, err := http.Post(
		targetURL,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("exchange failed with status: %d", resp.StatusCode)
	}

	var res struct {
		AccessToken string `json:"access_token"`
	}
	json.NewDecoder(resp.Body).Decode(&res)
	return res.AccessToken, nil
}

func (h *Handler) fetchIDPUser(token string) (*IDPUser, error) {
	targetURL := "https://identity-provider.isaxbsit2027.com/api/v1/auth/me"

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IDP /me failed: %d", resp.StatusCode)
	}

	var user IDPUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (h *Handler) finalizeIDPLogin(c *gin.Context, idpUser *IDPUser) {
	localUser, err := h.service.SyncIDPUser(c.Request.Context(), idpUser)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: " + err.Error()})
		return
	}

	token, refresh, err := h.service.GenerateTokens(localUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create local session"})
		return
	}

	h.setAuthCookies(c, token, refresh)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func (h *Handler) setAuthCookies(c *gin.Context, token, refresh string) {
	if h.cfg.IsProduction {
		c.SetSameSite(http.SameSiteNoneMode)
	} else {
		c.SetSameSite(http.SameSiteLaxMode)
	}
	c.SetCookie("access_token", token, int(AccessTokenTTL), "/", "", h.cfg.IsProduction, true)
	c.SetCookie("refresh_token", refresh, int(RefreshTokenTTL), "/", "", h.cfg.IsProduction, true)
}
