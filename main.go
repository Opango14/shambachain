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

	// Serve static JS files
	router.Static("/scripts", "./front-end/scripts")

	// Silence favicon 404s
	router.GET("/favicon.ico", func(c *gin.Context) { c.Status(204) })

	// Serve HTML templates
	router.LoadHTMLGlob("front-end/templates/*")

	// Homepage route
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	// Serve all HTML pages — only handle .html requests, ignore favicon/devtools/etc.
	router.GET("/:page", func(c *gin.Context) {
		page := c.Param("page")
		// Only serve files that end in .html
		if len(page) < 5 || page[len(page)-5:] != ".html" {
			c.Status(404)
			return
		}
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
