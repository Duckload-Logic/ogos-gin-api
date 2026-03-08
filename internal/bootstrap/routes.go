package bootstrap

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/features/analytics"
	"github.com/olazo-johnalbert/duckload-api/internal/features/apikeys"
	"github.com/olazo-johnalbert/duckload-api/internal/features/appointments"
	"github.com/olazo-johnalbert/duckload-api/internal/features/auth"
	"github.com/olazo-johnalbert/duckload-api/internal/features/locations"
	"github.com/olazo-johnalbert/duckload-api/internal/features/logs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/notifications"
	"github.com/olazo-johnalbert/duckload-api/internal/features/slips"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students/external"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	"github.com/olazo-johnalbert/duckload-api/internal/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/olazo-johnalbert/duckload-api/docs"
)

func SetupRoutes(db *sqlx.DB, handlers *Handlers) *gin.Engine {
	g := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:8080",
		"http://127.0.0.1:8080",
		"http://localhost:5173",
	}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Authorization")

	g.Use(cors.New(corsConfig))

	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	g.Use(func(c *gin.Context) {
		c.Set(middleware.SecurityLoggerContextKey, handlers.SystemLogService)
		c.Next()
	})

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
	locations.RegisterRoutes(apiV1Routes, handlers.LocationsHandler)
	students.RegisterRoutes(db, apiV1Routes, handlers.StudentHandler)
	appointments.RegisterRoutes(apiV1Routes, handlers.AppointmentHandler)
	slips.RegisterRoutes(apiV1Routes, handlers.SlipHandler)
	analytics.RegisterRoutes(apiV1Routes, handlers.AnalyticsHandler)
	apikeys.RegisterRoutes(apiV1Routes, handlers.APIKeyHandler)
	notifications.RegisterRoutes(db, apiV1Routes, handlers.NotificationsHandler)
	logs.RegisterRoutes(apiV1Routes, handlers.SystemLogHandler)

	external.RegisterRoutes(apiV1Routes, handlers.ExternalStudentHandler, handlers.APIKeyService.ValidateKeyFunc())

	return g
}
