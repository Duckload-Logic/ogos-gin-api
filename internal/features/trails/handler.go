package trails

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

// HandleListAuditTrails godoc
// @Summary      List audit trail entries
// @Description  Retrieves a paginated list of audit trail entries with optional filters.
// @Tags         AuditTrails
// @Accept       json
// @Produce      json
// @Param        page        query     int    false "Page number"
// @Param        page_size   query     int    false "Number of entries per page"
// @Param        action      query     string false "Filter by action (CREATE, UPDATE, DELETE)"
// @Param        entity_type query     string false "Filter by entity type (e.g., appointment, slip)"
// @Param        entity_id   query     int    false "Filter by entity ID"
// @Param        user_id     query     int    false "Filter by user ID"
// @Param        start_date  query     string false "Filter from date (YYYY-MM-DD)"
// @Param        end_date    query     string false "Filter to date (YYYY-MM-DD)"
// @Param        search      query     string false "Search by user name or entity type"
// @Success      200         {object}  ListAuditTrailsDTO
// @Failure      400         {object}  map[string]string "Bad request"
// @Failure      500         {object}  map[string]string "Internal server error"
// @Router       /audit-trails [get]
func (h *Handler) HandleListAuditTrails(c *gin.Context) {
	var req ListAuditTrailsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.ListAuditTrails(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve audit trails"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// HandleGetEntityAuditTrails godoc
// @Summary      Get audit trails for a specific entity
// @Description  Retrieves all audit trail entries for a given entity type and ID.
// @Tags         AuditTrails
// @Accept       json
// @Produce      json
// @Param        entityType path     string true "Entity type (e.g., appointment, slip)"
// @Param        entityId   path     int    true "Entity ID"
// @Success      200        {array}  AuditTrailDTO
// @Failure      400        {object} map[string]string "Bad request"
// @Failure      500        {object} map[string]string "Internal server error"
// @Router       /audit-trails/{entityType}/{entityId} [get]
func (h *Handler) HandleGetEntityAuditTrails(c *gin.Context) {
	entityType := c.Param("entityType")
	if entityType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Entity type is required"})
		return
	}

	entityID, err := strconv.Atoi(c.Param("entityId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity ID"})
		return
	}

	trails, err := h.service.GetByEntity(c.Request.Context(), entityType, entityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve audit trails"})
		return
	}

	c.JSON(http.StatusOK, trails)
}
