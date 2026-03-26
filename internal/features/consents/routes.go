package consents

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/middleware"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

func RegisterRoutes(
	rg *gin.RouterGroup,
	h *Handler,
	redis *datastore.RedisClient,
) {
	consentRoutes := rg.Group("/consents")
	consentRoutes.Use(middleware.AuthMiddleware(redis))

	superAdminOnly := consentRoutes.Group("")
	superAdminOnly.Use(middleware.RoleMiddleware(
		int(constants.SuperAdminRoleID),
	))
	{
		superAdminOnly.POST("/upload", h.PostDocument)
	}

	all := consentRoutes.Group("")
	{
		all.GET("/latest/:type", h.GetLatestDocument)
		all.GET("/latest/:type/content", h.GetDocumentContent)
		all.GET("/check/doc/:docID", h.GetConsentCheck)

		all.POST("/:type/doc/:docID", h.PostConsent)
	}
}
