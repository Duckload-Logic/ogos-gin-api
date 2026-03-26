package users

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/response"
)

type Handler struct {
	service ServiceInterface
}

// NewHandler creates a new users handler.
func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}

// ========================================
// |                                      |
// |      RETRIEVE HANDLER FUNCTIONS      |
// |                                      |
// ========================================

// GetMe godoc
// @Summary      Get current user
// @Description  Retrieves information about the currently authenticated user.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200      {object}  GetUserResponse        "Returns current user details"
// @Failure      500      {object}  map[string]string     "Failed to get current user"
// @Router       /users/me [get]
// GetMe retrieves the currently authenticated user's information.
func (h *Handler) GetMe(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	resp, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("[GetMe] {GetUserByID}: %v", err)
		response.SendError(
			c,
			"Failed to get current user",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, resp)
}

// GetUserByEmail godoc
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
// GetUserByEmail retrieves user information by their email address.
func (h *Handler) GetUserByEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		log.Printf(
			"[GetUserByEmail] {Check Query Email}: Email query parameter is required",
		)
		response.SendFail(
			c,
			gin.H{"error": "Email query parameter is required"},
		)
		return
	}

	authType := c.DefaultQuery("auth_type", "native")

	resp, err := h.service.GetUserByEmail(c.Request.Context(), email, authType)
	if err != nil {
		log.Printf("[GetUserByEmail] {GetUserByEmail}: %v", err)
		response.SendError(
			c,
			"Failed to get user by email",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, resp)
}
