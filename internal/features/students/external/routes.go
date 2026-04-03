package external

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/middleware"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

func RegisterRoutes(
	rg *gin.RouterGroup,
	h *Handler,
	validator middleware.APIKeyValidator,
	redis *datastore.RedisClient,
) {
	routes := rg.Group("/students/external")
	routes.Use(middleware.APIKeyMiddleware(validator))

	routes.GET("/profiles", h.GetStudents)
	routes.GET("/:studentNumber", h.GetStudentByStudentNumber)
	routes.GET(
		"/:studentNumber/personal-info",
		h.GetPersonalInfoByStudentNumber,
	)
	routes.GET("/:studentNumber/addresses", h.GetAddressByStudentNumber)
}
