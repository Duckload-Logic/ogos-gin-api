package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ========================================
// |                                      |
// |      RETRIEVE HANDLER FUNCTIONS      |
// |                                      |
// ========================================

// HandleGetCurrentUser godoc
// @Summary      Get current user
// @Description  Retrieves information about the currently authenticated user.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200      {object}  GetUserResponse        "Returns current user details"
// @Failure      500      {object}  map[string]string     "Failed to get current user"
// @Router       /users/me [get]
func (h *Handler) HandleGetCurrentUser(c *gin.Context) {
	email, exists := c.Get("userEmail")
	if !exists {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get current user"},
		)
		return
	}

	userEmail, ok := email.(string)
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get current user"},
		)
		return
	}

	resp, err := h.service.GetUserByEmail(c.Request.Context(), userEmail)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get current user"},
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// HandleGetUserByEmail godoc
// @Summary      Get user by email
// @Description  Retrieves user information based on the provided email.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        email   query     string true "User Email"
// @Success      200      {object}  GetUserResponse        "Returns user details"
// @Failure      400      {object}  map[string]string     "Email query parameter is required"
// @Failure      500      {object}  map[string]string     "Failed to get user by email"
// @Router       /users [get]
func (h *Handler) HandleGetUserByEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Email query parameter is required"},
		)
		return
	}

	resp, err := h.service.GetUserByEmail(c.Request.Context(), email)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get user by email"},
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}
