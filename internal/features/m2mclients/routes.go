package m2mclients

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/middleware"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

// RegisterRoutes sets up management and auth endpoints for M2M clients.
func RegisterRoutes(
	rg *gin.RouterGroup,
	h *Handler,
	redis *datastore.RedisClient,
) {
	// Public M2M Auth Routes
	m2mAuth := rg.Group("/auth/m2m")
	{
		m2mAuth.POST("/token", h.PostM2MToken)
		m2mAuth.POST("/refresh", h.PostM2MRefresh)
	}

	// Protected Management Routes
	m2mMgmt := rg.Group("/m2m-clients")
	m2mMgmt.Use(middleware.AuthMiddleware(redis))
	m2mMgmt.Use(middleware.RoleMiddleware(
		int(constants.SuperAdminRoleID),
		int(constants.DeveloperRoleID),
	))
	{
		m2mMgmt.GET("", h.GetM2MClients)
		m2mMgmt.POST("", h.PostM2MClient)
		m2mMgmt.POST("/:id/secret", h.PostM2MSecret)
		m2mMgmt.DELETE("/:id", h.DeleteM2MClient)
	}
}
