// main.go
package main

import (
	"book/database"
	"book/handlers"
	"book/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Set up the database connection
	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Conn.Close() // Close the database connection when the application exits

	// Initialize the database connection in handlers
	handlers.InitializeDB(db.Conn)

	// Set up the Gin router
	router := gin.Default()

	// Use the routes defined in routes.go
	routes.SetupRoutes(router, db)

	// Run the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err = router.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}


