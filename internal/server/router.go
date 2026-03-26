package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/bootstrap"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/middleware"
	"github.com/olazo-johnalbert/duckload-api/internal/features/analytics"
	"github.com/olazo-johnalbert/duckload-api/internal/features/apikeys"
	"github.com/olazo-johnalbert/duckload-api/internal/features/appointments"
	"github.com/olazo-johnalbert/duckload-api/internal/features/auth"
	"github.com/olazo-johnalbert/duckload-api/internal/features/consents"
	"github.com/olazo-johnalbert/duckload-api/internal/features/locations"
	"github.com/olazo-johnalbert/duckload-api/internal/features/logs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/notes"
	"github.com/olazo-johnalbert/duckload-api/internal/features/notifications"
	"github.com/olazo-johnalbert/duckload-api/internal/features/slips"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students/external"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "github.com/olazo-johnalbert/duckload-api/docs"
	externalDocs "github.com/olazo-johnalbert/duckload-api/docs/external"
)

func NewRouter(
	db *sqlx.DB,
	handlers *bootstrap.Handlers,
	cfg *config.Config,
) *gin.Engine {
	g := gin.Default()

	localOrigins := []string{
		"http://localhost:8080",
		"http://127.0.0.1:8080",
		"http://localhost:5173",
	}

	prodOrigins := []string{
		"https://pupt-ogos.dllbsit2027.com",
		"https://lemon-field-0c62e2800.1.azurestaticapps.net",
	}

	var origins []string
	if cfg.IsProduction {
		gin.SetMode(gin.ReleaseMode)
		origins = prodOrigins
	} else {
		gin.SetMode(gin.DebugMode)
		origins = localOrigins
	}

	corsConfig := cors.Config{
		AllowOrigins: origins,
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
			"x-api-key",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	g.Use(cors.New(corsConfig))
	g.Use(func(c *gin.Context) {
		c.Set(middleware.SecurityLoggerContextKey, handlers.SystemLogService)
		c.Next()
	})

	limiter := middleware.NewIPRateLimiter(5, 30)
	g.Use(middleware.RateLimitMiddleware(limiter))

	apiV1Routes := g.Group("/api/v1")

	apiV1Routes.GET("/docs/internal/*any", func(c *gin.Context) {
		docs.SwaggerInfo.Host = c.Request.Host
		ginSwagger.WrapHandler(swaggerFiles.Handler,
			ginSwagger.InstanceName("internal"),
		)(c)
	})
	apiV1Routes.GET("/docs/external/*any", func(c *gin.Context) {
		externalDocs.SwaggerInfoexternal.Host = c.Request.Host
		ginSwagger.WrapHandler(swaggerFiles.Handler,
			ginSwagger.InstanceName("external"),
		)(c)
	})

	if !cfg.IsProduction {
		apiV1Routes.Static("./uploads", "./uploads")
	}

	apiV1Routes.GET("/", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{"message": "PUPT-OGOS API version 1.0 initialized"},
		)
	})
	apiV1Routes.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	auth.RegisterRoutes(apiV1Routes, handlers.AuthHandler, handlers.Redis)
	users.RegisterRoutes(db, apiV1Routes, handlers.UserHandler, handlers.Redis)
	locations.RegisterRoutes(
		apiV1Routes,
		handlers.LocationsHandler,
		handlers.Redis,
	)
	students.RegisterRoutes(
		db,
		apiV1Routes,
		handlers.StudentHandler,
		handlers.Redis,
	)
	appointments.RegisterRoutes(
		db,
		apiV1Routes,
		handlers.AppointmentHandler,
		handlers.Redis,
	)
	slips.RegisterRoutes(db, apiV1Routes, handlers.SlipHandler, handlers.Redis)
	analytics.RegisterRoutes(
		apiV1Routes,
		handlers.AnalyticsHandler,
		handlers.Redis,
	)
	apikeys.RegisterRoutes(apiV1Routes, handlers.APIKeyHandler, handlers.Redis)
	notifications.RegisterRoutes(
		db,
		apiV1Routes,
		handlers.NotificationsHandler,
		handlers.Redis,
	)
	logs.RegisterRoutes(apiV1Routes, handlers.SystemLogHandler, handlers.Redis)
	consents.RegisterRoutes(
		apiV1Routes,
		handlers.ConsentHandler,
		handlers.Redis,
	)
	notes.RegisterRoutes(db, apiV1Routes, handlers.NoteHandler, handlers.Redis)

	external.RegisterRoutes(
		apiV1Routes,
		handlers.ExternalStudentHandler,
		handlers.APIKeyService.ValidateKeyFunc(),
		handlers.Redis,
	)
	return g
}
