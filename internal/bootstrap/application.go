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
	// initialize repos
	repos := GetRepositories(db)

	// initialize handlers
	handlers := GetHandlers(repos)

	// initialize router
	router := SetupRoutes(handlers)

	// get port
	portStr := os.Getenv("API_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid port in API_PORT %q: %w", portStr, err)
	}

	// return Application instance
	return &Application{
		Port:   port,
		Server: router,
	}, nil
}
