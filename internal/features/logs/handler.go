package logs

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
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

// HandleListSystemLogs godoc
// @Summary      List system logs
// @Description  Retrieves a paginated list of audit, system, and security logs.
// Super Admin only.
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
func (h *Handler) GetLogs(c *gin.Context) {
	var req ListSystemLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Printf("[GetLogs] {Bind Query}: %v", err)
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.ListLogs(c.Request.Context(), req)
	if err != nil {
		log.Printf("[GetLogs] {ListLogs}: %v", err)
		response.SendError(
			c,
			"Failed to retrieve system logs",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, result)
}

// GetAuditLogs returns only AUDIT category logs
func (h *Handler) GetAuditLogs(c *gin.Context) {
	var req ListSystemLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Printf("[GetAuditLogs] {Bind Query}: %v", err)
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	req.Category = audit.CategoryAudit

	result, err := h.service.ListLogs(c.Request.Context(), req)
	if err != nil {
		log.Printf("[GetAuditLogs] {ListLogs}: %v", err)
		response.SendError(
			c,
			"Failed to retrieve audit logs",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, result)
}

// GetSystemLogs returns only SYSTEM category logs
func (h *Handler) GetSystemLogs(c *gin.Context) {
	var req ListSystemLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Printf("[GetSystemLogs] {Bind Query}: %v", err)
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	req.Category = audit.CategorySystem

	result, err := h.service.ListLogs(c.Request.Context(), req)
	if err != nil {
		log.Printf("[GetSystemLogs] {ListLogs}: %v", err)
		response.SendError(
			c,
			"Failed to retrieve system logs",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, result)
}

// GetSecurityLogs returns only SECURITY category logs
func (h *Handler) GetSecurityLogs(c *gin.Context) {
	var req ListSystemLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Printf("[GetSecurityLogs] {Bind Query}: %v", err)
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	req.Category = audit.CategorySecurity

	result, err := h.service.ListLogs(c.Request.Context(), req)
	if err != nil {
		log.Printf("[GetSecurityLogs] {ListLogs}: %v", err)
		response.SendError(
			c,
			"Failed to retrieve security logs",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, result)
}

// GetLogStats returns log counts by category
func (h *Handler) GetLogStats(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	stats, err := h.service.GetStats(c.Request.Context(), startDate, endDate)
	if err != nil {
		log.Printf("[GetLogStats] {GetStats}: %v", err)
		response.SendError(
			c,
			"Failed to retrieve log stats",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, stats)
}
