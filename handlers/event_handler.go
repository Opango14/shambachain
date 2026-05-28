package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"shambachain/database"
	"shambachain/models"
	"shambachain/services"
)

// AddEventHandler handles POST /api/batches/:batchID/events for event recording
//
// This handler performs the following operations:
// 1. Extracts batch ID from URL parameter
// 2. Extracts actor ID and role from authentication context
// 3. Binds and validates AddEventRequest from request body
// 4. Validates event type and event data structure
// 5. Calls AddEvent service function
// 6. Returns success response with 200 status
// 7. Handles errors with appropriate HTTP status codes
//
// Preconditions:
//   - Request must be authenticated (actor ID and role in context)
//   - Batch ID must be provided in URL parameter
//   - Request body must contain valid AddEventRequest JSON
//   - Database connection must be available
//
// Postconditions:
//   - Returns 200 OK with success message on success
//   - Returns 400 Bad Request for validation errors
//   - Returns 401 Unauthorized if actor ID or role not in context
//   - Returns 404 Not Found if batch does not exist
//   - Returns 500 Internal Server Error for service layer errors
//
// Validates Requirements: 2.4, 2.5, 3.7, 7.5, 8.5, 11.6, 11.7, 11.8
func AddEventHandler(ctx *gin.Context) {
	// Step 1: Extract batch ID from URL parameter
	batchID := ctx.Param("batchID")
	if batchID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "batch ID is required",
		})
		return
	}

	// Step 2: Extract actor ID from authentication context
	actorID, exists := ctx.Get("user_id")
	if !exists {
		// Try alternative context keys
		actorID, exists = ctx.Get("actor_id")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "actor ID not found in authentication context",
			})
			return
		}
	}

	// Convert actor ID to string
	actorIDStr, ok := actorID.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid actor ID format",
		})
		return
	}

	// Step 2b: Extract actor role from authentication context
	actorRole, exists := ctx.Get("role")
	if !exists {
		// Try alternative context key
		actorRole, exists = ctx.Get("actor_role")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "actor role not found in authentication context",
			})
			return
		}
	}

	// Convert actor role to string
	actorRoleStr, ok := actorRole.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid actor role format",
		})
		return
	}

	// Step 3: Bind and validate AddEventRequest
	var req models.AddEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Step 4: Validate event type
	validEventTypes := map[string]bool{
		"registration":  true,
		"transport":     true,
		"quality_check": true,
		"transfer":      true,
		"sale":          true,
	}

	if !validEventTypes[req.EventType] {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid event type",
			"details": "event type must be one of: registration, transport, quality_check, transfer, sale",
		})
		return
	}

	// Step 4b: Validate event data structure (ensure it's not nil and can be processed)
	if req.EventData == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "event data cannot be empty",
		})
		return
	}

	// Step 5: Call AddEvent service function
	db := database.GetDB()
	err := services.AddEvent(db, batchID, actorIDStr, actorRoleStr, req.EventType, req.EventData)
	if err != nil {
		// Determine appropriate status code based on error type
		statusCode := http.StatusInternalServerError
		errorMessage := err.Error()

		// Check for validation errors
		if contains(errorMessage, "cannot be empty") ||
			contains(errorMessage, "cannot be nil") ||
			contains(errorMessage, "invalid") ||
			contains(errorMessage, "failed to marshal") {
			statusCode = http.StatusBadRequest
		}

		// Check for not found errors
		if contains(errorMessage, "not found") ||
			contains(errorMessage, "no blocks found") {
			statusCode = http.StatusNotFound
		}

		ctx.JSON(statusCode, gin.H{
			"error": errorMessage,
		})
		return
	}

	// Step 6: Return success response with 200 status
	ctx.JSON(http.StatusOK, gin.H{
		"message":    "event added successfully",
		"batch_id":   batchID,
		"event_type": req.EventType,
	})
}
