package blockchain

import (
	"fmt"
	"time"

	"shambachain/models"
)

// Valid event types
var validEventTypes = map[string]bool{
	"registration":  true,
	"transport":     true,
	"quality_check": true,
	"transfer":      true,
	"sale":          true,
}

// Valid actor roles
var validActorRoles = map[string]bool{
	"farmer":      true,
	"transporter": true,
	"inspector":   true,
	"buyer":       true,
}

// CreateBlock creates a new blockchain block with all fields populated.
//
// Preconditions:
//   - batchID is non-empty string
//   - index is non-negative integer (>= 0)
//   - eventType is valid event type (registration, transport, quality_check, transfer, sale)
//   - eventData is valid JSON string
//   - actorID is non-empty string
//   - actorRole is valid role (farmer, transporter, buyer, inspector)
//   - prevHash is valid hash string (64 hex characters for SHA-256)
//
// Postconditions:
//   - Returns Block with all fields populated
//   - block.Timestamp is set to current UTC time
//   - block.Hash is computed from block contents using SHA-256
//   - block.Hash is deterministic: same inputs always produce same hash
//   - block.PrevHash equals input prevHash
//   - Block is immutable once created
func CreateBlock(batchID string, index int, eventType string, eventData string,
	actorID string, actorRole string, prevHash string) (models.Block, error) {

	// Validate preconditions
	if batchID == "" {
		return models.Block{}, fmt.Errorf("batchID cannot be empty")
	}

	if index < 0 {
		return models.Block{}, fmt.Errorf("index must be non-negative, got %d", index)
	}

	if !validEventTypes[eventType] {
		return models.Block{}, fmt.Errorf("invalid event type: %s (must be one of: registration, transport, quality_check, transfer, sale)", eventType)
	}

	if actorID == "" {
		return models.Block{}, fmt.Errorf("actorID cannot be empty")
	}

	if !validActorRoles[actorRole] {
		return models.Block{}, fmt.Errorf("invalid actor role: %s (must be one of: farmer, transporter, inspector, buyer)", actorRole)
	}

	if prevHash == "" {
		return models.Block{}, fmt.Errorf("prevHash cannot be empty")
	}

	if len(prevHash) != 64 && prevHash != "0" {
		return models.Block{}, fmt.Errorf("prevHash must be 64 hex characters or '0' for genesis block, got length %d", len(prevHash))
	}

	// Step 1: Create block structure with all fields
	block := models.Block{
		BatchID:   batchID,
		Index:     index,
		Timestamp: time.Now().UTC(),
		EventType: eventType,
		EventData: eventData,
		ActorID:   actorID,
		ActorRole: actorRole,
		PrevHash:  prevHash,
		CreatedAt: time.Now().UTC(),
	}

	// Step 2: Compute and assign hash
	block.Hash = ComputeBlockHash(block)

	return block, nil
}
