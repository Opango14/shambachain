package models

// RegisterBatchResponse represents the response body for batch registration
type RegisterBatchResponse struct {
	BatchID     string `json:"batch_id"`
	QRCodeURL   string `json:"qr_code_url"`
	QRCodeData  string `json:"qr_code_data"` // base64 encoded image
	GenesisHash string `json:"genesis_hash"`
}

// TraceabilityResponse represents the response body for traceability queries
type TraceabilityResponse struct {
	Batch      Batch   `json:"batch"`
	Blockchain []Block `json:"blockchain"`
	Verified   bool    `json:"verified"`
	ChainValid bool    `json:"chain_valid"`
}

// ErrorResponse represents a generic error response
type ErrorResponse struct {
	Error string `json:"error"`
}
