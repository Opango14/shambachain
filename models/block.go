package models

import "time"

// Block represents a blockchain block for traceability
type Block struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	BatchID   string    `gorm:"index" json:"batch_id"`
	Index     int       `json:"index"`
	Timestamp time.Time `json:"timestamp"`
	EventType string    `json:"event_type"` // registration, transport, quality_check, transfer, sale
	EventData string    `json:"event_data"` // JSON string with event details
	ActorID   string    `json:"actor_id"`   // user who created this event
	ActorRole string    `json:"actor_role"` // farmer, transporter, buyer, inspector
	PrevHash  string    `json:"prev_hash"`
	Hash      string    `json:"hash"`
	CreatedAt time.Time `json:"created_at"`
}
