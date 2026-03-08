package apikeys

import (
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

// HandleCreateAPIKey creates a new API key and returns the plaintext key (shown only once).
func (h *Handler) HandleCreateAPIKey(c *gin.Context) {
	var req CreateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.GenerateKey(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create API key"})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// HandleListAPIKeys lists all API keys (prefix only, never the full key).
func (h *Handler) HandleListAPIKeys(c *gin.Context) {
	includeRevoked := c.Query("include_revoked") == "true"

	keys, err := h.service.ListKeys(c.Request.Context(), includeRevoked)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list API keys"})
		return
	}

	c.JSON(http.StatusOK, keys)
}

// HandleRevokeAPIKey deactivates an API key.
func (h *Handler) HandleRevokeAPIKey(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid API key ID"})
		return
	}

	if err := h.service.RevokeKey(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke API key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API key revoked successfully"})
}
