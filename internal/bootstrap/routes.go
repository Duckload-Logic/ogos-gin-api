package bootstrap

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(handlers *Handlers) *gin.Engine {
	g := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:8080",
		"http://127.0.0.1:8080",
	}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Authorization")

	g.Use(cors.New(corsConfig))

	apiV1Routes := g.Group("/api/v1")

	// ==============================
	// |                            |
	// |        HOME ROUTES         |
	// |                            |
	// ==============================
	apiV1Routes.GET("/", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{"message": "PUPT-OGOS API version 1.0 initialized"},
		)
	})
	apiV1Routes.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Module routes should be below here
	// ==============================
	// |                            |
	// |       MODULE ROUTES        |
	// |                            |
	// ==============================

	apiV1Routes.POST("/appointments", handlers.AppointmentHandler.Create)

	return g
}
