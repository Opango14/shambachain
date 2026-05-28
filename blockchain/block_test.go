package blockchain

import (
	"encoding/json"
	"testing"
	"time"

	"shambachain/models"
)

func TestCreateBlock_Success(t *testing.T) {
	// Test data
	batchID := "batch-123"
	index := 1
	eventType := "transport"
	eventData := `{"from_location":"Farm A","to_location":"Market B"}`
	actorID := "actor-456"
	actorRole := "transporter"
	prevHash := "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"

	block, err := CreateBlock(batchID, index, eventType, eventData, actorID, actorRole, prevHash)

	if err != nil {
		t.Fatalf("CreateBlock failed: %v", err)
	}

	// Verify all fields are set correctly
	if block.BatchID != batchID {
		t.Errorf("Expected BatchID %s, got %s", batchID, block.BatchID)
	}

	if block.Index != index {
		t.Errorf("Expected Index %d, got %d", index, block.Index)
	}

	if block.EventType != eventType {
		t.Errorf("Expected EventType %s, got %s", eventType, block.EventType)
	}

	if block.EventData != eventData {
		t.Errorf("Expected EventData %s, got %s", eventData, block.EventData)
	}

	if block.ActorID != actorID {
		t.Errorf("Expected ActorID %s, got %s", actorID, block.ActorID)
	}

	if block.ActorRole != actorRole {
		t.Errorf("Expected ActorRole %s, got %s", actorRole, block.ActorRole)
	}

	if block.PrevHash != prevHash {
		t.Errorf("Expected PrevHash %s, got %s", prevHash, block.PrevHash)
	}

	// Verify timestamp is recent (within last second)
	now := time.Now().UTC()
	if block.Timestamp.After(now) || block.Timestamp.Before(now.Add(-1*time.Second)) {
		t.Errorf("Timestamp not set to current UTC time: %v", block.Timestamp)
	}

	// Verify hash is computed and has correct length
	if len(block.Hash) != 64 {
		t.Errorf("Expected hash length 64, got %d", len(block.Hash))
	}

	// Verify hash is deterministic by recomputing
	recomputedHash := ComputeBlockHash(block)
	if block.Hash != recomputedHash {
		t.Errorf("Hash mismatch: stored=%s, recomputed=%s", block.Hash, recomputedHash)
	}
}

func TestCreateBlock_GenesisBlock(t *testing.T) {
	// Test genesis block creation with prevHash "0"
	eventData := `{"produce_type":"potatoes","quantity":500}`

	block, err := CreateBlock("batch-001", 0, "registration", eventData, "farmer-001", "farmer", "0")

	if err != nil {
		t.Fatalf("CreateBlock failed for genesis block: %v", err)
	}

	if block.Index != 0 {
		t.Errorf("Expected genesis block index 0, got %d", block.Index)
	}

	if block.PrevHash != "0" {
		t.Errorf("Expected genesis block prevHash '0', got %s", block.PrevHash)
	}

	if len(block.Hash) != 64 {
		t.Errorf("Expected hash length 64, got %d", len(block.Hash))
	}
}

func TestCreateBlock_AllEventTypes(t *testing.T) {
	eventTypes := []string{"registration", "transport", "quality_check", "transfer", "sale"}
	prevHash := "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"

	for _, eventType := range eventTypes {
		block, err := CreateBlock("batch-123", 1, eventType, "{}", "actor-123", "farmer", prevHash)
		if err != nil {
			t.Errorf("CreateBlock failed for event type %s: %v", eventType, err)
		}
		if block.EventType != eventType {
			t.Errorf("Expected event type %s, got %s", eventType, block.EventType)
		}
	}
}

func TestCreateBlock_AllActorRoles(t *testing.T) {
	actorRoles := []string{"farmer", "transporter", "inspector", "buyer"}
	prevHash := "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"

	for _, role := range actorRoles {
		block, err := CreateBlock("batch-123", 1, "transport", "{}", "actor-123", role, prevHash)
		if err != nil {
			t.Errorf("CreateBlock failed for actor role %s: %v", role, err)
		}
		if block.ActorRole != role {
			t.Errorf("Expected actor role %s, got %s", role, block.ActorRole)
		}
	}
}

func TestCreateBlock_InvalidEventType(t *testing.T) {
	prevHash := "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"

	_, err := CreateBlock("batch-123", 1, "invalid_event", "{}", "actor-123", "farmer", prevHash)

	if err == nil {
		t.Error("Expected error for invalid event type, got nil")
	}
}

