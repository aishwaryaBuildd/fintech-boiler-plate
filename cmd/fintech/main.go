package main

import (
	"fintech/pkg/vdo"
	"fintech/routes/auth"
	"fintech/routes/courses"
	"fintech/routes/folders"
	"fintech/store/mysql"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := LoadEnv()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Connect to the database
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	))

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Printf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	defer db.Close()

	r := gin.Default() // Create a new Gin engine

	// Disable proxy trusting by passing an empty slice
	r.SetTrustedProxies(nil)

	mysqlStore := mysql.NewMySQLStore(db)

	vdo := vdo.NewVideoCipherClient()

	// Set up routes
	auth.AuthRoutes(r, mysqlStore)
	courses.CourseRoutes(r, mysqlStore, vdo)
	folders.FolderRoutes(r, mysqlStore, vdo)
	// routes.VideoRoutes(r, db)
	// routes.UserActionRoutes(r, db)

	// Start the server
	r.Run(":8080")
}

// LoadEnv loads environment variables from the .env file
func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	return nil
}
