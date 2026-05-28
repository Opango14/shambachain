package blockchain

import (
	"shambachain/models"
)

// ValidateChain validates the integrity of a blockchain.
//
// Preconditions:
//   - blocks is a slice of Block structures
//   - blocks is ordered by Index (ascending)
//   - len(blocks) >= 1 (at least genesis block exists)
//
// Postconditions:
//   - Returns true if and only if chain is valid
//   - Chain is valid when:
//   - Genesis block (index 0) has PrevHash = "0"
//   - Each block's Hash matches recomputed hash from its data
//   - Each block's PrevHash matches previous block's Hash
//   - Block indices are sequential (0, 1, 2, ...)
//   - Block timestamps are monotonically non-decreasing
//   - Returns false if any validation check fails
//   - No mutations to input blocks
//
// Loop Invariants:
//   - For each iteration i (1 to len(blocks)-1):
//   - All blocks from 0 to i-1 have been validated
//   - blocks[i].PrevHash == blocks[i-1].Hash
//   - blocks[i].Hash == ComputeBlockHash(blocks[i])
//   - blocks[i].Index == i
func ValidateChain(blocks []models.Block) bool {
	// Precondition checks
	if blocks == nil || len(blocks) < 1 {
		return false
	}

	// Step 1: Validate genesis block
	genesisBlock := blocks[0]

	// Genesis block must have index 0
	if genesisBlock.Index != 0 {
		return false
	}

	// Genesis block must have previous hash "0"
	if genesisBlock.PrevHash != "0" {
		return false
	}

	// Genesis block hash must match recomputed hash
	recomputedHash := ComputeBlockHash(genesisBlock)
	if genesisBlock.Hash != recomputedHash {
		return false
	}

	// Step 2: Validate chain links
	// Loop Invariant: All blocks from 0 to i-1 are valid
	for i := 1; i < len(blocks); i++ {
		currentBlock := blocks[i]
		previousBlock := blocks[i-1]

		// Check sequential indices
		if currentBlock.Index != i {
			return false
		}

		// Check hash link: current block's PrevHash must equal previous block's Hash
		if currentBlock.PrevHash != previousBlock.Hash {
			return false
		}

		// Check hash integrity: stored hash must match recomputed hash
		recomputedHash := ComputeBlockHash(currentBlock)
		if currentBlock.Hash != recomputedHash {
			return false
		}

		// Check timestamp ordering: timestamps must be monotonically non-decreasing
		if currentBlock.Timestamp.Before(previousBlock.Timestamp) {
			return false
		}
	}

	// All validations passed
	return true
}
