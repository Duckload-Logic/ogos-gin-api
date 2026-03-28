package apikeys

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

// PostAPIKey creates a new API key and returns the plaintext key (shown only
// once).
func (h *Handler) PostAPIKey(c *gin.Context) {
	var req CreateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.GenerateKey(c.Request.Context(), req)
	if err != nil {
		log.Printf("[PostAPIKey] {GenerateKey}: %v", err)
		response.SendError(
			c,
			"Failed to create API key",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, resp)
}

// GetAPIKeys lists all API keys (prefix only, never the full key).
func (h *Handler) GetAPIKeys(c *gin.Context) {
	includeRevoked := c.Query("include_revoked") == "true"

	keys, err := h.service.ListKeys(c.Request.Context(), includeRevoked)
	if err != nil {
		response.SendError(
			c,
			"Failed to list API keys",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, keys)
}

// DeleteAPIKey deactivates an API key.
func (h *Handler) DeleteAPIKey(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.SendFail(c, gin.H{"error": "Invalid API key ID"})
		return
	}

	if err := h.service.RevokeKey(c.Request.Context(), id); err != nil {
		response.SendError(
			c,
			"Failed to revoke API key",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, gin.H{"message": "API key revoked successfully"})
}
