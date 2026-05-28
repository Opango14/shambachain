package services

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"shambachain/blockchain"
	"shambachain/models"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDBForTraceability(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(&models.Batch{}, &models.Block{}); err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	return db
}

// createTestBatchWithBlocks creates a test batch with a blockchain for testing
func createTestBatchWithBlocks(t *testing.T, db *gorm.DB, numBlocks int) (string, []models.Block) {
	batchID := uuid.New().String()
	farmerID := "test-farmer-123"

	// Create genesis block
	eventData := models.RegistrationEvent{
		ProduceType: "potatoes",
		Quantity:    100.0,
		Unit:        "kg",
		HarvestDate: time.Now().AddDate(0, 0, -1),
		Location:    "Test Farm",
		FarmName:    "Test Farm Name",
	}
	eventJSON, _ := json.Marshal(eventData)

	genesisBlock, err := blockchain.CreateBlock(
		batchID,
		0,
		"registration",
		string(eventJSON),
		farmerID,
		"farmer",
		"0",
	)
	if err != nil {
		t.Fatalf("failed to create genesis block: %v", err)
	}

	// Create batch
	batch := models.Batch{
		ID:          batchID,
		FarmerID:    farmerID,
		ProduceType: "potatoes",
		Quantity:    100.0,
		Unit:        "kg",
		HarvestDate: time.Now().AddDate(0, 0, -1),
		Location:    "Test Farm",
		Status:      "registered",
		GenesisHash: genesisBlock.Hash,
		CurrentHash: genesisBlock.Hash,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := db.Create(&batch).Error; err != nil {
		t.Fatalf("failed to create test batch: %v", err)
	}

	if err := db.Create(&genesisBlock).Error; err != nil {
		t.Fatalf("failed to create genesis block: %v", err)
	}

	blocks := []models.Block{genesisBlock}

	// Create additional blocks if requested
	for i := 1; i < numBlocks; i++ {
		prevBlock := blocks[i-1]

		transportData := models.TransportEvent{
			FromLocation:     "Test Farm",
			ToLocation:       "Test Market",
			TransportID:      "TRK-001",
			VehicleInfo:      "Truck ABC 123",
			DepartureTime:    time.Now(),
			EstimatedArrival: time.Now().Add(2 * time.Hour),
		}
		transportJSON, _ := json.Marshal(transportData)

		newBlock, err := blockchain.CreateBlock(
			batchID,
			i,
			"transport",
			string(transportJSON),
			"transporter-123",
			"transporter",
			prevBlock.Hash,
		)
		if err != nil {
			t.Fatalf("failed to create block %d: %v", i, err)
		}

		if err := db.Create(&newBlock).Error; err != nil {
			t.Fatalf("failed to save block %d: %v", i, err)
		}

		blocks = append(blocks, newBlock)

		// Update batch current hash
		batch.CurrentHash = newBlock.Hash
		if err := db.Save(&batch).Error; err != nil {
			t.Fatalf("failed to update batch: %v", err)
		}
	}

	return batchID, blocks
}

func TestGetTraceability_ValidBatch(t *testing.T) {
	db := setupTestDBForTraceability(t)
	batchID, expectedBlocks := createTestBatchWithBlocks(t, db, 3)

	response, err := GetTraceability(db, batchID)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if response == nil {
		t.Fatal("expected response, got nil")
	}

	// Verify batch data
	if response.Batch.ID != batchID {
		t.Errorf("expected batch ID %s, got %s", batchID, response.Batch.ID)
	}

	// Verify blockchain is returned
	if len(response.Blockchain) != 3 {
		t.Errorf("expected 3 blocks, got %d", len(response.Blockchain))
	}

	// Verify blocks are ordered by index
	for i, block := range response.Blockchain {
		if block.Index != i {
			t.Errorf("expected block index %d, got %d", i, block.Index)
		}
	}

	// Verify chain is valid
	if !response.ChainValid {
		t.Error("expected chain to be valid")
	}

	// Verify batch is verified (current hash matches last block)
	if !response.Verified {
		t.Error("expected batch to be verified")
	}

	// Verify last block hash matches batch current hash
	lastBlock := response.Blockchain[len(response.Blockchain)-1]
	if response.Batch.CurrentHash != lastBlock.Hash {
		t.Errorf("expected batch current hash %s to match last block hash %s",
			response.Batch.CurrentHash, lastBlock.Hash)
	}

	// Verify blocks match expected blocks
	for i, block := range response.Blockchain {
		if block.Hash != expectedBlocks[i].Hash {
			t.Errorf("block %d: expected hash %s, got %s", i, expectedBlocks[i].Hash, block.Hash)
		}
	}
}

func TestGetTraceability_NonExistentBatch(t *testing.T) {
	db := setupTestDBForTraceability(t)
	nonExistentID := "non-existent-batch-id"

	response, err := GetTraceability(db, nonExistentID)

	if err == nil {
		t.Fatal("expected error for non-existent batch, got nil")
	}

	if response != nil {
		t.Errorf("expected nil response, got: %v", response)
	}

	// Verify error message mentions batch not found
	expectedMsg := "batch not found"
	if err.Error()[:len(expectedMsg)] != expectedMsg {
		t.Errorf("expected error message to start with '%s', got: %s", expectedMsg, err.Error())
	}
}

func TestGetTraceability_EmptyBatchID(t *testing.T) {
	db := setupTestDBForTraceability(t)

	response, err := GetTraceability(db, "")

	if err == nil {
		t.Fatal("expected error for empty batch ID, got nil")
	}

	if response != nil {
		t.Errorf("expected nil response, got: %v", response)
	}

	expectedMsg := "batchID cannot be empty"
	if err.Error() != expectedMsg {
		t.Errorf("expected error message '%s', got: %s", expectedMsg, err.Error())
	}
}

func TestGetTraceability_NilDatabase(t *testing.T) {
	response, err := GetTraceability(nil, "some-batch-id")

	if err == nil {
		t.Fatal("expected error for nil database, got nil")
	}

	if response != nil {
		t.Errorf("expected nil response, got: %v", response)
	}

	expectedMsg := "database connection cannot be nil"
	if err.Error() != expectedMsg {
		t.Errorf("expected error message '%s', got: %s", expectedMsg, err.Error())
	}
}

func TestGetTraceability_TamperedChain(t *testing.T) {
	db := setupTestDBForTraceability(t)
	batchID, _ := createTestBatchWithBlocks(t, db, 3)

	// Tamper with a block's hash
	var block models.Block
	if err := db.Where("batch_id = ? AND `index` = ?", batchID, 1).First(&block).Error; err != nil {
		t.Fatalf("failed to fetch block for tampering: %v", err)
	}

	block.Hash = "tampered_hash_1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	if err := db.Save(&block).Error; err != nil {
		t.Fatalf("failed to save tampered block: %v", err)
	}

	response, err := GetTraceability(db, batchID)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Chain should be invalid due to tampering
	if response.ChainValid {
		t.Error("expected chain to be invalid after tampering")
	}

	// Batch should not be verified
	if response.Verified {
		t.Error("expected batch to not be verified after tampering")
	}
}

