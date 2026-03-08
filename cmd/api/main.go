package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/olazo-johnalbert/duckload-api/internal/bootstrap"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/database"
)

// @title           DuckLoad API
// @version         1.0
// @BasePath        /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-api-key

// @securityDefinitions.authorization BearerAuth
// @in header
// @name Authorization
func main() {
	config := config.LoadConfig()
	dbUrl := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)

	if config.DBTLS {
		dbUrl += "&tls=true"
	}

	db, err := database.GetDBConnection(dbUrl)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Get application
	app, err := bootstrap.GetNewApplication(db)
	if err != nil {
		log.Fatal(err)
	}

	// Serve Application
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}
