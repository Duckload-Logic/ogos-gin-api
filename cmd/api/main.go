package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	"github.com/olazo-johnalbert/duckload-api/internal/database"
)

type application struct {
	port int
	db   *sql.DB
}

func (app *application) serve() error {
	// I used a separate struct function to avoid future bloating -Albert
	router := app.routes()

	serverAddr := fmt.Sprintf(":%d", app.port)
	log.Printf("Starting server on port %s", serverAddr)

	return router.Run(serverAddr)
}

func main() {
	db, err := database.GetDBConnection()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Connected to the Database Successfully!")

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Migration error: %v", err)
		return
	}

	log.Println("Database migrations completed successfully")

	portStr := os.Getenv("API_PORT")
	if portStr == "" {
		portStr = "8080"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal("Invalid API_PORT:", err)
	}

	app := &application{
		port: port,
		db:   db,
	}

	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}
