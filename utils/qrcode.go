package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/skip2/go-qrcode"
)

// GenerateQRCode generates a QR code for a batch ID and saves it as a PNG file.
//
// Preconditions:
//   - batchID is non-empty string
//   - outputPath is valid writable directory path
//   - System has write permissions to outputPath
//
// Postconditions:
//   - Returns file path to generated QR code image on success
//   - QR code encodes URL: https://domain.com/trace/{batchID}
//   - Image is saved as PNG format
//   - File size is reasonable (< 100KB)
//   - Returns error if file write fails or QR generation fails
//   - No partial files left on error
func GenerateQRCode(batchID string, outputPath string) (string, error) {
	// Validate inputs
	if batchID == "" {
		return "", fmt.Errorf("batchID cannot be empty")
	}

	if outputPath == "" {
		return "", fmt.Errorf("outputPath cannot be empty")
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Construct the traceability URL
	url := fmt.Sprintf("https://domain.com/trace/%s", batchID)

	// Generate the file path
	filename := fmt.Sprintf("%s.png", batchID)
	filePath := filepath.Join(outputPath, filename)

	// Generate QR code with medium recovery level and size 256x256
	// This should keep the file size well under 100KB
	err := qrcode.WriteFile(url, qrcode.Medium, 256, filePath)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Verify the file was created and check its size
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		// Clean up the file if stat fails
		os.Remove(filePath)
		return "", fmt.Errorf("failed to verify QR code file: %w", err)
	}

	// Check file size (should be < 100KB = 102400 bytes)
	const maxFileSize = 100 * 1024 // 100KB in bytes
	if fileInfo.Size() > maxFileSize {
		// Clean up the oversized file
		os.Remove(filePath)
		return "", fmt.Errorf("QR code file size (%d bytes) exceeds maximum allowed size (%d bytes)", fileInfo.Size(), maxFileSize)
	}

	return filePath, nil
}
