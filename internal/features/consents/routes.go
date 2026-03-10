package consents

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	routes := rg.Group("/consents")
	routes.Use(middleware.AuthMiddleware())

	superAdminOnly := routes.Group("")
	superAdminOnly.Use(middleware.RoleMiddleware(
		int(constants.SuperAdminRoleID),
	))
	{
		superAdminOnly.POST("/upload", h.HandleUploadDocument)
	}

	all := routes.Group("")
	{
		all.GET("/latest/:type", h.HandleGetLatestDocument)
		all.GET("/check/email/:userEmail/doc/:docID", h.HandleCheckUserConsent)

		all.POST("/save", h.HandleSaveConsent)
	}
}
