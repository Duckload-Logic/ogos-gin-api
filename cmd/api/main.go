package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/olazo-johnalbert/duckload-api/internal/bootstrap"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
	"github.com/olazo-johnalbert/duckload-api/internal/server"
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
		config.DBPass,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)

	if config.DBTLS {
		dbUrl += "&tls=true"
	}

	db, err := datastore.GetDBConnection(dbUrl)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Get application dependencies
	app, err := bootstrap.Initialize(db, config)
	if err != nil {
		log.Fatal(err)
	}

	// Setup Router
	router := server.NewRouter(db, app.Handlers, config)

	// Start Server
	log.Printf("Server starting on port %s", config.WebsitesPort)
	if err := router.Run(":" + config.WebsitesPort); err != nil {
		log.Fatal(err)
	}
}
