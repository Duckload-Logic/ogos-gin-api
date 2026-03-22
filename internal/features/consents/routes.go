package consents

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/database"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, redis *database.RedisClient) {
	consentRoutes := rg.Group("/consents")
	consentRoutes.Use(middleware.AuthMiddleware(redis))

	superAdminOnly := consentRoutes.Group("")
	superAdminOnly.Use(middleware.RoleMiddleware(
		int(constants.SuperAdminRoleID),
	))
	{
		superAdminOnly.POST("/upload", h.HandleUploadDocument)
	}

	all := consentRoutes.Group("")
	{
		all.GET("/latest/:type", h.HandleGetLatestDocument)
		all.GET("/latest/:type/content", h.HandleGetDocumentContent)
		all.GET("/check/doc/:docID", h.HandleCheckUserConsent)

		all.POST("/:type/doc/:docID", h.HandleSaveConsent)
	}
}
