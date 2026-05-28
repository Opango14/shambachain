package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateQRCode_Success(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	batchID := "test-batch-123"
	filePath, err := GenerateQRCode(batchID, tempDir)

	if err != nil {
		t.Fatalf("GenerateQRCode failed: %v", err)
	}

	// Verify file path is returned
	if filePath == "" {
		t.Error("Expected non-empty file path")
	}

	// Verify file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("QR code file does not exist at path: %s", filePath)
	}

	// Verify file has correct name
	expectedFilename := batchID + ".png"
	if !strings.HasSuffix(filePath, expectedFilename) {
		t.Errorf("Expected filename to end with %s, got %s", expectedFilename, filePath)
	}

	// Verify file size is reasonable (< 100KB)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	const maxFileSize = 100 * 1024 // 100KB
	if fileInfo.Size() > maxFileSize {
		t.Errorf("File size %d exceeds maximum %d bytes", fileInfo.Size(), maxFileSize)
	}

	if fileInfo.Size() == 0 {
		t.Error("File size is 0, expected non-empty file")
	}
}

func TestGenerateQRCode_EmptyBatchID(t *testing.T) {
	tempDir := t.TempDir()

	_, err := GenerateQRCode("", tempDir)

	if err == nil {
		t.Error("Expected error for empty batchID, got nil")
	}

	if !strings.Contains(err.Error(), "batchID cannot be empty") {
		t.Errorf("Expected error message about empty batchID, got: %v", err)
	}
}

func TestGenerateQRCode_EmptyOutputPath(t *testing.T) {
	_, err := GenerateQRCode("test-batch-123", "")

	if err == nil {
		t.Error("Expected error for empty outputPath, got nil")
	}

	if !strings.Contains(err.Error(), "outputPath cannot be empty") {
		t.Errorf("Expected error message about empty outputPath, got: %v", err)
	}
}

func TestGenerateQRCode_CreatesDirectory(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Use a subdirectory that doesn't exist yet
	outputPath := filepath.Join(tempDir, "qrcodes", "nested")
	batchID := "test-batch-456"

	filePath, err := GenerateQRCode(batchID, outputPath)

	if err != nil {
		t.Fatalf("GenerateQRCode failed: %v", err)
	}

	// Verify the nested directory was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Output directory was not created: %s", outputPath)
	}

	// Verify file exists in the nested directory
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("QR code file does not exist at path: %s", filePath)
	}
}

func TestGenerateQRCode_UniqueFiles(t *testing.T) {
	tempDir := t.TempDir()

	// Generate QR codes for different batch IDs
	batchID1 := "batch-001"
	batchID2 := "batch-002"

	filePath1, err1 := GenerateQRCode(batchID1, tempDir)
	if err1 != nil {
		t.Fatalf("Failed to generate first QR code: %v", err1)
	}

	filePath2, err2 := GenerateQRCode(batchID2, tempDir)
	if err2 != nil {
		t.Fatalf("Failed to generate second QR code: %v", err2)
	}

	// Verify file paths are different
	if filePath1 == filePath2 {
		t.Error("Expected unique file paths for different batch IDs")
	}

	// Verify both files exist
	if _, err := os.Stat(filePath1); os.IsNotExist(err) {
		t.Errorf("First QR code file does not exist: %s", filePath1)
	}

	if _, err := os.Stat(filePath2); os.IsNotExist(err) {
		t.Errorf("Second QR code file does not exist: %s", filePath2)
	}
}

func TestGenerateQRCode_InvalidOutputPath(t *testing.T) {
	// Try to write to a path that should fail (e.g., a file instead of directory)
	tempDir := t.TempDir()
	
	// Create a file
	filePath := filepath.Join(tempDir, "notadir")
	if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Try to use the file as output directory
	_, err := GenerateQRCode("test-batch", filePath)

	if err == nil {
		t.Error("Expected error when output path is a file, got nil")
	}
}
