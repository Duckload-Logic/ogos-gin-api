package external

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, validator middleware.APIKeyValidator) {
	routes := rg.Group("/students/external")
	routes.Use(middleware.APIKeyMiddleware(validator))

	routes.GET("/by-email/:email", h.HandleGetStudentByEmail)
	routes.GET("/personal-info/:studentNumber", h.HandleGetPersonalInfoByStudentNumber)
	routes.GET("/addresses/:studentNumber", h.HandleGetAddressByStudentNumber)
}
