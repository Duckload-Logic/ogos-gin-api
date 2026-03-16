package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	authRoutes := rg.Group("/auth")
	{
		authRoutes.GET("/login", h.GetAuthRedirect)
		authRoutes.GET("/callback", h.GetAuthCallback)
		authRoutes.POST("/refresh", h.HandleRefreshToken)
		authRoutes.POST("/logout", 
			middleware.AuthMiddleware(), h.HandleLogout)
	}
}