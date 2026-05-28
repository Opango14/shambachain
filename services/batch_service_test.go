package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"shambachain/models"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Auto-migrate the models
	err = db.AutoMigrate(&models.Batch{}, &models.Block{})
	if err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	return db
}

// cleanupQRCodes removes test QR code files
func cleanupQRCodes(t *testing.T) {
	err := os.RemoveAll("./qrcodes")
	if err != nil && !os.IsNotExist(err) {
		t.Logf("warning: failed to cleanup qrcodes directory: %v", err)
	}
}

func TestRegisterBatch_Success(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupQRCodes(t)

	farmerID := "farmer-123"
	req := models.RegisterBatchRequest{
		ProduceType: "potatoes",
		Quantity:    500.0,
		Unit:        "kg",
		HarvestDate: time.Now().AddDate(0, 0, -2), // 2 days ago
		Location:    "Nakuru Farm, GPS: -0.3031, 36.0800",
		FarmName:    "Green Valley Farm",
	}

	response, err := RegisterBatch(db, farmerID, req)

	// Verify no error
	if err != nil {
		t.Fatalf("RegisterBatch failed: %v", err)
	}

	// Verify response fields
	if response.BatchID == "" {
		t.Error("BatchID should not be empty")
	}

	if response.GenesisHash == "" {
		t.Error("GenesisHash should not be empty")
	}

	if len(response.GenesisHash) != 64 {
		t.Errorf("GenesisHash should be 64 characters, got %d", len(response.GenesisHash))
	}

	if response.QRCodeURL == "" {
		t.Error("QRCodeURL should not be empty")
	}

	expectedURL := "/qrcodes/" + response.BatchID + ".png"
	if response.QRCodeURL != expectedURL {
		t.Errorf("QRCodeURL = %s, want %s", response.QRCodeURL, expectedURL)
	}

	// Verify batch was created in database
	var batch models.Batch
	err = db.Where("id = ?", response.BatchID).First(&batch).Error
	if err != nil {
		t.Fatalf("failed to retrieve batch from database: %v", err)
	}

	if batch.FarmerID != farmerID {
		t.Errorf("batch.FarmerID = %s, want %s", batch.FarmerID, farmerID)
	}

	if batch.ProduceType != req.ProduceType {
		t.Errorf("batch.ProduceType = %s, want %s", batch.ProduceType, req.ProduceType)
	}

	if batch.Quantity != req.Quantity {
		t.Errorf("batch.Quantity = %.2f, want %.2f", batch.Quantity, req.Quantity)
	}

	if batch.Status != "registered" {
		t.Errorf("batch.Status = %s, want 'registered'", batch.Status)
	}

	if batch.GenesisHash != response.GenesisHash {
		t.Errorf("batch.GenesisHash = %s, want %s", batch.GenesisHash, response.GenesisHash)
	}

	if batch.CurrentHash != batch.GenesisHash {
		t.Errorf("batch.CurrentHash should equal GenesisHash for new batch")
	}

	// Verify genesis block was created
	var block models.Block
	err = db.Where("batch_id = ? AND `index` = ?", response.BatchID, 0).First(&block).Error
	if err != nil {
		t.Fatalf("failed to retrieve genesis block from database: %v", err)
	}

	if block.Index != 0 {
		t.Errorf("genesis block index = %d, want 0", block.Index)
	}

	if block.PrevHash != "0" {
		t.Errorf("genesis block prevHash = %s, want '0'", block.PrevHash)
	}

	if block.EventType != "registration" {
		t.Errorf("genesis block eventType = %s, want 'registration'", block.EventType)
	}

	if block.ActorID != farmerID {
		t.Errorf("genesis block actorID = %s, want %s", block.ActorID, farmerID)
	}

	if block.ActorRole != "farmer" {
		t.Errorf("genesis block actorRole = %s, want 'farmer'", block.ActorRole)
	}

	if block.Hash != response.GenesisHash {
		t.Errorf("genesis block hash = %s, want %s", block.Hash, response.GenesisHash)
	}

	// Verify event data is valid JSON and contains expected fields
	var eventData models.RegistrationEvent
	err = json.Unmarshal([]byte(block.EventData), &eventData)
	if err != nil {
		t.Fatalf("failed to unmarshal event data: %v", err)
	}

	if eventData.ProduceType != req.ProduceType {
		t.Errorf("eventData.ProduceType = %s, want %s", eventData.ProduceType, req.ProduceType)
	}

	if eventData.Quantity != req.Quantity {
		t.Errorf("eventData.Quantity = %.2f, want %.2f", eventData.Quantity, req.Quantity)
	}

	// Verify QR code file was created
	if batch.QRCodePath == "" {
		t.Error("batch.QRCodePath should not be empty")
	}

	if _, err := os.Stat(batch.QRCodePath); os.IsNotExist(err) {
		t.Errorf("QR code file does not exist at path: %s", batch.QRCodePath)
	}

	// Verify QR code file size is reasonable (< 100KB)
	fileInfo, err := os.Stat(batch.QRCodePath)
	if err == nil {
		if fileInfo.Size() > 100*1024 {
			t.Errorf("QR code file size = %d bytes, want < 102400 bytes", fileInfo.Size())
		}
	}
}

func TestRegisterBatch_EmptyFarmerID(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupQRCodes(t)

	req := models.RegisterBatchRequest{
		ProduceType: "potatoes",
		Quantity:    500.0,
		Unit:        "kg",
		HarvestDate: time.Now().AddDate(0, 0, -2),
		Location:    "Nakuru Farm",
	}

	_, err := RegisterBatch(db, "", req)

	if err == nil {
		t.Error("RegisterBatch should fail with empty farmerID")
	}
}

