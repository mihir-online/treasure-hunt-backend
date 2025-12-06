package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mihir/treasure-hunt-backend/internal/models"
	"github.com/mihir/treasure-hunt-backend/internal/service"
)

// TreasureHandler handles HTTP requests for treasure operations
type TreasureHandler struct {
	service service.TreasureService
}

// NewTreasureHandler creates a new instance of TreasureHandler
func NewTreasureHandler(service service.TreasureService) *TreasureHandler {
	return &TreasureHandler{
		service: service,
	}
}

// GetStatus handles the GET /status endpoint
func (h *TreasureHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	response := models.StatusResponse{
		Status:  "healthy",
		Message: "Treasure Hunt API is running",
		Version: "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// CreateQRCode handles the POST /api/v1/qr-codes endpoint
func (h *TreasureHandler) CreateQRCode(w http.ResponseWriter, r *http.Request) {
	var req models.CreateQRCodeRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call service layer (business logic will be implemented later)
	response, err := h.service.CreateQRCode(r.Context(), &req)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// ClaimChest handles the POST /treasure/claim endpoint
func (h *TreasureHandler) ClaimChest(w http.ResponseWriter, r *http.Request) {
	var req models.ClaimChestRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call service layer
	response, err := h.service.ClaimChest(r.Context(), &req)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if it's a duplicate owner scenario (409 Conflict)
	if !response.Success && response.Message == "This chest already has an owner" {
		h.sendConflict(w, response.Message)
		return
	}

	// Return appropriate status code
	statusCode := http.StatusOK
	if !response.Success {
		statusCode = http.StatusBadRequest
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// sendError sends an error response
func (h *TreasureHandler) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{
		Success: false,
		Error:   message,
	})
}

// sendConflict sends a 409 conflict response
func (h *TreasureHandler) sendConflict(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusConflict)
	json.NewEncoder(w).Encode(models.ClaimChestResponse{
		Success: false,
		Message: message,
	})
}

// GetLeaderboard handles the POST /leaderboard endpoint
func (h *TreasureHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	var req models.LeaderboardRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set default values if not provided
	if req.Limit == 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// Call service layer
	response, err := h.service.GetLeaderboard(r.Context(), &req)
	if err != nil {
		h.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
