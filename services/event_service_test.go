package services

import (
	"testing"
	"time"

	"shambachain/models"
)

// TestAddEvent_Success tests successful event addition
func TestAddEvent_Success(t *testing.T) {
	db := setupTestDB(t)

	// First, register a batch to have a genesis block
	farmerID := "farmer-test-001"
	req := models.RegisterBatchRequest{
		ProduceType: "potatoes",
		Quantity:    100.0,
		Unit:        "kg",
		HarvestDate: time.Now().AddDate(0, 0, -1),
		Location:    "Test Farm",
		FarmName:    "Test Farm Name",
	}

	resp, err := RegisterBatch(db, farmerID, req)
	if err != nil {
		t.Fatalf("failed to register batch: %v", err)
	}

	batchID := resp.BatchID

	// Now add a transport event
	transportData := map[string]interface{}{
		"from_location":     "Test Farm",
		"to_location":       "Test Market",
		"transport_id":      "TRK-001",
		"vehicle_info":      "Truck ABC 123",
		"departure_time":    time.Now(),
		"estimated_arrival": time.Now().Add(2 * time.Hour),
	}

	err = AddEvent(db, batchID, "transporter-001", "transporter", "transport", transportData)
	if err != nil {
		t.Fatalf("failed to add transport event: %v", err)
	}

	// Verify the batch status was updated
	var batch models.Batch
	if err := db.Where("id = ?", batchID).First(&batch).Error; err != nil {
		t.Fatalf("failed to fetch batch: %v", err)
	}

	if batch.Status != "in_transit" {
		t.Errorf("expected batch status 'in_transit', got '%s'", batch.Status)
	}

	// Verify we now have 2 blocks (genesis + transport)
	var blockCount int64
	db.Model(&models.Block{}).Where("batch_id = ?", batchID).Count(&blockCount)
	if blockCount != 2 {
		t.Errorf("expected 2 blocks, got %d", blockCount)
	}

	// Verify the new block has correct index and prevHash
	var latestBlock models.Block
	if err := db.Where("batch_id = ?", batchID).Order("`index` DESC").First(&latestBlock).Error; err != nil {
		t.Fatalf("failed to fetch latest block: %v", err)
	}

	if latestBlock.Index != 1 {
		t.Errorf("expected block index 1, got %d", latestBlock.Index)
	}

	if latestBlock.EventType != "transport" {
		t.Errorf("expected event type 'transport', got '%s'", latestBlock.EventType)
	}

	// Verify batch current hash matches latest block hash
	if batch.CurrentHash != latestBlock.Hash {
		t.Errorf("batch current hash does not match latest block hash")
	}
}

// TestAddEvent_NonExistentBatch tests error handling for non-existent batch
func TestAddEvent_NonExistentBatch(t *testing.T) {
	db := setupTestDB(t)

	eventData := map[string]interface{}{
		"test": "data",
	}

	err := AddEvent(db, "non-existent-batch-id", "actor-001", "farmer", "transport", eventData)
	if err == nil {
		t.Fatal("expected error for non-existent batch, got nil")
	}

	if err.Error() != "batch not found: non-existent-batch-id" {
		t.Errorf("unexpected error message: %v", err)
	}
}

// TestAddEvent_StatusTransitions tests all status transitions
func TestAddEvent_StatusTransitions(t *testing.T) {
	db := setupTestDB(t)

	// Register a batch
	farmerID := "farmer-test-002"
	req := models.RegisterBatchRequest{
		ProduceType: "maize",
		Quantity:    500.0,
		Unit:        "kg",
		HarvestDate: time.Now().AddDate(0, 0, -2),
		Location:    "Test Farm 2",
		FarmName:    "Test Farm 2 Name",
	}

	resp, err := RegisterBatch(db, farmerID, req)
	if err != nil {
		t.Fatalf("failed to register batch: %v", err)
	}

	batchID := resp.BatchID

	// Test transport event -> in_transit
	transportData := map[string]interface{}{
		"from_location": "Farm",
		"to_location":   "Market",
	}
	err = AddEvent(db, batchID, "transporter-001", "transporter", "transport", transportData)
	if err != nil {
		t.Fatalf("failed to add transport event: %v", err)
	}

	var batch models.Batch
	db.Where("id = ?", batchID).First(&batch)
	if batch.Status != "in_transit" {
		t.Errorf("expected status 'in_transit', got '%s'", batch.Status)
	}

	// Test quality_check event -> status should remain in_transit
	qualityData := map[string]interface{}{
		"inspector_id": "inspector-001",
		"grade":        "A",
		"passed":       true,
	}
	err = AddEvent(db, batchID, "inspector-001", "inspector", "quality_check", qualityData)
	if err != nil {
		t.Fatalf("failed to add quality_check event: %v", err)
	}

	db.Where("id = ?", batchID).First(&batch)
	if batch.Status != "in_transit" {
		t.Errorf("expected status to remain 'in_transit', got '%s'", batch.Status)
	}

	// Test transfer event -> delivered
	transferData := map[string]interface{}{
		"from_owner_id": farmerID,
		"to_owner_id":   "buyer-001",
		"transfer_type": "sale",
	}
	err = AddEvent(db, batchID, farmerID, "farmer", "transfer", transferData)
	if err != nil {
		t.Fatalf("failed to add transfer event: %v", err)
	}

	db.Where("id = ?", batchID).First(&batch)
	if batch.Status != "delivered" {
		t.Errorf("expected status 'delivered', got '%s'", batch.Status)
	}

	// Test sale event -> sold
	saleData := map[string]interface{}{
		"price":    50000.0,
		"currency": "KES",
	}
	err = AddEvent(db, batchID, "buyer-001", "buyer", "sale", saleData)
	if err != nil {
		t.Fatalf("failed to add sale event: %v", err)
	}

	db.Where("id = ?", batchID).First(&batch)
	if batch.Status != "sold" {
		t.Errorf("expected status 'sold', got '%s'", batch.Status)
	}
}

