package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"

	"shambachain/models"
)

// ComputeBlockHash computes the SHA-256 hash of a block
// by concatenating block fields in deterministic order.
//
// Preconditions:
//   - block is a valid Block structure with all required fields populated
//
// Postconditions:
//   - Returns 64-character hexadecimal string (SHA-256 hash)
//   - Hash is deterministic: same block data always produces same hash
//   - Hash computation includes: BatchID, Index, Timestamp, EventType, EventData, ActorID, ActorRole, PrevHash
//   - No side effects on input block
func ComputeBlockHash(block models.Block) string {
	// Step 1: Concatenate block fields in deterministic order
	// Format timestamp as RFC3339 for consistency
	timestampStr := block.Timestamp.UTC().Format("2006-01-02T15:04:05.999999999Z07:00")

	record := fmt.Sprintf("%s%s%s%s%s%s%s%s",
		block.BatchID,
		strconv.Itoa(block.Index),
		timestampStr,
		block.EventType,
		block.EventData,
		block.ActorID,
		block.ActorRole,
		block.PrevHash,
	)

	// Step 2: Compute SHA-256 hash
	hashBytes := sha256.Sum256([]byte(record))

	// Step 3: Convert to hexadecimal string
	hash := hex.EncodeToString(hashBytes[:])

	return hash
}
