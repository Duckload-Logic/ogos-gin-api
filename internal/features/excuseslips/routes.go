package excuseslips

import (
    "github.com/gin-gonic/gin"
    "github.com/olazo-johnalbert/duckload-api/internal/core/constants"
    "github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

func RegisterRoutes(api *gin.RouterGroup, h *Handler) {
    excuseslipGroup := api.Group("/excuseslips")
    excuseslipGroup.Use(middleware.RoleMiddleware(
        int(constants.StudentRoleID),
        int(constants.CounselorRoleID),
        int(constants.FrontDeskRoleID),
    ))
    {
        excuseslipGroup.GET("", h.GetAll)
        excuseslipGroup.GET("/:id", h.GetByID)
        excuseslipGroup.POST("", h.Submit)
        excuseslipGroup.PATCH("/:id/status", h.UpdateStatus)
        excuseslipGroup.DELETE("/:id", h.Delete)
    }
}