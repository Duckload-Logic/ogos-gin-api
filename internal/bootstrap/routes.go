package bootstrap

import (
	"database/sql"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/features/appointments"
	"github.com/olazo-johnalbert/duckload-api/internal/features/auth"
	"github.com/olazo-johnalbert/duckload-api/internal/features/excuseslips"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/olazo-johnalbert/duckload-api/docs"
)

func SetupRoutes(db *sql.DB, handlers *Handlers) *gin.Engine {
	g := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:8080",
		"http://127.0.0.1:8080",
		"http://localhost:5173", // Address for Vite dev server
	}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Authorization")

	g.Use(cors.New(corsConfig))

	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	limiter := middleware.NewIPRateLimiter(5, 30)
	g.Use(middleware.RateLimitMiddleware(limiter))

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

	// ==============================
	// |                            |
	// |       MODULE ROUTES        |
	// |                            |
	// ==============================

	auth.RegisterRoutes(apiV1Routes, handlers.AuthHandler)
	users.RegisterRoutes(db, apiV1Routes, handlers.UserHandler)
	students.RegisterRoutes(db, apiV1Routes, handlers.StudentHandler)
	appointments.RegisterRoutes(apiV1Routes, handlers.AppointmentHandler)
	excuseslips.RegisterRoutes(apiV1Routes, handlers.ExcuseSlipHandler)

	return g
}
