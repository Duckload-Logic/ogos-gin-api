package m2mclients

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/response"
)

type Handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetService() ServiceInterface {
	return h.service
}

func (h *Handler) PostM2MClient(c *gin.Context) {
	var req CreateM2MClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(string)
	resp, err := h.service.CreateClient(c.Request.Context(), userID, req)
	if err != nil {
		if err.Error() == "developer already has an active m2m client" {
			response.SendError(c, err.Error(), http.StatusConflict, nil)
			return
		}
		fmt.Printf("[PostM2MClient] {Create Client}: %v\n", err)
		response.SendError(
			c,
			"Failed to create M2M client",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, resp)
}

// PostM2MToken godoc
// @Summary      M2M Token Exchange
// @Description  Exchanges client credentials (client_id and client_secret) for an M2M access token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body      M2MTokenRequest true "M2M Credentials"
// @Success      200     {object}  M2MTokenSuccessResponse
// @Failure      400     {object}  response.CommonErrorResponse
// @Failure      401     {object}  response.CommonErrorResponse
// @Router       /auth/m2m/token [post]
func (h *Handler) PostM2MToken(c *gin.Context) {
	var req M2MTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	client, err := h.service.Authenticate(
		c.Request.Context(),
		req.ClientID,
		req.ClientSecret,
	)
	if err != nil {
		response.SendFail(
			c,
			gin.H{"error": err.Error()},
			http.StatusUnauthorized,
		)
		return
	}

	tokens, err := h.service.IssueToken(c.Request.Context(), client)
	if err != nil {
		response.SendError(
			c,
			string(constants.ErrInternalServerError),
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, tokens)
}

func (h *Handler) PatchM2MClientVerify(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.SendFail(c, gin.H{"error": "invalid client ID"})
		return
	}

	err = h.service.VerifyClient(c.Request.Context(), id)
	if err != nil {
		fmt.Printf("[PatchM2MClientVerify] {Service Verify}: %v\n", err)
		response.SendError(
			c,
			string(constants.ErrInternalServerError),
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, gin.H{"message": "Client verified successfully"})
}

func (h *Handler) PostM2MTokenRefresh(c *gin.Context) {
	var req M2MRefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	tokenResp, err := h.service.RefreshToken(
		c.Request.Context(),
		req.RefreshToken,
	)
	if err != nil {
		fmt.Printf("[PostM2MTokenRefresh] {RefreshToken}: %v\n", err)
		response.SendError(c, err.Error(), http.StatusUnauthorized, nil)
		return
	}

	response.SendSuccess(c, tokenResp)
}

func (h *Handler) GetM2MClients(c *gin.Context) {
	includeRevoked := c.Query("include_revoked") == "true"
	roleIDs := c.MustGet("roleIDs").([]int)

	isSuperAdmin := false
	for _, rid := range roleIDs {
		if rid == int(constants.SuperAdminRoleID) {
			isSuperAdmin = true
			break
		}
	}

	var userID string
	if !isSuperAdmin {
		userID = c.MustGet("userID").(string)
	}

	clients, err := h.service.ListClients(
		c.Request.Context(),
		userID,
		includeRevoked,
		roleIDs,
	)
	if err != nil {
		fmt.Printf("[GetM2MClients] {List Clients}: %v\n", err)
		response.SendError(
			c,
			"Failed to list M2M clients",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, clients)
}

func (h *Handler) GetMyM2MClient(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	clients, err := h.service.ListClients(
		c.Request.Context(),
		userID,
		false,
		[]int{},
	)
	if err != nil {
		fmt.Printf("[GetMyM2MClient] {List Clients}: %v\n", err)
		response.SendError(
			c,
			"Failed to list M2M clients",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	if len(clients) == 0 {
		response.SendFail(c, gin.H{"error": "You do not have an M2M client"})
		return
	}

	response.SendSuccess(c, clients[0])
}

func (h *Handler) PostM2MClientSecret(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.SendFail(c, gin.H{"error": "Invalid M2M client ID"})
		return
	}

	secret, err := h.service.RegenerateSecret(c.Request.Context(), id)
	if err != nil {
		fmt.Printf("[PostM2MClientSecret] {Regenerate Secret}: %v\n", err)
		response.SendError(
			c,
			"Failed to regenerate client secret",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, gin.H{"clientSecret": secret})
}

func (h *Handler) DeleteM2MClient(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.SendFail(c, gin.H{"error": "Invalid M2M client ID"})
		return
	}

	if err := h.service.RevokeClient(c.Request.Context(), id); err != nil {
		fmt.Printf("[DeleteM2MClient] {Revoke Client}: %v\n", err)
		response.SendError(
			c,
			"Failed to revoke M2M client",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, gin.H{"message": "M2M client revoked successfully"})
}
