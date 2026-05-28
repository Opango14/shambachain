package services

import (
	"fmt"

	"gorm.io/gorm"

	"shambachain/blockchain"
	"shambachain/models"
)

// GetTraceability retrieves the complete traceability information for a batch.
//
// This function performs the following operations:
// 1. Fetches the batch record by ID
// 2. Fetches all blocks for the batch ordered by index (ascending)
// 3. Validates the blockchain using ValidateChain
// 4. Verifies that the batch's current hash matches the last block's hash
// 5. Returns TraceabilityResponse with batch, blockchain, and validation results
//
// Preconditions:
//   - db is a valid, open database connection
//   - batchID is non-empty string
//
// Postconditions:
//   - Returns TraceabilityResponse with batch and full blockchain on success
//   - response.Blockchain is ordered by Index (ascending)
//   - response.ChainValid is result of ValidateChain(response.Blockchain)
//   - response.Verified is true if chain is valid and batch.CurrentHash matches last block hash
//   - Returns error if batch not found
//   - No mutations to database
//
// Validates Requirements: 6.1, 6.2, 6.3, 6.4, 6.5, 6.6
func GetTraceability(db *gorm.DB, batchID string) (*models.TraceabilityResponse, error) {
	// Validate preconditions
	if db == nil {
		return nil, fmt.Errorf("database connection cannot be nil")
	}

	if batchID == "" {
		return nil, fmt.Errorf("batchID cannot be empty")
	}

	// Step 1: Fetch batch record by ID
	var batch models.Batch
	if err := db.Where("id = ?", batchID).First(&batch).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("batch not found: %s", batchID)
		}
		return nil, fmt.Errorf("failed to fetch batch: %w", err)
	}

	// Step 2: Fetch all blocks for batch ordered by index (ascending)
	var blocks []models.Block
	if err := db.Where("batch_id = ?", batchID).Order("`index` ASC").Find(&blocks).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch blocks: %w", err)
	}

	// Step 3: Validate blockchain using ValidateChain
	chainValid := blockchain.ValidateChain(blocks)

	// Step 4: Verify batch current hash matches last block hash
	verified := false
	if chainValid && len(blocks) > 0 {
		lastBlock := blocks[len(blocks)-1]
		verified = batch.CurrentHash == lastBlock.Hash
	}

	// Step 5: Return TraceabilityResponse with validation results
	response := &models.TraceabilityResponse{
		Batch:      batch,
		Blockchain: blocks,
		Verified:   verified,
		ChainValid: chainValid,
	}

	return response, nil
}