func TestGetTraceability_MismatchedCurrentHash(t *testing.T) {
	db := setupTestDBForTraceability(t)
	batchID, _ := createTestBatchWithBlocks(t, db, 2)

	// Update batch current hash to mismatch last block
	var batch models.Batch
	if err := db.Where("id = ?", batchID).First(&batch).Error; err != nil {
		t.Fatalf("failed to fetch batch: %v", err)
	}

	batch.CurrentHash = "mismatched_hash_1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	if err := db.Save(&batch).Error; err != nil {
		t.Fatalf("failed to save batch: %v", err)
	}

	response, err := GetTraceability(db, batchID)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Chain should still be valid (blocks are intact)
	if !response.ChainValid {
		t.Error("expected chain to be valid")
	}

	// But batch should not be verified (current hash mismatch)
	if response.Verified {
		t.Error("expected batch to not be verified with mismatched current hash")
	}
}

func TestGetTraceability_SingleBlock(t *testing.T) {
	db := setupTestDBForTraceability(t)
	batchID, _ := createTestBatchWithBlocks(t, db, 1)

	response, err := GetTraceability(db, batchID)

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify single block (genesis)
	if len(response.Blockchain) != 1 {
		t.Errorf("expected 1 block, got %d", len(response.Blockchain))
	}

	if response.Blockchain[0].Index != 0 {
		t.Errorf("expected genesis block with index 0, got %d", response.Blockchain[0].Index)
	}

	if !response.ChainValid {
		t.Error("expected chain to be valid")
	}

	if !response.Verified {
		t.Error("expected batch to be verified")
	}
}

func TestGetTraceability_NoMutations(t *testing.T) {
	db := setupTestDBForTraceability(t)
	batchID, _ := createTestBatchWithBlocks(t, db, 2)

	// Fetch original batch and blocks
	var originalBatch models.Batch
	db.Where("id = ?", batchID).First(&originalBatch)

	var originalBlocks []models.Block
	db.Where("batch_id = ?", batchID).Order("`index` ASC").Find(&originalBlocks)

	// Call GetTraceability
	_, err := GetTraceability(db, batchID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify no mutations occurred
	var afterBatch models.Batch
	db.Where("id = ?", batchID).First(&afterBatch)

	if originalBatch.CurrentHash != afterBatch.CurrentHash {
		t.Error("batch current hash was modified")
	}

	if originalBatch.UpdatedAt != afterBatch.UpdatedAt {
		t.Error("batch updated_at was modified")
	}

	var afterBlocks []models.Block
	db.Where("batch_id = ?", batchID).Order("`index` ASC").Find(&afterBlocks)

	if len(originalBlocks) != len(afterBlocks) {
		t.Error("number of blocks changed")
	}

	for i := range originalBlocks {
		if originalBlocks[i].Hash != afterBlocks[i].Hash {
			t.Errorf("block %d hash was modified", i)
		}
	}
}
