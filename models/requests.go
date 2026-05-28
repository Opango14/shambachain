package models

import "time"

// RegisterBatchRequest represents the request body for batch registration
type RegisterBatchRequest struct {
	ProduceType string    `json:"produce_type" binding:"required"`
	Quantity    float64   `json:"quantity" binding:"required,gt=0"`
	Unit        string    `json:"unit" binding:"required"`
	HarvestDate time.Time `json:"harvest_date" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	FarmName    string    `json:"farm_name"`
}

// AddEventRequest represents the request body for adding events to a batch
type AddEventRequest struct {
	EventType string                 `json:"event_type" binding:"required"`
	EventData map[string]interface{} `json:"event_data" binding:"required"`
}
