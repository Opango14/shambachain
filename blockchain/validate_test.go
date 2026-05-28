package blockchain

import (
	"testing"
	"time"

	"shambachain/models"
)

// TestValidateChain_ValidSingleBlock tests validation of a valid single-block chain (genesis only)
func TestValidateChain_ValidSingleBlock(t *testing.T) {
	timestamp := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	genesisBlock := models.Block{
		BatchID:   "batch-123",
		Index:     0,
		Timestamp: timestamp,
		EventType: "registration",
		EventData: `{"produce_type":"potatoes","quantity":500}`,
		ActorID:   "farmer-001",
		ActorRole: "farmer",
		PrevHash:  "0",
	}
	genesisBlock.Hash = ComputeBlockHash(genesisBlock)

	blocks := []models.Block{genesisBlock}

	// Validates: Requirements 4.1, 4.2, 4.7 (Property 20)
	if !ValidateChain(blocks) {
		t.Error("Valid single-block chain should validate successfully")
	}
}

// TestValidateChain_ValidMultipleBlocks tests validation of a valid multi-block chain
func TestValidateChain_ValidMultipleBlocks(t *testing.T) {
	timestamp1 := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	timestamp2 := time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC)
	timestamp3 := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	// Genesis block
	block0 := models.Block{
		BatchID:   "batch-456",
		Index:     0,
		Timestamp: timestamp1,
		EventType: "registration",
		EventData: `{"produce_type":"maize","quantity":1000}`,
		ActorID:   "farmer-002",
		ActorRole: "farmer",
		PrevHash:  "0",
	}
	block0.Hash = ComputeBlockHash(block0)

	// Transport block
	block1 := models.Block{
		BatchID:   "batch-456",
		Index:     1,
		Timestamp: timestamp2,
		EventType: "transport",
		EventData: `{"from_location":"Farm A","to_location":"Market B"}`,
		ActorID:   "transporter-001",
		ActorRole: "transporter",
		PrevHash:  block0.Hash,
	}
	block1.Hash = ComputeBlockHash(block1)

	// Quality check block
	block2 := models.Block{
		BatchID:   "batch-456",
		Index:     2,
		Timestamp: timestamp3,
		EventType: "quality_check",
		EventData: `{"inspector_id":"inspector-001","grade":"A","passed":true}`,
		ActorID:   "inspector-001",
		ActorRole: "inspector",
		PrevHash:  block1.Hash,
	}
	block2.Hash = ComputeBlockHash(block2)

	blocks := []models.Block{block0, block1, block2}

	// Validates: Requirements 4.1, 4.2, 4.3, 4.4, 4.5, 4.7 (Property 20)
	if !ValidateChain(blocks) {
		t.Error("Valid multi-block chain should validate successfully")
	}
}

// TestValidateChain_EmptyChain tests validation of empty chain
func TestValidateChain_EmptyChain(t *testing.T) {
	blocks := []models.Block{}

	// Validates: Requirements 4.6, 4.7
	if ValidateChain(blocks) {
		t.Error("Empty chain should fail validation")
	}
}

// TestValidateChain_NilChain tests validation of nil chain
func TestValidateChain_NilChain(t *testing.T) {
	// Validates: Requirements 4.6, 4.7
	if ValidateChain(nil) {
		t.Error("Nil chain should fail validation")
	}
}

// TestValidateChain_InvalidGenesisIndex tests genesis block with wrong index
func TestValidateChain_InvalidGenesisIndex(t *testing.T) {
	timestamp := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	genesisBlock := models.Block{
		BatchID:   "batch-789",
		Index:     1, // Wrong index - should be 0
		Timestamp: timestamp,
		EventType: "registration",
		EventData: `{"produce_type":"fish","quantity":200}`,
		ActorID:   "farmer-003",
		ActorRole: "farmer",
		PrevHash:  "0",
	}
	genesisBlock.Hash = ComputeBlockHash(genesisBlock)

	blocks := []models.Block{genesisBlock}

	// Validates: Requirements 4.1, 4.6, 4.7 (Property 21)
	if ValidateChain(blocks) {
		t.Error("Chain with invalid genesis index should fail validation")
	}
}

