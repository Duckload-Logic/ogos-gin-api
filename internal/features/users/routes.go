package users

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

func RegisterRoutes(db *sqlx.DB, rg *gin.RouterGroup, h *Handler) {
	userRoutes := rg.Group("/users")
	userRoutes.Use(middleware.AuthMiddleware())
	userRoutes.Use(middleware.RoleMiddleware(
		int(constants.SuperAdminRoleID),
		int(constants.CounselorRoleID),
		int(constants.StudentRoleID),
	))
	userRoutes.GET("/me", h.HandleGetCurrentUser)
}
