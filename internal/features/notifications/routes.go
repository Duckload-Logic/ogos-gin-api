package notifications

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/database"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

func RegisterRoutes(db *sqlx.DB, rg *gin.RouterGroup, h *Handler, redis *database.RedisClient) {
    routes := rg.Group("/notifications")
    routes.Use(middleware.AuthMiddleware(redis))

    userResourceLookup := middleware.OwnershipMiddleware(db, "userId")

    userRoutes := routes.Group("/")
    userRoutes.Use(middleware.RoleMiddleware(
        int(constants.StudentRoleID),
        int(constants.CounselorRoleID),
    ))
    {
        userRoutes.GET("/:userId", userResourceLookup, h.GetUserNotifications)
        
        userRoutes.PATCH("/:id/read", h.MarkAsRead) 
    }
}