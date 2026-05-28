package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"shambachain/models"
)

func TestAddEventHandler_MissingBatchID(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Setup route - when batch ID is empty, Gin will return 404
	// This test verifies the route behavior
	router.POST("/api/batches/:batchID/events", func(ctx *gin.Context) {
		ctx.Set("user_id", "test-actor-123")
		ctx.Set("role", "transporter")
		AddEventHandler(ctx)
	})

	// Prepare request
	reqBody := models.AddEventRequest{
		EventType: "transport",
		EventData: map[string]interface{}{
			"from_location": "Farm A",
			"to_location":   "Market B",
		},
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Create request with empty batch ID (will result in 404 from router)
	req, _ := http.NewRequest("POST", "/api/batches//events", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions - Handler returns 400 when batch ID is empty
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAddEventHandler_MissingActorID(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Setup route without setting actor_id in context
	router.POST("/api/batches/:batchID/events", AddEventHandler)

	// Prepare request
	reqBody := models.AddEventRequest{
		EventType: "transport",
		EventData: map[string]interface{}{
			"from_location": "Farm A",
			"to_location":   "Market B",
		},
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Create request
	req, _ := http.NewRequest("POST", "/api/batches/test-batch-123/events", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAddEventHandler_MissingActorRole(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Setup route with actor_id but without role
	router.POST("/api/batches/:batchID/events", func(ctx *gin.Context) {
		ctx.Set("user_id", "test-actor-123")
		AddEventHandler(ctx)
	})

	// Prepare request
	reqBody := models.AddEventRequest{
		EventType: "transport",
		EventData: map[string]interface{}{
			"from_location": "Farm A",
			"to_location":   "Market B",
		},
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Create request
	req, _ := http.NewRequest("POST", "/api/batches/test-batch-123/events", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAddEventHandler_InvalidRequestBody(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Setup route
	router.POST("/api/batches/:batchID/events", func(ctx *gin.Context) {
		ctx.Set("user_id", "test-actor-123")
		ctx.Set("role", "transporter")
		AddEventHandler(ctx)
	})

	// Create request with invalid JSON
	req, _ := http.NewRequest("POST", "/api/batches/test-batch-123/events", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAddEventHandler_InvalidEventType(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Setup route
	router.POST("/api/batches/:batchID/events", func(ctx *gin.Context) {
		ctx.Set("user_id", "test-actor-123")
		ctx.Set("role", "transporter")
		AddEventHandler(ctx)
	})

	// Prepare request with invalid event type
	reqBody := models.AddEventRequest{
		EventType: "invalid_event",
		EventData: map[string]interface{}{
			"from_location": "Farm A",
			"to_location":   "Market B",
		},
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Create request
	req, _ := http.NewRequest("POST", "/api/batches/test-batch-123/events", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Verify error message
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	if response["error"] != "invalid event type" {
		t.Errorf("expected error message about invalid event type, got: %v", response["error"])
	}
}

func TestAddEventHandler_EmptyEventData(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Setup route
	router.POST("/api/batches/:batchID/events", func(ctx *gin.Context) {
		ctx.Set("user_id", "test-actor-123")
		ctx.Set("role", "transporter")
		AddEventHandler(ctx)
	})

	// Prepare request with nil event data
	reqBody := models.AddEventRequest{
		EventType: "transport",
		EventData: nil,
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Create request
	req, _ := http.NewRequest("POST", "/api/batches/test-batch-123/events", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAddEventHandler_ValidEventTypes(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Note: This test requires a real database connection
	t.Skip("Skipping integration test - requires database initialization")

	validEventTypes := []string{"registration", "transport", "quality_check", "transfer", "sale"}

	for _, eventType := range validEventTypes {
		t.Run(eventType, func(t *testing.T) {
			router := gin.New()

			// Setup route
			router.POST("/api/batches/:batchID/events", func(ctx *gin.Context) {
				ctx.Set("user_id", "test-actor-123")
				ctx.Set("role", "farmer")
				AddEventHandler(ctx)
			})

			// Prepare request
			reqBody := models.AddEventRequest{
				EventType: eventType,
				EventData: map[string]interface{}{
					"test_field": "test_value",
				},
			}
			jsonBody, _ := json.Marshal(reqBody)

			// Create request
			req, _ := http.NewRequest("POST", "/api/batches/test-batch-123/events", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Record response
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// For this test, we expect either 200 (success) or 404 (batch not found)
			// but NOT 400 (bad request for invalid event type)
			if w.Code == http.StatusBadRequest {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				if response["error"] == "invalid event type" {
					t.Errorf("event type %s should be valid but was rejected", eventType)
				}
			}
		})
	}
}

func TestAddEventHandler_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Note: This test requires a real database connection and existing batch
	t.Skip("Skipping integration test - requires database initialization")

	// Setup route
	router.POST("/api/batches/:batchID/events", func(ctx *gin.Context) {
		ctx.Set("user_id", "test-transporter-123")
		ctx.Set("role", "transporter")
		AddEventHandler(ctx)
	})

	// Prepare request
	reqBody := models.AddEventRequest{
		EventType: "transport",
		EventData: map[string]interface{}{
			"from_location":     "Nakuru Farm",
			"to_location":       "Nairobi Market",
			"transport_id":      "TRK-001",
			"vehicle_info":      "Truck KBZ 123A",
			"departure_time":    "2024-01-15T08:00:00Z",
			"estimated_arrival": "2024-01-15T12:00:00Z",
		},
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Create request (use a real batch ID in actual test)
	req, _ := http.NewRequest("POST", "/api/batches/existing-batch-id/events", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		t.Logf("response body: %s", w.Body.String())
	}

	// Parse response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	// Verify response fields
	if response["message"] != "event added successfully" {
		t.Errorf("expected success message, got: %v", response["message"])
	}
	if response["batch_id"] == "" {
		t.Error("expected non-empty batch ID in response")
	}
	if response["event_type"] != "transport" {
		t.Errorf("expected event_type 'transport', got: %v", response["event_type"])
	}
}

func TestAddEventHandler_AlternativeContextKeys(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Note: This test requires a real database connection
	t.Skip("Skipping integration test - requires database initialization")

	// Setup route with alternative context keys (actor_id and actor_role)
	router.POST("/api/batches/:batchID/events", func(ctx *gin.Context) {
		ctx.Set("actor_id", "test-actor-456")
		ctx.Set("actor_role", "inspector")
		AddEventHandler(ctx)
	})

	// Prepare request
	reqBody := models.AddEventRequest{
		EventType: "quality_check",
		EventData: map[string]interface{}{
			"inspector_id":   "test-actor-456",
			"inspector_name": "John Doe",
			"grade":          "A",
			"passed":         true,
		},
	}
	jsonBody, _ := json.Marshal(reqBody)

	// Create request
	req, _ := http.NewRequest("POST", "/api/batches/existing-batch-id/events", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should not return unauthorized error
	if w.Code == http.StatusUnauthorized {
		t.Errorf("handler should accept alternative context keys, got unauthorized")
		t.Logf("response body: %s", w.Body.String())
	}
}
