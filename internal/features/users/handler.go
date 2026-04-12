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

// GetUsers godoc
// @Summary      List all users
// @Description  Retrieves a paginated list of all users with filtering options.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "Page number"
// @Param        page_size  query     int     false  "Items per page"
// @Param        role_id    query     int     false  "Filter by role"
// @Param        search     query     string  false  "Search by name/email"
// @Param        active     query     bool    false  "Filter by status"
// @Success      200        {object}  ListUsersResponse
// @Router       /users/all [get]
func (h *Handler) GetUsers(c *gin.Context) {
	var params ListUsersParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.ListUsers(c.Request.Context(), params)
	if err != nil {
		log.Printf("[GetUsers] {ListUsers}: %v", err)
		response.SendError(
			c,
			"Failed to list users",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, resp)
}

// GetRoleDistribution godoc
// @Summary      Get user role distribution
// @Description  Returns the count of users for each role in the system.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200      {array}   RoleDistributionDTO
// @Router       /users/distribution [get]
func (h *Handler) GetRoleDistribution(c *gin.Context) {
	resp, err := h.service.GetRoleDistribution(c.Request.Context())
	if err != nil {
		log.Printf("[GetRoleDistribution] {GetRoleDistribution}: %v", err)
		response.SendError(
			c,
			"Failed to get role distribution",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, resp)
}

func (h *Handler) PostBlockUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		log.Printf(
			"[BlockUser] {Check Param ID}: User ID parameter is required",
		)
		response.SendFail(
			c,
			gin.H{"error": "User ID parameter is required"},
		)
		return
	}

	err := h.service.BlockUser(c.Request.Context(), userID)
	if err != nil {
		log.Printf("[BlockUser] {BlockUser}: %v", err)
		response.SendError(
			c,
			"Failed to block user",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, gin.H{"message": "User blocked successfully"})
}

func (h *Handler) PostUnblockUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		log.Printf(
			"[PostUnblockUser] {Check Param ID}: User ID parameter is required",
		)
		response.SendFail(
			c,
			gin.H{"error": "User ID parameter is required"},
		)
		return
	}

	err := h.service.UnblockUser(c.Request.Context(), userID)
	if err != nil {
		log.Printf("[UnblockUser] {UnblockUser}: %v", err)
		response.SendError(
			c,
			"Failed to unblock user",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, gin.H{"message": "User unblocked successfully"})
}
