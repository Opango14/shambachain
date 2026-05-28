package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"shambachain/blockchain"
	"shambachain/models"
	"shambachain/utils"
)

// RegisterBatch registers a new produce batch with blockchain traceability.
//
// This function performs the following operations atomically:
// 1. Begins a database transaction
// 2. Generates a unique batch ID using UUID
// 3. Creates a batch record with status "registered"
// 4. Creates a genesis block (index 0, prevHash "0") with registration event data
// 5. Generates a QR code and stores the file path
// 6. Updates the batch with genesis hash and current hash
// 7. Commits the transaction or rolls back on error
//
// Preconditions:
//   - db is a valid, open database connection
//   - farmerID is non-empty string representing authenticated farmer
//   - req.ProduceType is non-empty and valid produce type
//   - req.Quantity is positive number (> 0)
//   - req.Unit is non-empty string (kg, tons, pieces)
//   - req.HarvestDate is valid date not in future
//
// Postconditions:
//   - Returns RegisterBatchResponse with unique BatchID on success
//   - Creates new Batch record in database with status "registered"
//   - Creates genesis block (index 0) in blockchain for this batch
//   - Generates QR code image and stores path in batch record
//   - response.GenesisHash matches the hash of genesis block
//   - Returns error if any database operation fails
//   - No partial state: either all operations succeed or all rollback
//
// Validates Requirements: 1.1, 1.2, 1.3, 1.6, 1.7, 1.8, 10.1, 10.2, 10.3
func RegisterBatch(db *gorm.DB, farmerID string, req models.RegisterBatchRequest) (*models.RegisterBatchResponse, error) {
	// Validate preconditions
	if db == nil {
		return nil, fmt.Errorf("database connection cannot be nil")
	}

	if farmerID == "" {
		return nil, fmt.Errorf("farmerID cannot be empty")
	}

	if req.ProduceType == "" {
		return nil, fmt.Errorf("produce_type cannot be empty")
	}

	if req.Quantity <= 0 {
		return nil, fmt.Errorf("quantity must be positive, got %.2f", req.Quantity)
	}

	if req.Unit == "" {
		return nil, fmt.Errorf("unit cannot be empty")
	}

	if req.Location == "" {
		return nil, fmt.Errorf("location cannot be empty")
	}

	// Validate harvest date is not in the future
	if req.HarvestDate.After(time.Now()) {
		return nil, fmt.Errorf("harvest_date cannot be in the future")
	}

	// Step 1: Begin database transaction
	tx := db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Ensure rollback on panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Step 2: Generate unique batch ID using UUID
	batchID := uuid.New().String()

	// Step 3: Create batch record with status "registered"
	batch := models.Batch{
		ID:          batchID,
		FarmerID:    farmerID,
		ProduceType: req.ProduceType,
		Quantity:    req.Quantity,
		Unit:        req.Unit,
		HarvestDate: req.HarvestDate,
		Location:    req.Location,
		Status:      "registered",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	// Step 4: Create genesis block (index 0, prevHash "0") with registration event data
	eventData := models.RegistrationEvent{
		ProduceType: req.ProduceType,
		Quantity:    req.Quantity,
		Unit:        req.Unit,
		HarvestDate: req.HarvestDate,
		Location:    req.Location,
		FarmName:    req.FarmName,
	}

	eventJSON, err := json.Marshal(eventData)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to marshal registration event data: %w", err)
	}

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
		tx.Rollback()
		return nil, fmt.Errorf("failed to create genesis block: %w", err)
	}

	// Step 5: Update batch with genesis hash and current hash
	batch.GenesisHash = genesisBlock.Hash
	batch.CurrentHash = genesisBlock.Hash

	// Step 6: Save batch to database
	if err := tx.Create(&batch).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create batch record: %w", err)
	}

	// Step 7: Save genesis block to database
	if err := tx.Create(&genesisBlock).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create genesis block record: %w", err)
	}

	// Step 8: Generate QR code and store path
	qrPath, err := utils.GenerateQRCode(batchID, "./qrcodes")
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Step 9: Update batch with QR code path
	batch.QRCodePath = qrPath
	if err := tx.Save(&batch).Error; err != nil {
		tx.Rollback()
		// Cleanup QR code file on failure
		os.Remove(qrPath)
		return nil, fmt.Errorf("failed to update batch with QR code path: %w", err)
	}

	// Step 10: Commit transaction
	if err := tx.Commit().Error; err != nil {
		// Cleanup QR code file on commit failure
		os.Remove(qrPath)
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Step 11: Prepare response with QR code data
	qrCodeData, err := readAndEncodeQRCode(qrPath)
	if err != nil {
		// Transaction already committed, but we can still return the response
		// without the base64 encoded data
		qrCodeData = ""
	}

	response := &models.RegisterBatchResponse{
		BatchID:     batchID,
		QRCodeURL:   fmt.Sprintf("/qrcodes/%s.png", batchID),
		QRCodeData:  qrCodeData,
		GenesisHash: genesisBlock.Hash,
	}

	return response, nil
}

// readAndEncodeQRCode reads a QR code file and returns its base64-encoded content
func readAndEncodeQRCode(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read QR code file: %w", err)
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	return encoded, nil
}
