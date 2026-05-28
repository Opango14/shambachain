package services

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"

	"shambachain/blockchain"
	"shambachain/models"
)

// AddEvent adds a new event to a batch's blockchain.
//
// This function performs the following operations atomically:
// 1. Begins a database transaction with row locking
// 2. Fetches the batch record (with lock to prevent concurrent modifications)
// 3. Fetches the latest block for the batch
// 4. Creates a new block with incremented index
// 5. Sets prevHash to the latest block's hash
// 6. Saves the new block to the database
// 7. Updates the batch's current hash and status based on event type
// 8. Commits the transaction or rolls back on error
//
// Preconditions:
//   - db is valid, open database connection
//   - batchID exists in database
//   - actorID is non-empty string
//   - actorRole is valid role (farmer, transporter, inspector, buyer)
//   - eventType is valid event type (registration, transport, quality_check, transfer, sale)
//   - eventData is valid map that can be marshaled to JSON
//
// Postconditions:
//   - Creates new block appended to batch's blockchain
//   - New block's index = previous max index + 1
//   - New block's PrevHash = previous block's Hash
//   - Updates batch's CurrentHash to new block's Hash
//   - Updates batch's Status based on event type:
//   - "transport" -> "in_transit"
//   - "transfer" -> "delivered"
//   - "sale" -> "sold"
//   - Returns nil on success, error on failure
//   - Transaction is atomic: either all updates succeed or all rollback
//
// Validates Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6, 3.7, 3.8, 10.1, 10.2, 10.3, 10.5
func AddEvent(db *gorm.DB, batchID string, actorID string, actorRole string,
	eventType string, eventData map[string]interface{}) error {

	// Validate preconditions
	if db == nil {
		return fmt.Errorf("database connection cannot be nil")
	}

	if batchID == "" {
		return fmt.Errorf("batchID cannot be empty")
	}

	if actorID == "" {
		return fmt.Errorf("actorID cannot be empty")
	}

	if actorRole == "" {
		return fmt.Errorf("actorRole cannot be empty")
	}

	if eventType == "" {
		return fmt.Errorf("eventType cannot be empty")
	}

	if eventData == nil {
		return fmt.Errorf("eventData cannot be nil")
	}

	// Step 1: Begin database transaction
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Ensure rollback on panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Step 2: Fetch batch with row lock to prevent concurrent modifications
	// Using Clauses(clause.Locking{Strength: "UPDATE"}) for row-level locking
	var batch models.Batch
	if err := tx.Set("gorm:query_option", "FOR UPDATE").
		Where("id = ?", batchID).
		First(&batch).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("batch not found: %s", batchID)
		}
		return fmt.Errorf("failed to fetch batch: %w", err)
	}

	// Step 3: Get latest block for this batch
	var latestBlock models.Block
	if err := tx.Where("batch_id = ?", batchID).
		Order("`index` DESC").
		First(&latestBlock).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("no blocks found for batch: %s (batch may be corrupted)", batchID)
		}
		return fmt.Errorf("failed to fetch latest block: %w", err)
	}

	// Step 4: Marshal event data to JSON
	eventJSON, err := json.Marshal(eventData)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to marshal event data to JSON: %w", err)
	}

	// Step 5: Create new block with incremented index
	newIndex := latestBlock.Index + 1

	newBlock, err := blockchain.CreateBlock(
		batchID,
		newIndex,
		eventType,
		string(eventJSON),
		actorID,
		actorRole,
		latestBlock.Hash, // prevHash = latest block's hash
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create new block: %w", err)
	}

	// Step 6: Save new block to database
	if err := tx.Create(&newBlock).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to save new block to database: %w", err)
	}

	// Step 7: Update batch current hash and status based on event type
	batch.CurrentHash = newBlock.Hash
	batch.UpdatedAt = time.Now().UTC()

	// Update status based on event type
	switch eventType {
	case "transport":
		batch.Status = "in_transit"
	case "transfer":
		batch.Status = "delivered"
	case "sale":
		batch.Status = "sold"
		// For "registration" and "quality_check", status remains unchanged
	}

	if err := tx.Save(&batch).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update batch: %w", err)
	}

	// Step 8: Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
