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
		// Protected routes (require authentication)
		protected := api.Group("")
		protected.Use(auth.AuthMiddleware())
		{
			// Batch registration
			protected.POST("/batches", handlers.RegisterBatchHandler)

			// Event recording
			protected.POST("/batches/:batchID/events", handlers.AddEventHandler)
		}

		// Public routes (no authentication required)
		// Traceability retrieval - public so buyers can scan QR codes
		api.GET("/trace/:batchID", handlers.GetTraceabilityHandler)
	}

	// QR code serving (public)
	router.GET("/qrcodes/:filename", handlers.ServeQRCodeHandler)
}
