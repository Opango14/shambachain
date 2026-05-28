package main

import (
	"crypto/sha256"
	"fmt"
)

type Block struct {
	Index    int
	Data     string
	PrevHash string
}

func (b Block) Hash() string {
	record := fmt.Sprintf("%d%s%s", b.Index, b.Data, b.PrevHash)
	h := sha256.Sum256([]byte(record))
	return fmt.Sprintf("%x", h)
}

func main() {
	genesis := Block{0, "Genesis Block", "0"}
	block1 := Block{1, "Transfer 10 BTC", genesis.Hash()}

	fmt.Println("Block 0 hash:", genesis.Hash())
	fmt.Println("Block 1 hash:", block1.Hash())
}
