package main

import (
	"fmt"
	"log"
	"os"

	"shambachain/database"
	"shambachain/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	fmt.Println("Initializing database...")
	database.InitDB()

	// Create Gin router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router)

	// Start server
	fmt.Println("Starting server on :8080...")
	port := "8080"

	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	fmt.Println("Starting server on :" + port + "...")

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
