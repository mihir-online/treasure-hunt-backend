package handlers

import (
	"github.com/gorilla/mux"
)

// NewRouter creates and configures the application router
func NewRouter(treasureHandler *TreasureHandler) *mux.Router {
	router := mux.NewRouter()

	// Apply global middleware
	router.Use(RecoveryMiddleware)
	router.Use(LoggingMiddleware)
	router.Use(CORSMiddleware)

	// Health check endpoint
	router.HandleFunc("/status", treasureHandler.GetStatus).Methods("GET")

	// API v1 routes
	v1 := router.PathPrefix("/api/v1").Subrouter()

	// QR Code endpoints
	v1.HandleFunc("/qr-codes", treasureHandler.CreateQRCode).Methods("POST")

	// Chest endpoints - Note: The actual path will be /treasure/claim from base URL
	router.HandleFunc("/treasure/claim", treasureHandler.ClaimChest).Methods("POST")

	return router
}
