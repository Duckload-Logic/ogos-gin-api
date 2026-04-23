package users

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/response"
	"github.com/olazo-johnalbert/duckload-api/internal/core/sessions"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

type Handler struct {
	service        ServiceInterface
	sessionService *sessions.Service
	logger         audit.Logger
	logReader      audit.LogReader
}

// NewHandler creates a new users handler.
func NewHandler(
	service ServiceInterface,
	sessionService *sessions.Service,
	logger audit.Logger,
	logReader audit.LogReader,
) *Handler {
	return &Handler{
		service:        service,
		sessionService: sessionService,
		logger:         logger,
		logReader:      logReader,
	}
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
		fmt.Printf("[GetMe] {GetUserByID}: %v\n", err)
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
		fmt.Printf(
			"[GetUserByEmail] {Check Query}: Email parameter is required\n",
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
		fmt.Printf("[GetUserByEmail] {GetUserByEmail}: %v\n", err)
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
		fmt.Printf("[GetUsers] {ListUsers}: %v\n", err)
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
		fmt.Printf("[GetRoleDistribution] {GetRoleDistribution}: %v\n", err)
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

func (h *Handler) PostUserBlock(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		fmt.Printf(
			"[PostUserBlock] {Check Param}: User ID is required\n",
		)
		response.SendFail(
			c,
			gin.H{"error": "User ID parameter is required"},
		)
		return
	}

	err := h.service.BlockUser(c.Request.Context(), userID)
	if err != nil {
		fmt.Printf("[PostUserBlock] {BlockUser}: %v\n", err)
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

func (h *Handler) PostUserUnblock(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		fmt.Printf(
			"[PostUserUnblock] {Check Param}: User ID is required\n",
		)
		response.SendFail(
			c,
			gin.H{"error": "User ID parameter is required"},
		)
		return
	}

	err := h.service.UnblockUser(c.Request.Context(), userID)
	if err != nil {
		fmt.Printf("[PostUserUnblock] {UnblockUser}: %v\n", err)
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

func (h *Handler) GetUserSessions(c *gin.Context) {
	targetUserID := c.Param("id")
	if targetUserID == "" {
		response.SendFail(c, gin.H{"error": "User ID is required"})
		return
	}

	sessions, err := h.sessionService.ListUserSessions(
		c.Request.Context(),
		targetUserID,
	)
	if err != nil {
		fmt.Printf("[GetUserSessions] {ListSessions}: %v\n", err)
		response.SendError(
			c,
			"Failed to list user sessions",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, sessions)
}

func (h *Handler) DeleteUserSession(c *gin.Context) {
	targetUserID := c.Param("id")
	jti := c.Param("session_id")
	if targetUserID == "" || jti == "" {
		response.SendFail(
			c,
			gin.H{"error": "User ID and Session ID are required"},
		)
		return
	}

	err := h.sessionService.DeleteUserToken(
		c.Request.Context(),
		targetUserID,
		sessions.NewJTI(jti),
	)
	if err != nil {
		fmt.Printf("[DeleteUserSession] {DeleteToken}: %v\n", err)
		response.SendError(
			c,
			"Failed to revoke session",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	// Audit log for session revocation
	adminEmail := c.MustGet("userEmail").(string)
	adminID := c.MustGet("userID").(string)

	h.logger.Record(
		c.Request.Context(),
		nil, // Use default DB in logger implementation
		audit.LogEntry{
			Level:    audit.LevelWarning,
			Category: audit.CategorySecurity,
			Action:   audit.ActionLogout, // Or add a new ActionSessionRevoked
			Message: fmt.Sprintf(
				"Superadmin %s revoked session %s for user %s",
				adminEmail,
				jti,
				targetUserID,
			),
			UserID:   structs.StringToNullableString(adminID),
			TargetID: structs.StringToNullableString(targetUserID),
		},
	)

	response.SendSuccess(c, gin.H{"message": "Session revoked successfully"})
}

func (h *Handler) GetUserActivity(c *gin.Context) {
	targetUserID := c.Param("id")
	if targetUserID == "" {
		response.SendFail(c, gin.H{"error": "User ID is required"})
		return
	}

	// Fetch user email first because the logs repo mostly filters by email
	user, err := h.service.GetUserByID(c.Request.Context(), targetUserID)
	if err != nil {
		fmt.Printf("[GetUserActivity] {GetUserByID}: %v\n", err)
		response.SendError(c, "Failed to find user", http.StatusNotFound, nil)
		return
	}

	var req audit.ListSystemLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	// Force filter by target user's email
	req.UserEmail = user.Email

	result, err := h.logReader.ListLogs(c.Request.Context(), req)
	if err != nil {
		fmt.Printf("[GetUserActivity] {ListLogs}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve user activity",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, result)
}

func (h *Handler) PostUpdateRoles(c *gin.Context) {
	var req UpdateRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	adminID := c.MustGet("userID").(string)

	err := h.service.UpdateUserRoles(c.Request.Context(), req, adminID)
	if err != nil {
		fmt.Printf("[PostUpdateRoles] {UpdateUserRoles}: %v\n", err)
		response.SendError(
			c,
			"Failed to update user roles",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	// Audit log for security visibility
	adminEmail := c.MustGet("userEmail").(string)
	h.logger.Record(
		c.Request.Context(),
		nil,
		audit.LogEntry{
			Level:    audit.LevelWarning,
			Category: audit.CategorySecurity,
			Action:   audit.ActionElevateRoles,
			Message: fmt.Sprintf(
				"Superadmin %s elevated/updated roles for user %s. Reason: %s, Ref: %s",
				adminEmail,
				req.UserID,
				req.Reason,
				req.ReferenceID,
			),
			UserID:   structs.StringToNullableString(adminID),
			TargetID: structs.StringToNullableString(req.UserID),
		},
	)

	response.SendSuccess(c, gin.H{"message": "User roles updated successfully"})
}

// func (h *Handler) PostProfilePicture(c *gin.Context) {
// 	userID := c.Param("id")
// 	if userID == "" {
// 		response.SendFail(c, gin.H{"error": "User ID is required"})
// 		return
// 	}

// 	var req struct {
// 		FileID string `json:"file_id" binding:"required"`
// 	}

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		response.SendFail(c, gin.H{"error": err.Error()})
// 		return
// 	}

// 	err := h.service.PostProfilePicture(c.Request.Context(), userID, req.FileID)
// 	if err != nil {
// 		fmt.Printf("[PostProfilePicture] {PostProfilePicture}: %v\n", err)
// 		response.SendError(
// 			c,
// 			"Failed to set profile picture",
// 			http.StatusInternalServerError,
// 			nil,
// 		)
// 		return
// 	}

// 	response.SendSuccess(
// 		c,
// 		gin.H{"message": "Profile picture set successfully"},
// 	)
// }
