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

	// Serve static CSS files
	router.Static("/styles", "./front-end/styles")

	// Serve HTML templates
	router.LoadHTMLGlob("front-end/templates/*")

	// Homepage route
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	// Serve all HTML pages
	router.GET("/:page", func(c *gin.Context) {
		page := c.Param("page")
		c.HTML(200, page, nil)
	})

	// Setup API routes
	routes.SetupRoutes(router)

	// Start server
	port := "8080"

	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	fmt.Println("Starting server on :" + port + "...")

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}