package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/features/logs"
)

type Handler struct {
	service    *Service
	logService *logs.Service
}

func NewHandler(s *Service, logService *logs.Service) *Handler {
	return &Handler{service: s, logService: logService}
}

// HandleLogin godoc
// @Summary      User login
// @Description  Authenticates a user and returns JWT tokens.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body      LoginRequest true "Login Credentials"
// @Success      200     {object}  TokenResponse          "Returns {token, refreshToken}"
// @Failure      400     {object}  map[string]string      "Invalid request format"
// @Failure      401     {object}  map[string]string      "Unauthorized"
// @Router       /auth/login [post]
func (h *Handler) HandleLogin(c *gin.Context) {
	// Map request body to struct
	var req LoginDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	ip := c.ClientIP()
	ua := c.Request.UserAgent()

	// Authenticate user
	token, refreshToken, err := h.service.AuthenticateUser(
		c, req.Email, req.Password,
	)
	if err != nil {
		// Log failed login attempt
		h.logService.Record(c.Request.Context(), logs.LogEntry{
			Category:  logs.CategorySecurity,
			Action:    logs.ActionLoginFailed,
			Message:   fmt.Sprintf("Failed login attempt for %s: %s", req.Email, err.Error()),
			UserEmail: req.Email,
			IPAddress: ip,
			UserAgent: ua,
		})
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Log successful login
	h.logService.Record(c.Request.Context(), logs.LogEntry{
		Category:  logs.CategorySecurity,
		Action:    logs.ActionLoginSuccess,
		Message:   fmt.Sprintf("User %s logged in successfully", req.Email),
		UserEmail: req.Email,
		IPAddress: ip,
		UserAgent: ua,
	})

	// Return tokens
	c.JSON(http.StatusOK, TokenDTO{
		Token:        token,
		RefreshToken: refreshToken,
	})
}

// HandleRefreshToken godoc
// @Summary      Refresh JWT token
// @Description  Refreshes the JWT token using a valid refresh token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body      RefreshTokenRequest true "Refresh Token"
// @Success      200     {object}  TokenResponse          "Returns {token, refreshToken}"
// @Failure      400     {object}  map[string]string      "Invalid request format"
// @Failure      401     {object}  map[string]string      "Unauthorized"
// @Router       /auth/refresh-token [post]
func (h *Handler) HandleRefreshToken(c *gin.Context) {
	// Map request body to struct
	var req RefreshTokenDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	ip := c.ClientIP()
	ua := c.Request.UserAgent()

	// Refresh token
	newToken, newRefreshToken, err := h.service.RefreshToken(
		c, req.RefreshToken,
	)
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

	// Return new tokens
	c.JSON(http.StatusOK, TokenDTO{
		Token:        newToken,
		RefreshToken: newRefreshToken,
	})
}

// HandleLogout godoc
// @Summary      User logout
// @Description  Invalidates the user's JWT token (if token blacklisting is implemented).
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200     {object}  map[string]string      "Logout successful"
// @Failure      401     {object}  map[string]string      "Unauthorized"
// @Router       /auth/logout [post]
func (h *Handler) HandleLogout(c *gin.Context) {
	// var req RefreshTokenDTO
	// if err := c.ShouldBindJSON(&req); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
	// 	return
	// }

	// h.service.Logout(c, req.RefreshToken)

	// Log logout event
	userEmail, exists := c.Get("userEmail")
	if exists {
		h.logService.Record(c.Request.Context(), logs.LogEntry{
			Category:  logs.CategorySecurity,
			Action:    logs.ActionLogout,
			Message:   fmt.Sprintf("User %s logged out", userEmail.(string)),
			UserEmail: userEmail.(string),
			IPAddress: c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