// TestValidateChain_InvalidGenesisPrevHash tests genesis block with wrong prevHash
func TestValidateChain_InvalidGenesisPrevHash(t *testing.T) {
	timestamp := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	genesisBlock := models.Block{
		BatchID:   "batch-999",
		Index:     0,
		Timestamp: timestamp,
		EventType: "registration",
		EventData: `{"produce_type":"vegetables","quantity":300}`,
		ActorID:   "farmer-004",
		ActorRole: "farmer",
		PrevHash:  "abc123", // Wrong prevHash - should be "0"
	}
	genesisBlock.Hash = ComputeBlockHash(genesisBlock)

	blocks := []models.Block{genesisBlock}

	// Validates: Requirements 4.1, 4.6, 4.7 (Property 21)
	if ValidateChain(blocks) {
		t.Error("Chain with invalid genesis prevHash should fail validation")
	}
}

// TestValidateChain_TamperedHash tests detection of tampered block hash
func TestValidateChain_TamperedHash(t *testing.T) {
	timestamp := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	genesisBlock := models.Block{
		BatchID:   "batch-111",
		Index:     0,
		Timestamp: timestamp,
		EventType: "registration",
		EventData: `{"produce_type":"potatoes","quantity":500}`,
		ActorID:   "farmer-005",
		ActorRole: "farmer",
		PrevHash:  "0",
	}
	genesisBlock.Hash = ComputeBlockHash(genesisBlock)

	// Tamper with the hash
	genesisBlock.Hash = "0000000000000000000000000000000000000000000000000000000000000000"

	blocks := []models.Block{genesisBlock}

	// Validates: Requirements 4.2, 4.6, 4.7, 12.4 (Property 21)
	if ValidateChain(blocks) {
		t.Error("Chain with tampered hash should fail validation")
	}
}

// TestValidateChain_TamperedData tests detection of tampered block data
func TestValidateChain_TamperedData(t *testing.T) {
	timestamp := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	genesisBlock := models.Block{
		BatchID:   "batch-222",
		Index:     0,
		Timestamp: timestamp,
		EventType: "registration",
		EventData: `{"produce_type":"potatoes","quantity":500}`,
		ActorID:   "farmer-006",
		ActorRole: "farmer",
		PrevHash:  "0",
	}
	genesisBlock.Hash = ComputeBlockHash(genesisBlock)

	// Tamper with the data after hash computation
	genesisBlock.EventData = `{"produce_type":"potatoes","quantity":5000}`

	blocks := []models.Block{genesisBlock}

	// Validates: Requirements 4.2, 4.6, 4.7, 12.3, 12.4 (Property 2, Property 21)
	if ValidateChain(blocks) {
		t.Error("Chain with tampered data should fail validation")
	}
}

// TestValidateChain_BrokenHashLink tests detection of broken hash link between blocks
func TestValidateChain_BrokenHashLink(t *testing.T) {
	timestamp1 := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	timestamp2 := time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC)

	block0 := models.Block{
		BatchID:   "batch-333",
		Index:     0,
		Timestamp: timestamp1,
		EventType: "registration",
		EventData: `{"produce_type":"maize","quantity":1000}`,
		ActorID:   "farmer-007",
		ActorRole: "farmer",
		PrevHash:  "0",
	}
	block0.Hash = ComputeBlockHash(block0)

	block1 := models.Block{
		BatchID:   "batch-333",
		Index:     1,
		Timestamp: timestamp2,
		EventType: "transport",
		EventData: `{"from_location":"Farm A","to_location":"Market B"}`,
		ActorID:   "transporter-002",
		ActorRole: "transporter",
		PrevHash:  "wronghash1234567890123456789012345678901234567890123456789012", // Wrong prevHash
	}
	block1.Hash = ComputeBlockHash(block1)

	blocks := []models.Block{block0, block1}

	// Validates: Requirements 4.3, 4.6, 4.7 (Property 4, Property 21)
	if ValidateChain(blocks) {
		t.Error("Chain with broken hash link should fail validation")
	}
}

