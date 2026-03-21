package logs

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// HandleListSystemLogs godoc
// @Summary      List system logs
// @Description  Retrieves a paginated list of audit, system, and security logs. Super Admin only.
// @Tags         SystemLogs
// @Accept       json
// @Produce      json
// @Param        page        query     int    false "Page number"
// @Param        page_size   query     int    false "Number of entries per page"
// @Param        category    query     string false "Filter by category (AUDIT, SYSTEM, SECURITY)"
// @Param        action      query     string false "Filter by action"
// @Param        user_email  query     string false "Filter by user email"
// @Param        start_date  query     string false "Filter from date (YYYY-MM-DD)"
// @Param        end_date    query     string false "Filter to date (YYYY-MM-DD)"
// @Param        search      query     string false "Search in message, action, or user email"
// @Success      200         {object}  ListSystemLogsDTO
// @Failure      400         {object}  map[string]string "Bad request"
// @Failure      500         {object}  map[string]string "Internal server error"
// @Router       /system-logs [get]
func (h *Handler) HandleListSystemLogs(c *gin.Context) {
	var req ListSystemLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Printf("[HandleListSystemLogs] {Bind Query}: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.ListLogs(c.Request.Context(), req)
	if err != nil {
		log.Printf("[HandleListSystemLogs] {ListLogs}: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve system logs"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// HandleGetAuditLogs returns only AUDIT category logs
func (h *Handler) HandleGetAuditLogs(c *gin.Context) {
	var req ListSystemLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Printf("[HandleGetAuditLogs] {Bind Query}: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Category = CategoryAudit

	result, err := h.service.ListLogs(c.Request.Context(), req)
	if err != nil {
		log.Printf("[HandleGetAuditLogs] {ListLogs}: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve audit logs"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// HandleGetSystemLogs returns only SYSTEM category logs
func (h *Handler) HandleGetSystemLogs(c *gin.Context) {
	var req ListSystemLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Printf("[HandleGetSystemLogs] {Bind Query}: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Category = CategorySystem

	result, err := h.service.ListLogs(c.Request.Context(), req)
	if err != nil {
		log.Printf("[HandleGetSystemLogs] {ListLogs}: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve system logs"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// HandleGetSecurityLogs returns only SECURITY category logs
func (h *Handler) HandleGetSecurityLogs(c *gin.Context) {
	var req ListSystemLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Printf("[HandleGetSecurityLogs] {Bind Query}: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Category = CategorySecurity

	result, err := h.service.ListLogs(c.Request.Context(), req)
	if err != nil {
		log.Printf("[HandleGetSecurityLogs] {ListLogs}: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve security logs"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// HandleGetLogStats returns log counts by category
func (h *Handler) HandleGetLogStats(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	stats, err := h.service.GetStats(c.Request.Context(), startDate, endDate)
	if err != nil {
		log.Printf("[HandleGetLogStats] {GetStats}: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve log stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
