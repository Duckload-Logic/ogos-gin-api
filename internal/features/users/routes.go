package users

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

func RegisterRoutes(db *sql.DB, rg *gin.RouterGroup, h *Handler) {
	userRoutes := rg.Group("/users")
	userRoutes.Use(middleware.AuthMiddleware())
	userRoutes.Use(middleware.RoleMiddleware(
		int(constants.CounselorRoleID),
		int(constants.FrontDeskRoleID),
		int(constants.StudentRoleID),
	))
	userRoutes.Use(middleware.OwnershipMiddleware(db, "userID"))
	{
		userRoutes.GET("/id/:userID", h.HandleGetUserByID)
	}
}
