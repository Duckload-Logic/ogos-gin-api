package users

import (
	"fmt"
	"net/http"
	"strconv"

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
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get current user"},
		)
		return
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get current user"},
		)
		return
	}

	resp, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		fmt.Println("Error getting current user:", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get current user"},
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// HandleGetUserByID godoc
// @Summary      Get user by ID
// @Description  Retrieves user information based on the provided user ID.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        userID   path      int true "User ID"
// @Success      200      {object}  GetUserResponse        "Returns user details"
// @Failure      400      {object}  map[string]string     "Invalid user ID"
// @Failure      500      {object}  map[string]string     "Failed to get user by ID"
// @Router       /users/{userID} [get]
func (h *Handler) HandleGetUserByID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	resp, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		fmt.Println("Error getting user by ID:", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get user by ID"},
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// HandleGetUserByEmail
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
		fmt.Println("Error getting user by email:", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to get user by email"},
		)
		return
	}

	c.JSON(http.StatusOK, resp)
}