// TestAddEvent_SequentialIndexing tests that blocks are indexed sequentially
func TestAddEvent_SequentialIndexing(t *testing.T) {
	db := setupTestDB(t)

	// Register a batch
	farmerID := "farmer-test-003"
	req := models.RegisterBatchRequest{
		ProduceType: "vegetables",
		Quantity:    200.0,
		Unit:        "kg",
		HarvestDate: time.Now().AddDate(0, 0, -1),
		Location:    "Test Farm 3",
		FarmName:    "Test Farm 3 Name",
	}

	resp, err := RegisterBatch(db, farmerID, req)
	if err != nil {
		t.Fatalf("failed to register batch: %v", err)
	}

	batchID := resp.BatchID

	// Add multiple events
	eventData := map[string]interface{}{"test": "data"}

	for i := 0; i < 5; i++ {
		err = AddEvent(db, batchID, "actor-001", "farmer", "quality_check", eventData)
		if err != nil {
			t.Fatalf("failed to add event %d: %v", i, err)
		}
	}

	// Verify we have 6 blocks total (genesis + 5 events)
	var blocks []models.Block
	db.Where("batch_id = ?", batchID).Order("`index` ASC").Find(&blocks)

	if len(blocks) != 6 {
		t.Fatalf("expected 6 blocks, got %d", len(blocks))
	}

	// Verify sequential indexing
	for i, block := range blocks {
		if block.Index != i {
			t.Errorf("expected block index %d, got %d", i, block.Index)
		}

		// Verify hash chain
		if i > 0 {
			if block.PrevHash != blocks[i-1].Hash {
				t.Errorf("block %d prevHash does not match previous block hash", i)
			}
		}
	}
}

// TestAddEvent_InvalidInputs tests validation of input parameters
func TestAddEvent_InvalidInputs(t *testing.T) {
	db := setupTestDB(t)

	eventData := map[string]interface{}{"test": "data"}

	tests := []struct {
		name      string
		batchID   string
		actorID   string
		actorRole string
		eventType string
		eventData map[string]interface{}
		wantErr   bool
	}{
		{
			name:      "empty batchID",
			batchID:   "",
			actorID:   "actor-001",
			actorRole: "farmer",
			eventType: "transport",
			eventData: eventData,
			wantErr:   true,
		},
		{
			name:      "empty actorID",
			batchID:   "batch-001",
			actorID:   "",
			actorRole: "farmer",
			eventType: "transport",
			eventData: eventData,
			wantErr:   true,
		},
		{
			name:      "empty actorRole",
			batchID:   "batch-001",
			actorID:   "actor-001",
			actorRole: "",
			eventType: "transport",
			eventData: eventData,
			wantErr:   true,
		},
		{
			name:      "empty eventType",
			batchID:   "batch-001",
			actorID:   "actor-001",
			actorRole: "farmer",
			eventType: "",
			eventData: eventData,
			wantErr:   true,
		},
		{
			name:      "nil eventData",
			batchID:   "batch-001",
			actorID:   "actor-001",
			actorRole: "farmer",
			eventType: "transport",
			eventData: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AddEvent(db, tt.batchID, tt.actorID, tt.actorRole, tt.eventType, tt.eventData)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