// TestValidateChain_NonSequentialIndices tests detection of non-sequential indices
func TestValidateChain_NonSequentialIndices(t *testing.T) {
	timestamp1 := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	timestamp2 := time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC)

	block0 := models.Block{
		BatchID:   "batch-444",
		Index:     0,
		Timestamp: timestamp1,
		EventType: "registration",
		EventData: `{"produce_type":"fish","quantity":200}`,
		ActorID:   "farmer-008",
		ActorRole: "farmer",
		PrevHash:  "0",
	}
	block0.Hash = ComputeBlockHash(block0)

	block1 := models.Block{
		BatchID:   "batch-444",
		Index:     3, // Wrong index - should be 1
		Timestamp: timestamp2,
		EventType: "transport",
		EventData: `{"from_location":"Farm C","to_location":"Market D"}`,
		ActorID:   "transporter-003",
		ActorRole: "transporter",
		PrevHash:  block0.Hash,
	}
	block1.Hash = ComputeBlockHash(block1)

	blocks := []models.Block{block0, block1}

	// Validates: Requirements 4.4, 4.6, 4.7 (Property 4, Property 21)
	if ValidateChain(blocks) {
		t.Error("Chain with non-sequential indices should fail validation")
	}
}

// TestValidateChain_InvalidTimestampOrdering tests detection of invalid timestamp ordering
func TestValidateChain_InvalidTimestampOrdering(t *testing.T) {
	timestamp1 := time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC)
	timestamp2 := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC) // Earlier than block0

	block0 := models.Block{
		BatchID:   "batch-555",
		Index:     0,
		Timestamp: timestamp1,
		EventType: "registration",
		EventData: `{"produce_type":"vegetables","quantity":300}`,
		ActorID:   "farmer-009",
		ActorRole: "farmer",
		PrevHash:  "0",
	}
	block0.Hash = ComputeBlockHash(block0)

	block1 := models.Block{
		BatchID:   "batch-555",
		Index:     1,
		Timestamp: timestamp2, // Invalid - earlier than previous block
		EventType: "transport",
		EventData: `{"from_location":"Farm E","to_location":"Market F"}`,
		ActorID:   "transporter-004",
		ActorRole: "transporter",
		PrevHash:  block0.Hash,
	}
	block1.Hash = ComputeBlockHash(block1)

	blocks := []models.Block{block0, block1}

	// Validates: Requirements 4.5, 4.6, 4.7 (Property 5, Property 21)
	if ValidateChain(blocks) {
		t.Error("Chain with invalid timestamp ordering should fail validation")
	}
}

// TestValidateChain_EqualTimestamps tests validation with equal timestamps (should be valid)
func TestValidateChain_EqualTimestamps(t *testing.T) {
	timestamp := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	block0 := models.Block{
		BatchID:   "batch-666",
		Index:     0,
		Timestamp: timestamp,
		EventType: "registration",
		EventData: `{"produce_type":"potatoes","quantity":500}`,
		ActorID:   "farmer-010",
		ActorRole: "farmer",
		PrevHash:  "0",
	}
	block0.Hash = ComputeBlockHash(block0)

	block1 := models.Block{
		BatchID:   "batch-666",
		Index:     1,
		Timestamp: timestamp, // Same timestamp - should be valid (monotonically non-decreasing)
		EventType: "quality_check",
		EventData: `{"inspector_id":"inspector-002","grade":"A","passed":true}`,
		ActorID:   "inspector-002",
		ActorRole: "inspector",
		PrevHash:  block0.Hash,
	}
	block1.Hash = ComputeBlockHash(block1)

	blocks := []models.Block{block0, block1}

	// Validates: Requirements 4.5, 4.7 (Property 5)
	if !ValidateChain(blocks) {
		t.Error("Chain with equal timestamps should validate successfully (monotonically non-decreasing)")
	}
}

