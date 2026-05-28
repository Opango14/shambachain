package models

import "time"

// RegistrationEvent represents the data for a produce registration event
type RegistrationEvent struct {
	ProduceType string    `json:"produce_type"`
	Quantity    float64   `json:"quantity"`
	Unit        string    `json:"unit"`
	HarvestDate time.Time `json:"harvest_date"`
	Location    string    `json:"location"`
	FarmName    string    `json:"farm_name"`
}

// TransportEvent represents the data for a transport event
type TransportEvent struct {
	FromLocation     string    `json:"from_location"`
	ToLocation       string    `json:"to_location"`
	TransportID      string    `json:"transport_id"`
	VehicleInfo      string    `json:"vehicle_info"`
	DepartureTime    time.Time `json:"departure_time"`
	EstimatedArrival time.Time `json:"estimated_arrival"`
}

// QualityCheckEvent represents the data for a quality inspection event
type QualityCheckEvent struct {
	InspectorID   string  `json:"inspector_id"`
	InspectorName string  `json:"inspector_name"`
	Grade         string  `json:"grade"` // A, B, C
	Notes         string  `json:"notes"`
	Temperature   float64 `json:"temperature"` // for perishables
	Passed        bool    `json:"passed"`
}

// TransferEvent represents the data for an ownership transfer event
type TransferEvent struct {
	FromOwnerID  string  `json:"from_owner_id"`
	ToOwnerID    string  `json:"to_owner_id"`
	TransferType string  `json:"transfer_type"` // sale, donation, return
	Price        float64 `json:"price"`
	Currency     string  `json:"currency"`
}
