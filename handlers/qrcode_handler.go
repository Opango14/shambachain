package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// ServeQRCodeHandler handles GET /qrcodes/:filename for serving QR code images
func ServeQRCodeHandler(ctx *gin.Context) {
	// Extract filename from URL parameter
	filename := ctx.Param("filename")
	if filename == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "filename is required",
		})
		return
	}

	// Construct file path
	filePath := filepath.Join("./qrcodes", filename)

	// Serve the file
	ctx.File(filePath)
}
