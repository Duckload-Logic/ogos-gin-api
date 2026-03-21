package bootstrap

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/storage"
	"github.com/olazo-johnalbert/duckload-api/internal/database"
)

type Application struct {
	Port   int
	Server *gin.Engine
}

func GetNewApplication(db *sqlx.DB, cfg *config.Config) (*Application, error) {
	var fileStorage storage.FileStorage

	if cfg.IsProduction {
		blobStorage, err := storage.NewBlobStorage(
			cfg.AzureStorageConnectionString,
			cfg.AzureContainerName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Azure Blob Storage: %w", err)
		}
		fileStorage = blobStorage
	} else {
		uploadDir := cfg.LocalUploadDIR
		fileStorage = storage.NewDiskStorage(uploadDir)
	}

	// Set Gin mode based on environment
	if cfg.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	repos := getRepositories(db)

	redis, err := database.NewRedisClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis: %w", err)
	}

	handlers := getHandlers(repos, fileStorage, cfg, redis)

	router := SetupRoutes(db, handlers, cfg)

	port, err := strconv.Atoi(cfg.WebsitesPort)
	if err != nil {
		return nil, fmt.Errorf("invalid port in WEBSITES_PORT %q: %w", cfg.WebsitesPort, err)
	}

	return &Application{
		Port:   port,
		Server: router,
	}, nil
}
