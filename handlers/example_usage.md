# Batch Registration Handler Usage

## Overview

The `RegisterBatchHandler` handles POST requests to `/api/batches` for registering new produce batches with blockchain traceability.

## Setup

To use this handler in your application, you need to:

1. Initialize the database
2. Set up authentication middleware
3. Register the route with your Gin router

## Example Integration

```go
package main

import (
	"github.com/gin-gonic/gin"
	"shambachain/auth"
	"shambachain/database"
	"shambachain/handlers"
)

func main() {
	// Initialize database
	database.InitDB()

	// Create Gin router
	router := gin.Default()

	// API routes
	api := router.Group("/api")
	{
		// Protected routes (require authentication)
		protected := api.Group("")
		protected.Use(auth.AuthMiddleware())
		{
			// Register batch endpoint
			protected.POST("/batches", handlers.RegisterBatchHandler)
		}
	}

	// Start server
	router.Run(":8080")
}
```

## API Endpoint

### POST /api/batches

Register a new produce batch with blockchain traceability.

**Authentication Required:** Yes (Bearer token)

**Request Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "produce_type": "potatoes",
  "quantity": 500.0,
  "unit": "kg",
  "harvest_date": "2024-01-15T00:00:00Z",
  "location": "Nakuru Farm, GPS: -0.3031, 36.0800",
  "farm_name": "Green Valley Farm"
}
```

**Success Response (201 Created):**
```json
{
  "batch_id": "550e8400-e29b-41d4-a716-446655440000",
  "qr_code_url": "/qrcodes/550e8400-e29b-41d4-a716-446655440000.png",
  "qr_code_data": "iVBORw0KGgoAAAANSUhEUgAA...",
  "genesis_hash": "a1b2c3d4e5f6..."
}
```

**Error Responses:**

- **400 Bad Request** - Invalid request body or validation error
```json
{
  "error": "invalid request body",
  "details": "quantity must be positive"
}
```

- **401 Unauthorized** - Missing or invalid authentication
```json
{
  "error": "farmer ID not found in authentication context"
}
```

- **500 Internal Server Error** - Server-side error
```json
{
  "error": "failed to create batch record"
}
```

## Validation Rules

The handler validates the following:

1. **produce_type**: Must be non-empty string
2. **quantity**: Must be positive number (> 0)
3. **unit**: Must be non-empty string (e.g., "kg", "tons", "pieces")
4. **harvest_date**: Must not be in the future
5. **location**: Must be non-empty string
6. **farm_name**: Optional string

## Authentication Context

The handler expects the authentication middleware to set one of the following in the Gin context:

- `farmer_id` (string) - Preferred
- `user_id` (string) - Fallback

Example middleware setup:
```go
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Extract and validate token
		token := extractToken(ctx)
		
		// Validate token and get user ID
		userID, err := validateToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}
		
		// Set farmer_id in context
		ctx.Set("farmer_id", userID)
		ctx.Next()
	}
}
```

## Testing

Run the handler tests:
```bash
go test -v ./handlers/... -run TestRegisterBatchHandler
```

Note: Integration tests that require database access are skipped by default. To run them, initialize the database before running tests.

## Requirements Validated

This handler validates the following requirements:

- **1.1**: Create unique batch record with status "registered"
- **1.4**: Validate quantity is positive
- **1.5**: Validate harvest date is not in future
- **11.1**: Validate produce_type is non-empty
- **11.2**: Validate quantity is positive
- **11.3**: Validate unit is non-empty
- **11.4**: Validate harvest_date is not in future
- **11.5**: Validate location is non-empty
