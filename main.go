package main

import (
	"database/sql"
	"fintech/routes"
	"fintech/utils"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Load environment variables from .env file
	err := utils.LoadEnv()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Connect to the database
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
	))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	r := gin.Default() // Create a new Gin engine

	// Disable proxy trusting by passing an empty slice
	r.SetTrustedProxies(nil)

	// Set up routes
	routes.AuthRoutes(r, db)
	routes.VideoRoutes(r, db)
	routes.UserActionRoutes(r, db)

	// Start the server
	r.Run(":8080")
}
