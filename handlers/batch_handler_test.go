package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"shambachain/models"
)

func TestRegisterBatchHandler_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Note: This test requires a real database connection
	// In a production environment, you would initialize the database before running tests
	// For now, we'll skip this test if the database is not initialized
	t.Skip("Skipping integration test - requires database initialization")

	// Setup route
	router.POST("/api/batches", func(ctx *gin.Context) {
		// Mock authentication context
		ctx.Set("farmer_id", "test-farmer-123")
		RegisterBatchHandler(ctx)
	})

	// Prepare request
	reqBody := models.RegisterBatchRequest{
		ProduceType: "potatoes",
		Quantity:    500.0,
		Unit:        "kg",
		HarvestDate: time.Now().AddDate(0, 0, -2), // 2 days ago
		Location:    "Nakuru Farm",
		FarmName:    "Green Valley Farm",
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Create request
	req, _ := http.NewRequest("POST", "/api/batches", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		t.Logf("response body: %s", w.Body.String())
	}

	// Parse response
	var response models.RegisterBatchResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	// Verify response fields
	if response.BatchID == "" {
		t.Error("expected non-empty batch ID")
	}
	if response.GenesisHash == "" {
		t.Error("expected non-empty genesis hash")
	}
	if response.QRCodeURL == "" {
		t.Error("expected non-empty QR code URL")
	}

	// Cleanup QR code if created
	if response.BatchID != "" {
		os.Remove("./qrcodes/" + response.BatchID + ".png")
	}
}

func TestRegisterBatchHandler_MissingFarmerID(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Setup route without setting farmer_id in context
	router.POST("/api/batches", RegisterBatchHandler)

	// Prepare request
	reqBody := models.RegisterBatchRequest{
		ProduceType: "potatoes",
		Quantity:    500.0,
		Unit:        "kg",
		HarvestDate: time.Now().AddDate(0, 0, -2),
		Location:    "Nakuru Farm",
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Create request
	req, _ := http.NewRequest("POST", "/api/batches", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestRegisterBatchHandler_InvalidRequestBody(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Setup route
	router.POST("/api/batches", func(ctx *gin.Context) {
		ctx.Set("farmer_id", "test-farmer-123")
		RegisterBatchHandler(ctx)
	})

	// Create request with invalid JSON
	req, _ := http.NewRequest("POST", "/api/batches", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestRegisterBatchHandler_ValidationErrors(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Note: This test requires a real database connection
	t.Skip("Skipping integration test - requires database initialization")

	// Setup route
	router.POST("/api/batches", func(ctx *gin.Context) {
		ctx.Set("farmer_id", "test-farmer-123")
		RegisterBatchHandler(ctx)
	})

	tests := []struct {
		name     string
		reqBody  models.RegisterBatchRequest
		wantCode int
	}{
		{
			name: "negative quantity",
			reqBody: models.RegisterBatchRequest{
				ProduceType: "potatoes",
				Quantity:    -10.0,
				Unit:        "kg",
				HarvestDate: time.Now().AddDate(0, 0, -2),
				Location:    "Nakuru Farm",
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "future harvest date",
			reqBody: models.RegisterBatchRequest{
				ProduceType: "potatoes",
				Quantity:    500.0,
				Unit:        "kg",
				HarvestDate: time.Now().AddDate(0, 0, 2), // 2 days in future
				Location:    "Nakuru Farm",
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "empty produce type",
			reqBody: models.RegisterBatchRequest{
				ProduceType: "",
				Quantity:    500.0,
				Unit:        "kg",
				HarvestDate: time.Now().AddDate(0, 0, -2),
				Location:    "Nakuru Farm",
			},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.reqBody)
			req, _ := http.NewRequest("POST", "/api/batches", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantCode {
				t.Errorf("expected status %d, got %d", tt.wantCode, w.Code)
				t.Logf("response body: %s", w.Body.String())
			}
		})
	}
}
