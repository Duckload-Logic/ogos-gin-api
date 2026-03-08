package bootstrap

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/core/azure"
	"github.com/olazo-johnalbert/duckload-api/internal/core/storage"
)

type Application struct {
	Port   int
	Server *gin.Engine
}

func GetNewApplication(db *sqlx.DB) (*Application, error) {

	var fileStorage storage.FileStorage

	azureConnStr := os.Getenv("AZURE_STORAGE_CONNECTION_STRING")
	if azureConnStr != "" {
		blobStorage, err := azure.NewBlobStorage(
			azureConnStr,
			os.Getenv("AZURE_CONTAINER_NAME"),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Azure Blob Storage: %w", err)
		}
		fileStorage = blobStorage
	} else {
		uploadDir := os.Getenv("UPLOAD_DIR")
		if uploadDir == "" {
			uploadDir = "./uploads"
		}
		fileStorage = storage.NewDiskStorage(uploadDir)
	}

	repos := getRepositories(db)

	handlers := getHandlers(repos, fileStorage)

	router := SetupRoutes(db, handlers)

	portStr := os.Getenv("API_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid port in API_PORT %q: %w", portStr, err)
	}

	return &Application{
		Port:   port,
		Server: router,
	}, nil
}
