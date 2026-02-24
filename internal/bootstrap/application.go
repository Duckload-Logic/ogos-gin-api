package bootstrap

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Application struct {
	Port   int
	Server *gin.Engine
}

func GetNewApplication(db *sqlx.DB) (*Application, error) {

	repos := getRepositories(db)

	handlers := getHandlers(repos)

	router := SetupRoutes(db, handlers)
	router.Static("/uploads", "./uploads")

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
