package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"shambachain/database"
	"shambachain/models"
	"shambachain/services"
)

// RegisterBatchHandler handles POST /api/batches for batch registration
//
// This handler performs the following operations:
// 1. Extracts farmer ID from authentication context
// 2. Binds and validates RegisterBatchRequest from request body
// 3. Calls RegisterBatch service function
// 4. Returns RegisterBatchResponse with 201 status on success
// 5. Handles errors with appropriate HTTP status codes
//
// Preconditions:
//   - Request must be authenticated (farmer ID in context)
//   - Request body must contain valid RegisterBatchRequest JSON
//   - Database connection must be available
//
// Postconditions:
//   - Returns 201 Created with RegisterBatchResponse on success
//   - Returns 400 Bad Request for validation errors
//   - Returns 401 Unauthorized if farmer ID not in context
//   - Returns 500 Internal Server Error for service layer errors
//
// Validates Requirements: 1.1, 1.4, 1.5, 11.1, 11.2, 11.3, 11.4, 11.5
func RegisterBatchHandler(ctx *gin.Context) {
	// Step 1: Extract farmer ID from authentication context
	// The auth middleware should set this in the context
	farmerID, exists := ctx.Get("farmer_id")
	if !exists {
		// Try to get from token or user_id as fallback
		farmerID, exists = ctx.Get("user_id")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "farmer ID not found in authentication context",
			})
			return
		}
	}

	// Convert farmer ID to string
	farmerIDStr, ok := farmerID.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid farmer ID format",
		})
		return
	}

	// Step 2: Bind and validate RegisterBatchRequest
	var req models.RegisterBatchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Step 3: Call RegisterBatch service function
	db := database.GetDB()
	response, err := services.RegisterBatch(db, farmerIDStr, req)
	if err != nil {
		// Determine appropriate status code based on error type
		statusCode := http.StatusInternalServerError
		errorMessage := err.Error()

		// Check for validation errors (these would be caught by service layer)
		if contains(errorMessage, "cannot be empty") ||
			contains(errorMessage, "must be positive") ||
			contains(errorMessage, "cannot be in the future") {
			statusCode = http.StatusBadRequest
		}

		ctx.JSON(statusCode, gin.H{
			"error": errorMessage,
		})
		return
	}

	// Step 4: Return RegisterBatchResponse with 201 status
	ctx.JSON(http.StatusCreated, response)
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
