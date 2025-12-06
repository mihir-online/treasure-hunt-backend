package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/mihir/treasure-hunt-backend/config"
	"github.com/mihir/treasure-hunt-backend/internal/handlers"
	"github.com/mihir/treasure-hunt-backend/internal/qrcode"
	"github.com/mihir/treasure-hunt-backend/internal/repository/postgres"
	"github.com/mihir/treasure-hunt-backend/internal/service"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database connection
	db, err := postgres.NewDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("✓ Database connection established")

	// Initialize repositories
	chestRepo := postgres.NewTreasureChestRepository(db)
	explorerRepo := postgres.NewTreasureExplorerRepository(db)
	ownerRepo := postgres.NewTreasureOwnerRepository(db)
	playerRepo := postgres.NewPlayerRepository(db)

	// Initialize QR code generator
	qrGenerator := qrcode.NewGenerator(&cfg.QRCode)

	// Initialize services
	treasureService := service.NewTreasureService(
		chestRepo,
		explorerRepo,
		ownerRepo,
		playerRepo,
		qrGenerator,
	)

	// Initialize handlers
	treasureHandler := handlers.NewTreasureHandler(treasureService)

	// Setup router
	router := handlers.NewRouter(treasureHandler)

	// Create server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 Starting server on %s:%s", cfg.Server.Host, cfg.Server.Port)
		log.Printf("📊 Status endpoint: http://%s:%s/status", cfg.Server.Host, cfg.Server.Port)
		log.Printf("🔌 API endpoints: http://%s:%s/api/v1", cfg.Server.Host, cfg.Server.Port)
		log.Printf("📂 QR codes will be stored in: %s", cfg.QRCode.StorageDir)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