// TestValidateChain_LongChain tests validation of a longer chain
func TestValidateChain_LongChain(t *testing.T) {
	timestamp1 := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	timestamp2 := time.Date(2024, 1, 15, 11, 0, 0, 0, time.UTC)
	timestamp3 := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	timestamp4 := time.Date(2024, 1, 15, 13, 0, 0, 0, time.UTC)
	timestamp5 := time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC)

	// Genesis block
	block0 := models.Block{
		BatchID:   "batch-777",
		Index:     0,
		Timestamp: timestamp1,
		EventType: "registration",
		EventData: `{"produce_type":"maize","quantity":1000}`,
		ActorID:   "farmer-011",
		ActorRole: "farmer",
		PrevHash:  "0",
	}
	block0.Hash = ComputeBlockHash(block0)

	// Transport block
	block1 := models.Block{
		BatchID:   "batch-777",
		Index:     1,
		Timestamp: timestamp2,
		EventType: "transport",
		EventData: `{"from_location":"Farm G","to_location":"Warehouse H"}`,
		ActorID:   "transporter-005",
		ActorRole: "transporter",
		PrevHash:  block0.Hash,
	}
	block1.Hash = ComputeBlockHash(block1)

	// Quality check block
	block2 := models.Block{
		BatchID:   "batch-777",
		Index:     2,
		Timestamp: timestamp3,
		EventType: "quality_check",
		EventData: `{"inspector_id":"inspector-003","grade":"B","passed":true}`,
		ActorID:   "inspector-003",
		ActorRole: "inspector",
		PrevHash:  block1.Hash,
	}
	block2.Hash = ComputeBlockHash(block2)

	// Transfer block
	block3 := models.Block{
		BatchID:   "batch-777",
		Index:     3,
		Timestamp: timestamp4,
		EventType: "transfer",
		EventData: `{"from_owner_id":"farmer-011","to_owner_id":"buyer-001","transfer_type":"sale"}`,
		ActorID:   "farmer-011",
		ActorRole: "farmer",
		PrevHash:  block2.Hash,
	}
	block3.Hash = ComputeBlockHash(block3)

	// Sale block
	block4 := models.Block{
		BatchID:   "batch-777",
		Index:     4,
		Timestamp: timestamp5,
		EventType: "sale",
		EventData: `{"price":50000,"currency":"KES"}`,
		ActorID:   "buyer-001",
		ActorRole: "buyer",
		PrevHash:  block3.Hash,
	}
	block4.Hash = ComputeBlockHash(block4)

	blocks := []models.Block{block0, block1, block2, block3, block4}

	// Validates: Requirements 4.1, 4.2, 4.3, 4.4, 4.5, 4.7 (Property 20)
	if !ValidateChain(blocks) {
		t.Error("Valid long chain should validate successfully")
	}
}

// TestValidateChain_NoMutations tests that validation does not modify input blocks
func TestValidateChain_NoMutations(t *testing.T) {
	timestamp := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	originalBlock := models.Block{
		BatchID:   "batch-888",
		Index:     0,
		Timestamp: timestamp,
		EventType: "registration",
		EventData: `{"produce_type":"potatoes","quantity":500}`,
		ActorID:   "farmer-012",
		ActorRole: "farmer",
		PrevHash:  "0",
	}
	originalBlock.Hash = ComputeBlockHash(originalBlock)

	// Create a copy to compare later
	originalHash := originalBlock.Hash
	originalData := originalBlock.EventData
	originalIndex := originalBlock.Index

	blocks := []models.Block{originalBlock}

	ValidateChain(blocks)

	// Validates: Requirements 4.8 (Property 24)
	if blocks[0].Hash != originalHash {
		t.Error("Validation modified block hash")
	}
	if blocks[0].EventData != originalData {
		t.Error("Validation modified block data")
	}
	if blocks[0].Index != originalIndex {
		t.Error("Validation modified block index")
	}
}
