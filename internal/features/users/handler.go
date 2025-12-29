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

// HandleGetUserByID
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
