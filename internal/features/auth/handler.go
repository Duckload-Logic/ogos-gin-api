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

	// Authenticate user
	token, refreshToken, err := h.service.AuthenticateUser(
		c, req.Email, req.Password,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

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

	// Refresh token
	newToken, newRefreshToken, err := h.service.RefreshToken(
		c, req.RefreshToken,
	)
	if err != nil {
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
	var req RefreshTokenDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	h.service.Logout(c, req.RefreshToken)

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
