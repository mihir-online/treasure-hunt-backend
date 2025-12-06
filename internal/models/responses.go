package models

// StatusResponse represents the health check response
type StatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Version string `json:"version"`
}

// QRCodeInfo represents information about a generated QR code
type QRCodeInfo struct {
	ChestID  string `json:"chest_id"`  // UUID of the chest
	QRMatch  string `json:"qr_match"`  // QR match code (e.g., "1-X7K9P2M4")
	URL      string `json:"url"`       // Full URL for the QR code
	FilePath string `json:"file_path"` // Path to the generated QR code image
}

// CreateQRCodeResponse represents the response after creating a QR code
type CreateQRCodeResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Source  string       `json:"source"`
	Count   int          `json:"count"`
	QRCodes []QRCodeInfo `json:"qr_codes"`
}

// ClaimChestResponse represents the response after claiming a chest
type ClaimChestResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ChestID string `json:"chest_id,omitempty"`
	Status  string `json:"status,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
}

// LeaderboardEntry represents a single entry in the leaderboard
type LeaderboardEntry struct {
	Name   string `json:"name"`   // Player's name
	Rank   int    `json:"rank"`   // Player's rank (1-indexed)
	Points int    `json:"points"` // Total points earned
}

// LeaderboardResponse represents the response for leaderboard request
type LeaderboardResponse struct {
	Leaderboard []LeaderboardEntry `json:"leaderboard"` // List of leaderboard entries
	Total       int                `json:"total"`       // Total number of players in this source
	Source      string             `json:"source"`      // Source/category of the treasure hunt
}