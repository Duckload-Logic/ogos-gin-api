package analytics

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/response"
)

type Handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAnalyticsDashboard(c *gin.Context) {
	ctx := c.Request.Context()

	dashboardData, err := h.service.GetDashboard(ctx)
	if err != nil {
		response.SendError(
			c,
			"Failed to generate analytics dashboard",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, dashboardData)
}
