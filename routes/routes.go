package routes

import (
	"github.com/gin-gonic/gin"
	"shambachain/auth"
	"shambachain/handlers"
)

// SetupRoutes configures all API routes for the application
func SetupRoutes(router *gin.Engine) {
	// API routes
	api := router.Group("/api")
	{
		// Auth routes
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", handlers.RegisterHandler)
			authGroup.POST("/login", handlers.LoginHandler)
		}

		// Protected routes (require authentication)
		protected := api.Group("")
		protected.Use(auth.AuthMiddleware())
		{
			// User profile
			protected.GET("/user/profile", handlers.GetProfileHandler)

			// Batch registration
			protected.POST("/batches", handlers.RegisterBatchHandler)

			// Event recording
			protected.POST("/batches/:batchID/events", handlers.AddEventHandler)
		}

		// Public routes (no authentication required)
		// Marketplace - public so buyers can see available products
		api.GET("/marketplace", handlers.GetMarketplaceHandler)

		// Traceability retrieval - public so buyers can scan QR codes
		api.GET("/trace/:batchID", handlers.GetTraceabilityHandler)
	}

	// QR code serving (public)
	router.GET("/qrcodes/:filename", handlers.ServeQRCodeHandler)
}
