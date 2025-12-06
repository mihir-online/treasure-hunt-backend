package models

// CreateQRCodeRequest represents the request to create a new QR code
type CreateQRCodeRequest struct {
	Count     int    `json:"count"      validate:"required,min=1,max=100"` // Number of QR codes to generate
	Source    string `json:"source"     validate:"required"`               // Source of the treasure chest (must be unique)
	IsSecure  bool   `json:"is_secure"`                                    // Whether the chest is secure (default: false)
	CreatedBy string `json:"created_by" validate:"required"`               // Username or ID of creator
}

// ClaimChestRequest represents the request to claim a chest
type ClaimChestRequest struct {
	ChestID string `json:"chest_id" validate:"required"` // UUID of the chest
	Source  string `json:"source"   validate:"required"` // Source to verify
	QRMatch string `json:"qr_match"`                     // Only required for secure chests
	Name    string `json:"name"     validate:"required"` // Claimer's name
	Email   string `json:"email"    validate:"required"` // Claimer's email
	Phone   string `json:"phone"    validate:"required"` // Claimer's phone
}

// LeaderboardRequest represents the request to fetch leaderboard
type LeaderboardRequest struct {
	Source string `json:"source" validate:"required"` // Source/category of the treasure hunt
	Limit  int    `json:"limit"  validate:"required,min=1,max=100"` // Number of entries to return
	Offset int    `json:"offset" validate:"min=0"`    // Offset for pagination
}