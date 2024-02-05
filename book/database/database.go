// database/database.go
package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Database holds the reference to the database connection
type Database struct {
	Conn *sql.DB
}

// InitDB initializes the database connection
func InitDB() (*Database, error) {
	// Get database configuration from environment variables
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Default port to "3306" if not provided
	if port == "" {
		port = "3306"
	}

	// Create the database connection string
	dbSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbName)
	fmt.Println("DB Source:", dbSource)

	// Open a connection to the database
	db, err := sql.Open("mysql", dbSource)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	// Check the connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}
	fmt.Println("Connected to the database")
	return &Database{Conn: db}, nil
}
