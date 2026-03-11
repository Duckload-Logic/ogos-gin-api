package auth

import (
	"fmt"
	"net/http"

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
	c.SetCookie("access_token", token, int(AccessTokenTTL), "/", "", h.cfg.IsProduction, true) // 1 hour
	// Refresh token: longer-lived, HTTP-only
	c.SetCookie("refresh_token", refreshToken, int(RefreshTokenTTL), "/", "", h.cfg.IsProduction, true) // 12 hours

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
