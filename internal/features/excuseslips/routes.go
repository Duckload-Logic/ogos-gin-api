package excuseslips

import (
    "github.com/gin-gonic/gin"
    "github.com/olazo-johnalbert/duckload-api/internal/core/constants"
    "github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

func RegisterRoutes(api *gin.RouterGroup, h *Handler) {
    excuseslipGroup := api.Group("/excuseslips")
    excuseslipGroup.Use(middleware.AuthMiddleware()) 

    excuseslipGroup.GET("", middleware.RoleMiddleware(
        int(constants.StudentRoleID), int(constants.CounselorRoleID), int(constants.FrontDeskRoleID),
    ), h.GetAll)

    excuseslipGroup.GET("/:id", middleware.RoleMiddleware(
        int(constants.StudentRoleID), int(constants.CounselorRoleID), int(constants.FrontDeskRoleID),
    ), h.GetByID)

    excuseslipGroup.POST("", middleware.RoleMiddleware(
        int(constants.StudentRoleID),
    ), h.Submit)

    excuseslipGroup.PATCH("/:id/status", middleware.RoleMiddleware(
        int(constants.CounselorRoleID), 
        int(constants.FrontDeskRoleID),
    ), h.UpdateStatus)
    
    excuseslipGroup.DELETE("/:id", middleware.RoleMiddleware(
        int(constants.CounselorRoleID), 
        int(constants.FrontDeskRoleID),
    ), h.Delete)
}