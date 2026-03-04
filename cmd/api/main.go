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
// @host            localhost:8080
// @BasePath        /api/v1
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

	db, err := database.GetDBConnection(dbUrl)
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
