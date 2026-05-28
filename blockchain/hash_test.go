package blockchain

import (
	"testing"
	"time"

	"shambachain/models"
)

func TestComputeBlockHash(t *testing.T) {
	// Test case 1: Hash determinism - same block data should produce same hash
	timestamp := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	block1 := models.Block{
		BatchID:   "batch-123",
		Index:     0,
		Timestamp: timestamp,
		EventType: "registration",
		EventData: `{"produce_type":"potatoes","quantity":500}`,
		ActorID:   "farmer-001",
		ActorRole: "farmer",
		PrevHash:  "0",
	}

	block2 := models.Block{
		BatchID:   "batch-123",
		Index:     0,
		Timestamp: timestamp,
		EventType: "registration",
		EventData: `{"produce_type":"potatoes","quantity":500}`,
		ActorID:   "farmer-001",
		ActorRole: "farmer",
		PrevHash:  "0",
	}

	hash1 := ComputeBlockHash(block1)
	hash2 := ComputeBlockHash(block2)

	// Verify hash determinism (Property 1)
	if hash1 != hash2 {
		t.Errorf("Hash determinism failed: expected identical hashes for identical blocks, got %s and %s", hash1, hash2)
	}

	// Verify hash format (Property 10) - should be 64 hex characters
	if len(hash1) != 64 {
		t.Errorf("Hash format invalid: expected 64 characters, got %d", len(hash1))
	}

	// Verify hash is hexadecimal
	for _, c := range hash1 {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Errorf("Hash contains non-hexadecimal character: %c", c)
		}
	}
}

func TestComputeBlockHash_DifferentData(t *testing.T) {
	// Test case 2: Different block data should produce different hashes
	timestamp := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	block1 := models.Block{
		BatchID:   "batch-123",
		Index:     0,
		Timestamp: timestamp,
		EventType: "registration",
		EventData: `{"produce_type":"potatoes","quantity":500}`,
		ActorID:   "farmer-001",
		ActorRole: "farmer",
		PrevHash:  "0",
	}

	block2 := models.Block{
		BatchID:   "batch-456", // Different batch ID
		Index:     0,
		Timestamp: timestamp,
		EventType: "registration",
		EventData: `{"produce_type":"potatoes","quantity":500}`,
		ActorID:   "farmer-001",
		ActorRole: "farmer",
		PrevHash:  "0",
	}

	hash1 := ComputeBlockHash(block1)
	hash2 := ComputeBlockHash(block2)

	if hash1 == hash2 {
		t.Errorf("Different blocks produced same hash: %s", hash1)
	}
}

func TestComputeBlockHash_GenesisBlock(t *testing.T) {
	// Test case 3: Genesis block with prevHash "0"
	timestamp := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	genesisBlock := models.Block{
		BatchID:   "batch-789",
		Index:     0,
		Timestamp: timestamp,
		EventType: "registration",
		EventData: `{"produce_type":"maize","quantity":1000}`,
		ActorID:   "farmer-002",
		ActorRole: "farmer",
		PrevHash:  "0",
	}

	hash := ComputeBlockHash(genesisBlock)

	// Verify hash is computed correctly
	if hash == "" {
		t.Error("Genesis block hash is empty")
	}

	if len(hash) != 64 {
		t.Errorf("Genesis block hash length invalid: expected 64, got %d", len(hash))
	}
}

func TestComputeBlockHash_SequentialBlocks(t *testing.T) {
	// Test case 4: Sequential blocks with linked hashes
	timestamp1 := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	timestamp2 := time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC)

	block1 := models.Block{
		BatchID:   "batch-999",
		Index:     0,
		Timestamp: timestamp1,
		EventType: "registration",
		EventData: `{"produce_type":"fish","quantity":200}`,
		ActorID:   "farmer-003",
		ActorRole: "farmer",
		PrevHash:  "0",
	}

	hash1 := ComputeBlockHash(block1)

	block2 := models.Block{
		BatchID:   "batch-999",
		Index:     1,
		Timestamp: timestamp2,
		EventType: "transport",
		EventData: `{"from_location":"Farm A","to_location":"Market B"}`,
		ActorID:   "transporter-001",
		ActorRole: "transporter",
		PrevHash:  hash1, // Link to previous block
	}

	hash2 := ComputeBlockHash(block2)

	// Verify hashes are different
	if hash1 == hash2 {
		t.Error("Sequential blocks produced identical hashes")
	}

	// Verify both hashes are valid
	if len(hash1) != 64 || len(hash2) != 64 {
		t.Errorf("Invalid hash lengths: hash1=%d, hash2=%d", len(hash1), len(hash2))
	}
}
