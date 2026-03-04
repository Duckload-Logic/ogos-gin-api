package database

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func GetDBConnection(dbUrl string) (*sqlx.DB, error) {
	// Open Database instance
	db, err := sqlx.Connect("mysql", dbUrl)
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
