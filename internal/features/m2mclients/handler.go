package m2mclients

import (
	"log"
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

// PostM2MClient creates a new M2M client and returns the client and secret.
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
		log.Printf("[PostM2MClient] {CreateClient}: %v", err)
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
		log.Printf("[PostM2MToken] {Authenticate}: %v", err)
		response.SendError(c, err.Error(), http.StatusUnauthorized, nil)
		return
	}

	tokenResp, err := h.service.IssueToken(c.Request.Context(), client)
	if err != nil {
		log.Printf("[PostM2MToken] {IssueToken}: %v", err)
		response.SendError(
			c,
			"Failed to issue token",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, tokenResp)
}

// PostM2MRefresh godoc
// @Summary      M2M Token Refresh
// @Description  Refreshes an existing M2M session using a valid refresh token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body      M2MRefreshTokenRequest true "Refresh Token"
// @Success      200     {object}  M2MTokenSuccessResponse
// @Failure      400     {object}  response.CommonErrorResponse
// @Failure      401     {object}  response.CommonErrorResponse
// @Router       /auth/m2m/refresh [post]
func (h *Handler) PostM2MRefresh(c *gin.Context) {
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
		log.Printf("[PostM2MRefresh] {RefreshToken}: %v", err)
		response.SendError(c, err.Error(), http.StatusUnauthorized, nil)
		return
	}

	response.SendSuccess(c, tokenResp)
}

// GetM2MClients lists all M2M clients.
func (h *Handler) GetM2MClients(c *gin.Context) {
	includeRevoked := c.Query("include_revoked") == "true"
	roleID := c.MustGet("roleID").(int)

	var userID string
	if roleID != int(constants.SuperAdminRoleID) {
		userID = c.MustGet("userID").(string)
	}

	clients, err := h.service.ListClients(
		c.Request.Context(),
		userID,
		includeRevoked,
	)
	if err != nil {
		log.Printf("[GetM2MClients] {ListClients}: %v", err)
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

// PostM2MSecret rotates the secret for an existing M2M client.
func (h *Handler) PostM2MSecret(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.SendFail(c, gin.H{"error": "Invalid M2M client ID"})
		return
	}

	secret, err := h.service.RegenerateSecret(c.Request.Context(), id)
	if err != nil {
		log.Printf("[PostM2MSecret] {RegenerateSecret}: %v", err)
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

// DeleteM2MClient deactivates an M2M client.
func (h *Handler) DeleteM2MClient(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.SendFail(c, gin.H{"error": "Invalid M2M client ID"})
		return
	}

	if err := h.service.RevokeClient(c.Request.Context(), id); err != nil {
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
