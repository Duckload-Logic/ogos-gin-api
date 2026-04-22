package users

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/middleware"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

func RegisterRoutes(
	db *sqlx.DB,
	rg *gin.RouterGroup,
	h *Handler,
	redis *datastore.RedisClient,
) {
	userRoutes := rg.Group("/users")
	userRoutes.Use(middleware.AuthMiddleware(redis))
	userRoutes.Use(middleware.RoleMiddleware(
		int(constants.SuperAdminRoleID),
		int(constants.AdminRoleID),
		int(constants.StudentRoleID),
	))

	userRoutes.GET("/me", h.GetMe)
	userRoutes.GET("", h.GetUserByEmail)
	userRoutes.GET("/all",
		middleware.RoleMiddleware(int(constants.SuperAdminRoleID)),
		h.GetUsers,
	)
	userRoutes.GET("/distribution",
		middleware.RoleMiddleware(int(constants.SuperAdminRoleID)),
		h.GetRoleDistribution,
	)
	userRoutes.POST("/:id/block",
		middleware.RoleMiddleware(int(constants.SuperAdminRoleID)),
		h.PostUserBlock,
	)
	userRoutes.POST("/:id/unblock",
		middleware.RoleMiddleware(int(constants.SuperAdminRoleID)),
		h.PostUserUnblock,
	)

	// Session & Activity Audit (Super Admin only)
	userRoutes.GET("/:id/sessions",
		middleware.RoleMiddleware(int(constants.SuperAdminRoleID)),
		h.GetUserSessions,
	)
	userRoutes.DELETE("/:id/sessions/:session_id",
		middleware.RoleMiddleware(int(constants.SuperAdminRoleID)),
		h.DeleteUserSession,
	)
	userRoutes.GET("/:id/activity",
		middleware.RoleMiddleware(int(constants.SuperAdminRoleID)),
		h.GetUserActivity,
	)
}