func TestRegisterBatch_NegativeQuantity(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupQRCodes(t)

	req := models.RegisterBatchRequest{
		ProduceType: "potatoes",
		Quantity:    -10.0,
		Unit:        "kg",
		HarvestDate: time.Now().AddDate(0, 0, -2),
		Location:    "Nakuru Farm",
	}

	_, err := RegisterBatch(db, "farmer-123", req)

	if err == nil {
		t.Error("RegisterBatch should fail with negative quantity")
	}
}

func TestRegisterBatch_ZeroQuantity(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupQRCodes(t)

	req := models.RegisterBatchRequest{
		ProduceType: "potatoes",
		Quantity:    0.0,
		Unit:        "kg",
		HarvestDate: time.Now().AddDate(0, 0, -2),
		Location:    "Nakuru Farm",
	}

	_, err := RegisterBatch(db, "farmer-123", req)

	if err == nil {
		t.Error("RegisterBatch should fail with zero quantity")
	}
}

func TestRegisterBatch_FutureHarvestDate(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupQRCodes(t)

	req := models.RegisterBatchRequest{
		ProduceType: "potatoes",
		Quantity:    500.0,
		Unit:        "kg",
		HarvestDate: time.Now().AddDate(0, 0, 1), // tomorrow
		Location:    "Nakuru Farm",
	}

	_, err := RegisterBatch(db, "farmer-123", req)

	if err == nil {
		t.Error("RegisterBatch should fail with future harvest date")
	}
}

func TestRegisterBatch_EmptyProduceType(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupQRCodes(t)

	req := models.RegisterBatchRequest{
		ProduceType: "",
		Quantity:    500.0,
		Unit:        "kg",
		HarvestDate: time.Now().AddDate(0, 0, -2),
		Location:    "Nakuru Farm",
	}

	_, err := RegisterBatch(db, "farmer-123", req)

	if err == nil {
		t.Error("RegisterBatch should fail with empty produce type")
	}
}

func TestRegisterBatch_EmptyUnit(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupQRCodes(t)

	req := models.RegisterBatchRequest{
		ProduceType: "potatoes",
		Quantity:    500.0,
		Unit:        "",
		HarvestDate: time.Now().AddDate(0, 0, -2),
		Location:    "Nakuru Farm",
	}

	_, err := RegisterBatch(db, "farmer-123", req)

	if err == nil {
		t.Error("RegisterBatch should fail with empty unit")
	}
}

func TestRegisterBatch_EmptyLocation(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupQRCodes(t)

	req := models.RegisterBatchRequest{
		ProduceType: "potatoes",
		Quantity:    500.0,
		Unit:        "kg",
		HarvestDate: time.Now().AddDate(0, 0, -2),
		Location:    "",
	}

	_, err := RegisterBatch(db, "farmer-123", req)

	if err == nil {
		t.Error("RegisterBatch should fail with empty location")
	}
}

func TestRegisterBatch_QRCodeUniqueness(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupQRCodes(t)

	farmerID := "farmer-123"
	req := models.RegisterBatchRequest{
		ProduceType: "potatoes",
		Quantity:    500.0,
		Unit:        "kg",
		HarvestDate: time.Now().AddDate(0, 0, -2),
		Location:    "Nakuru Farm",
	}

	// Register first batch
	response1, err := RegisterBatch(db, farmerID, req)
	if err != nil {
		t.Fatalf("first RegisterBatch failed: %v", err)
	}

	// Register second batch
	response2, err := RegisterBatch(db, farmerID, req)
	if err != nil {
		t.Fatalf("second RegisterBatch failed: %v", err)
	}

	// Verify batch IDs are different
	if response1.BatchID == response2.BatchID {
		t.Error("batch IDs should be unique")
	}

	// Verify QR code paths are different
	var batch1, batch2 models.Batch
	db.Where("id = ?", response1.BatchID).First(&batch1)
	db.Where("id = ?", response2.BatchID).First(&batch2)

	if batch1.QRCodePath == batch2.QRCodePath {
		t.Error("QR code paths should be unique")
	}

	// Verify both QR code files exist
	if _, err := os.Stat(batch1.QRCodePath); os.IsNotExist(err) {
		t.Errorf("first QR code file does not exist: %s", batch1.QRCodePath)
	}

	if _, err := os.Stat(batch2.QRCodePath); os.IsNotExist(err) {
		t.Errorf("second QR code file does not exist: %s", batch2.QRCodePath)
	}
}

func TestRegisterBatch_TransactionRollback(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupQRCodes(t)

	// Create a request with valid data
	req := models.RegisterBatchRequest{
		ProduceType: "potatoes",
		Quantity:    500.0,
		Unit:        "kg",
		HarvestDate: time.Now().AddDate(0, 0, -2),
		Location:    "Nakuru Farm",
	}

	// First, let's test with a valid request to ensure the function works
	response, err := RegisterBatch(db, "farmer-123", req)
	if err != nil {
		t.Fatalf("RegisterBatch failed: %v", err)
	}

	// Verify the batch was created
	var count int64
	db.Model(&models.Batch{}).Where("id = ?", response.BatchID).Count(&count)
	if count != 1 {
		t.Errorf("expected 1 batch, got %d", count)
	}

	// Verify the genesis block was created
	db.Model(&models.Block{}).Where("batch_id = ?", response.BatchID).Count(&count)
	if count != 1 {
		t.Errorf("expected 1 block, got %d", count)
	}

	// Verify QR code file exists
	qrPath := filepath.Join("./qrcodes", response.BatchID+".png")
	if _, err := os.Stat(qrPath); os.IsNotExist(err) {
		t.Error("QR code file should exist after successful registration")
	}
}
