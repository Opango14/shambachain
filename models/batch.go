package models

import "time"

// Batch represents a registered produce batch
type Batch struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	FarmerID    string    `gorm:"index" json:"farmer_id"`
	ProduceType string    `json:"produce_type"` // potatoes, fish, vegetables, maize
	Quantity    float64   `json:"quantity"`     // in kg
	Unit        string    `json:"unit"`         // kg, tons, pieces
	HarvestDate time.Time `json:"harvest_date"`
	Location    string    `json:"location"`     // GPS coordinates or farm name
	Status      string    `json:"status"`       // registered, in_transit, delivered, sold
	QRCodePath  string    `json:"qr_code_path"` // path to QR code image
	GenesisHash string    `json:"genesis_hash"` // first block hash
	CurrentHash string    `json:"current_hash"` // latest block hash
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
