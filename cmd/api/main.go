package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

type application struct {
	port int
	db   *sql.DB
}

func buildDBURL() string {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	dbPortStr := os.Getenv("DB_PORT")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Printf("Warning: Invalid DB_PORT '%s', defaulting to 3306", dbPortStr)
		dbPort = 3306
	}

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)
}

func (app *application) serve() error {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Connected to database via Application Struct!",
		})
	})

	serverAddr := fmt.Sprintf(":%d", app.port)
	log.Printf("Starting server on port %s", serverAddr)
	return r.Run(serverAddr)
}

func main() {
	db, err := sql.Open("mysql", buildDBURL())
	if err != nil {
		log.Fatal("Failed to open Database connection:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping DB:", err)
	}
	fmt.Println("Connected to the Database Successfully!")

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
