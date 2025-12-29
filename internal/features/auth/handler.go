package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) HandleLogin(c *gin.Context) {
	// Map request body to struct
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Authenticate user
	token, refreshToken, err := h.service.AuthenticateUser(
		c, req.Email, req.Password,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Return tokens
	c.JSON(http.StatusOK, TokenResponse{
		Token:        token,
		RefreshToken: refreshToken,
	})
}

func (h *Handler) HandleRefreshToken(c *gin.Context) {
	// Map request body to struct
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Refresh token
	newToken, newRefreshToken, err := h.service.RefreshToken(
		c, req.RefreshToken,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Return new tokens
	c.JSON(http.StatusOK, TokenResponse{
		Token:        newToken,
		RefreshToken: newRefreshToken,
	})
}
