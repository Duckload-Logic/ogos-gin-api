package bootstrap

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/storage"
)

type Application struct {
	Handlers *Handlers
}

func Initialize(db *sqlx.DB, cfg *config.Config) (*Application, error) {
	var fileStorage storage.FileStorage

	if cfg.IsProduction {
		blobStorage, err := storage.NewBlobStorage(
			cfg.AzureStorageConnectionString,
			cfg.AzureContainerName,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to initialize Azure Blob Storage: %w",
				err,
			)
		}
		fileStorage = blobStorage
	} else {
		uploadDir := cfg.LocalUploadDIR
		fileStorage = storage.NewDiskStorage(uploadDir)
	}

	repos := getRepositories(db)

	redis, err := datastore.NewRedisClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis: %w", err)
	}

	handlers := getHandlers(repos, fileStorage, cfg, redis)

	return &Application{
		Handlers: handlers,
	}, nil
}
