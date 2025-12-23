package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/olazo-johnalbert/duckload-api/internal/bootstrap"
	"github.com/olazo-johnalbert/duckload-api/internal/database"
)

func main() {
	// Get database connection
	db, err := database.GetDBConnection()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Connected to the Database Successfully!")

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Migration error: %v", err)
		return
	}

	log.Println("Database migrations completed successfully")

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
