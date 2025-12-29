package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	authRoutes := rg.Group("/auth")
	{
		authRoutes.POST("/login", h.HandleLogin)
		authRoutes.POST("/refresh", h.HandleRefreshToken)
	}
}
