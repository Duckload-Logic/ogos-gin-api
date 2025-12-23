package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func buildDBURL() string {
	// Parse .env variables here
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	// Get port string
	dbPortStr := os.Getenv("DB_PORT")
	dbPort, err := strconv.Atoi(dbPortStr)

	// Defaul to port: 3306 if port data is not given
	if err != nil {
		log.Printf(
			"Warning: Invalid DB_PORT '%s', defaulting to 3306",
			dbPortStr,
		)
		dbPort = 3306
	}

	// Return string url
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)
}

func GetDBConnection() (*sql.DB, error) {
	// Open Database instance
	db, err := sql.Open("mysql", buildDBURL())
	if err != nil {
		log.Fatal("Failed to open Database connection:", err)
	}

	// Check ping connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping DB:", err)
	}

	// Return database instance
	return db, nil
}
