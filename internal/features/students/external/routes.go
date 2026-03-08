package external

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, validator middleware.APIKeyValidator) {
	routes := rg.Group("/students/external")
	routes.Use(middleware.APIKeyMiddleware(validator))

	routes.GET("/by-email", h.GetStudentByEmail)
	routes.GET("/personal-info", h.GetPersonalInfoByStudentNumber)
	routes.GET("/addresses", h.GetAddressByStudentNumber)
}
