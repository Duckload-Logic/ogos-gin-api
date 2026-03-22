package notes

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/database"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

func RegisterRoutes(db *sqlx.DB, rg *gin.RouterGroup, h *Handler, redis *database.RedisClient) {
	routes := rg.Group("/notes")
	routes.Use(middleware.AuthMiddleware(redis))
	routes.Use(middleware.HydrateStudentContext(db))
	routes.Use(middleware.RoleMiddleware(int(constants.CounselorRoleID)))
	{
		routes.GET("/user/id/:iirID", h.HandleGetStudentSignificantNotes)
		routes.POST("/user/id/:iirID", h.HandlePostStudentSignificantNote)
	}
}