func TestCreateBlock_InvalidActorRole(t *testing.T) {
	prevHash := "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"

	_, err := CreateBlock("batch-123", 1, "transport", "{}", "actor-123", "invalid_role", prevHash)

	if err == nil {
		t.Error("Expected error for invalid actor role, got nil")
	}
}

func TestCreateBlock_EmptyBatchID(t *testing.T) {
	prevHash := "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"

	_, err := CreateBlock("", 1, "transport", "{}", "actor-123", "farmer", prevHash)

	if err == nil {
		t.Error("Expected error for empty batchID, got nil")
	}
}

func TestCreateBlock_NegativeIndex(t *testing.T) {
	prevHash := "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"

	_, err := CreateBlock("batch-123", -1, "transport", "{}", "actor-123", "farmer", prevHash)

	if err == nil {
		t.Error("Expected error for negative index, got nil")
	}
}

func TestCreateBlock_EmptyActorID(t *testing.T) {
	prevHash := "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"

	_, err := CreateBlock("batch-123", 1, "transport", "{}", "", "farmer", prevHash)

	if err == nil {
		t.Error("Expected error for empty actorID, got nil")
	}
}

func TestCreateBlock_EmptyPrevHash(t *testing.T) {
	_, err := CreateBlock("batch-123", 1, "transport", "{}", "actor-123", "farmer", "")

	if err == nil {
		t.Error("Expected error for empty prevHash, got nil")
	}
}

func TestCreateBlock_InvalidPrevHashLength(t *testing.T) {
	_, err := CreateBlock("batch-123", 1, "transport", "{}", "actor-123", "farmer", "short")

	if err == nil {
		t.Error("Expected error for invalid prevHash length, got nil")
	}
}

func TestCreateBlock_HashDeterminism(t *testing.T) {
	// Create the same block twice and verify hashes match
	batchID := "batch-123"
	index := 1
	eventType := "transport"
	eventData := `{"from":"A","to":"B"}`
	actorID := "actor-456"
	actorRole := "transporter"
	prevHash := "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"

	// Create first block
	block1, err1 := CreateBlock(batchID, index, eventType, eventData, actorID, actorRole, prevHash)
	if err1 != nil {
		t.Fatalf("First CreateBlock failed: %v", err1)
	}

	// Wait a tiny bit to ensure different timestamp
	time.Sleep(1 * time.Millisecond)

	// Create second block with same parameters
	block2, err2 := CreateBlock(batchID, index, eventType, eventData, actorID, actorRole, prevHash)
	if err2 != nil {
		t.Fatalf("Second CreateBlock failed: %v", err2)
	}

	// Hashes will be different because timestamps are different
	// But if we manually set the same timestamp, hashes should match
	block2.Timestamp = block1.Timestamp
	block2.CreatedAt = block1.CreatedAt

	hash1 := ComputeBlockHash(block1)
	hash2 := ComputeBlockHash(block2)

	if hash1 != hash2 {
		t.Errorf("Hash determinism failed: hash1=%s, hash2=%s", hash1, hash2)
	}
}

func TestCreateBlock_WithRealEventData(t *testing.T) {
	// Test with actual event data structures
	transportEvent := models.TransportEvent{
		FromLocation:     "Nakuru Farm",
		ToLocation:       "Nairobi Market",
		TransportID:      "TRK-001",
		VehicleInfo:      "Truck KBZ 123A",
		DepartureTime:    time.Now(),
		EstimatedArrival: time.Now().Add(4 * time.Hour),
	}

	eventDataBytes, err := json.Marshal(transportEvent)
	if err != nil {
		t.Fatalf("Failed to marshal event data: %v", err)
	}

	prevHash := "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"
	block, err := CreateBlock("batch-123", 1, "transport", string(eventDataBytes), "transporter-456", "transporter", prevHash)

	if err != nil {
		t.Fatalf("CreateBlock failed with real event data: %v", err)
	}

	// Verify we can unmarshal the event data back
	var decoded models.TransportEvent
	if err := json.Unmarshal([]byte(block.EventData), &decoded); err != nil {
		t.Errorf("Failed to unmarshal event data from block: %v", err)
	}

	if decoded.TransportID != transportEvent.TransportID {
		t.Errorf("Event data mismatch: expected TransportID %s, got %s", transportEvent.TransportID, decoded.TransportID)
	}
}
