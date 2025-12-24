package bootstrap

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Application struct {
	Port   int
	Server *gin.Engine
}

func GetNewApplication(db *sql.DB) (*Application, error) {

	repos := getRepositories(db)

	handlers := getHandlers(repos)

	router := SetupRoutes(handlers)

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
