package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"shambachain/database"
	"shambachain/services"
)

// GetTraceabilityHandler handles GET /api/trace/:batchID for traceability retrieval
func GetTraceabilityHandler(ctx *gin.Context) {
	// Extract batch ID from URL parameter
	batchID := ctx.Param("batchID")
	if batchID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "batch ID is required",
		})
		return
	}

	// Call GetTraceability service function
	db := database.GetDB()
	response, err := services.GetTraceability(db, batchID)
	if err != nil {
		// Determine appropriate status code
		statusCode := http.StatusInternalServerError
		errorMessage := err.Error()

		if contains(errorMessage, "not found") {
			statusCode = http.StatusNotFound
		}

		ctx.JSON(statusCode, gin.H{
			"error": errorMessage,
		})
		return
	}

	// Return traceability response
	ctx.JSON(http.StatusOK, response)
}
