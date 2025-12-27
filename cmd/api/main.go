package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/olazo-johnalbert/duckload-api/internal/bootstrap"
	"github.com/olazo-johnalbert/duckload-api/internal/database"
)

// @title           DuckLoad API
// @version         1.0
// @host            localhost:8080
// @BasePath        /api/v1

func main() {
	// Get database connection
	db, err := database.GetDBConnection()
	if err != nil {
		fmt.Println(err.Error())
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
