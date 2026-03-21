package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/database"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, redis *database.RedisClient) {
	authRoutes := rg.Group("/auth")
	{
		authRoutes.POST("/login", h.HandleLogin)
		authRoutes.POST("/refresh", h.HandleRefreshToken)
		authRoutes.GET(
			"/me",
			middleware.AuthMiddleware(redis),
			h.HandleGetMe,
		)
		authRoutes.POST(
			"/logout",
			middleware.AuthMiddleware(redis),
			h.HandleLogout,
		)

		// IDP OAuth 2.0 routes
		authRoutes.GET("/idp/authorize-url", h.GetAuthorizeURL)
		authRoutes.POST("/idp/token", h.PostIDPToken)
	}
}
