package excuseslips

import 	"github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup, h *Handler) {
	api.POST("/excuseslips", h.Submit)
}