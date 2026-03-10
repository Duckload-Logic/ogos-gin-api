package external

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, validator middleware.APIKeyValidator) {
	routes := rg.Group("/students/external")
	routes.Use(middleware.APIKeyMiddleware(validator))

	routes.GET("", h.ListStudents)
	routes.GET("/by-id/:userID", h.HandleGetStudentByUserID)
	routes.GET("/personal-info/:studentNumber", h.HandleGetPersonalInfoByStudentNumber)
	routes.GET("/addresses/:studentNumber", h.HandleGetAddressByStudentNumber)
}
