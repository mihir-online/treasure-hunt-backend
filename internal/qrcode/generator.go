package qrcode

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/mihir/treasure-hunt-backend/config"
	qrcode "github.com/skip2/go-qrcode"
)

// Generator handles QR code generation
type Generator struct {
	config *config.QRCodeConfig
}

// NewGenerator creates a new QR code generator
func NewGenerator(cfg *config.QRCodeConfig) *Generator {
	return &Generator{
		config: cfg,
	}
}

// GenerateURL creates the URL for a treasure chest based on security settings
func (g *Generator) GenerateURL(chestID, source, qrMatch string, isSecure bool) string {
	baseURL := g.config.BaseURL + g.config.QRPath

	params := url.Values{}
	params.Add("id", chestID)
	params.Add("source", source)

	if isSecure {
		params.Add("qr_match", qrMatch)
	}

	return baseURL + "?" + params.Encode()
}

// GenerateQRCode generates a QR code image and saves it to the storage directory
func (g *Generator) GenerateQRCode(url, source, qrMatch string) (string, error) {
	// Ensure storage directory exists
	if err := os.MkdirAll(g.config.StorageDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Create filename: SOURCE_QRMATCH.png
	// Replace spaces and special characters in source with underscores
	sanitizedSource := sanitizeFilename(source)
	filename := fmt.Sprintf("%s_%s.png", sanitizedSource, qrMatch)
	filepath := filepath.Join(g.config.StorageDir, filename)

	// Generate QR code
	err := qrcode.WriteFile(url, qrcode.Medium, g.config.Size, filepath)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	return filepath, nil
}

// GenerateBatch generates multiple QR codes for a batch of chests
func (g *Generator) GenerateBatch(chests []ChestInfo) ([]QRCodeResult, error) {
	results := make([]QRCodeResult, 0, len(chests))

	for _, chest := range chests {
		url := g.GenerateURL(chest.ID, chest.Source, chest.QRMatch, chest.IsSecure)

		filepath, err := g.GenerateQRCode(url, chest.Source, chest.QRMatch)
		if err != nil {
			return nil, fmt.Errorf("failed to generate QR code for chest %s: %w", chest.ID, err)
		}

		results = append(results, QRCodeResult{
			ChestID:  chest.ID,
			QRMatch:  chest.QRMatch,
			URL:      url,
			FilePath: filepath,
		})
	}

	return results, nil
}

// ChestInfo contains the information needed to generate a QR code
type ChestInfo struct {
	ID       string
	Source   string
	QRMatch  string
	IsSecure bool
}

// QRCodeResult contains the result of QR code generation
type QRCodeResult struct {
	ChestID  string
	QRMatch  string
	URL      string
	FilePath string
}

// sanitizeFilename removes or replaces characters that are problematic in filenames
func sanitizeFilename(s string) string {
	// Replace spaces with underscores
	s = strings.ReplaceAll(s, " ", "_")

	// Remove or replace other problematic characters
	problematic := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range problematic {
		s = strings.ReplaceAll(s, char, "_")
	}

	// Limit length to avoid filesystem issues
	if len(s) > 50 {
		s = s[:50]
	}

	return s
}
