package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"shambachain/database"
	"shambachain/routes"
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
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
