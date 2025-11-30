package service

import (
	"context"
	"fmt"

	"github.com/mihir/treasure-hunt-backend/internal/models"
	"github.com/mihir/treasure-hunt-backend/internal/qrcode"
	"github.com/mihir/treasure-hunt-backend/internal/repository"
)

// TreasureService defines the business logic interface for treasure operations
type TreasureService interface {
	// CreateQRCode creates a new treasure chest with QR code
	CreateQRCode(
		ctx context.Context,
		req *models.CreateQRCodeRequest,
	) (*models.CreateQRCodeResponse, error)

	// ClaimChest processes a chest claim by an explorer
	ClaimChest(
		ctx context.Context,
		req *models.ClaimChestRequest,
	) (*models.ClaimChestResponse, error)

	// GetChestByID retrieves a chest by its ID
	GetChestByID(ctx context.Context, chestID string) (*models.TreasureChest, error)
}

// treasureService implements the TreasureService interface
type treasureService struct {
	chestRepo    repository.TreasureChestRepository
	explorerRepo repository.TreasureExplorerRepository
	ownerRepo    repository.TreasureOwnerRepository
	qrGenerator  *qrcode.Generator
}

// NewTreasureService creates a new instance of TreasureService
func NewTreasureService(
	chestRepo repository.TreasureChestRepository,
	explorerRepo repository.TreasureExplorerRepository,
	ownerRepo repository.TreasureOwnerRepository,
	qrGenerator *qrcode.Generator,
) TreasureService {
	return &treasureService{
		chestRepo:    chestRepo,
		explorerRepo: explorerRepo,
		ownerRepo:    ownerRepo,
		qrGenerator:  qrGenerator,
	}
}

// CreateQRCode creates multiple treasure chests with QR codes
func (s *treasureService) CreateQRCode(
	ctx context.Context,
	req *models.CreateQRCodeRequest,
) (*models.CreateQRCodeResponse, error) {
	// Step 1: Ensure source is unique
	exists, err := s.chestRepo.CheckSourceExists(ctx, req.Source)
	if err != nil {
		return nil, fmt.Errorf("failed to check source uniqueness: %w", err)
	}
	if exists {
		return &models.CreateQRCodeResponse{
			Success: false,
			Message: fmt.Sprintf(
				"Source '%s' already exists. Please use a unique source.",
				req.Source,
			),
		}, nil
	}

	// Step 2: Insert the specified count of rows in the chest table
	chests, err := s.chestRepo.CreateBatch(
		ctx,
		req.Count,
		req.Source,
		req.IsSecure,
		req.CreatedBy,
		"UNCLAIMED",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create treasure chests: %w", err)
	}

	// Step 3: Verify count
	if len(chests) != req.Count {
		return nil, fmt.Errorf("expected %d chests but created %d", req.Count, len(chests))
	}

	// Step 4: Generate URLs and QR codes
	chestInfos := make([]qrcode.ChestInfo, len(chests))
	for i, chest := range chests {
		chestInfos[i] = qrcode.ChestInfo{
			ID:       chest.ID,
			Source:   chest.Source,
			QRMatch:  chest.QRMatch,
			IsSecure: chest.IsSecure,
		}
	}

	// Step 5: Generate QR code images
	qrResults, err := s.qrGenerator.GenerateBatch(chestInfos)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR codes: %w", err)
	}

	// Step 6: Build response
	qrCodes := make([]models.QRCodeInfo, len(qrResults))
	for i, result := range qrResults {
		qrCodes[i] = models.QRCodeInfo{
			ChestID:  result.ChestID,
			QRMatch:  result.QRMatch,
			URL:      result.URL,
			FilePath: result.FilePath,
		}
	}

	return &models.CreateQRCodeResponse{
		Success: true,
		Message: fmt.Sprintf("Successfully created %d treasure chests with QR codes", req.Count),
		Source:  req.Source,
		Count:   len(qrCodes),
		QRCodes: qrCodes,
	}, nil
}

// ClaimChest processes a chest claim by an explorer
func (s *treasureService) ClaimChest(
	ctx context.Context,
	req *models.ClaimChestRequest,
) (*models.ClaimChestResponse, error) {
	// Step 1: Find chest by ID
	chest, err := s.chestRepo.GetByID(ctx, req.ChestID)
	if err != nil {
		return &models.ClaimChestResponse{
			Success: false,
			Message: fmt.Sprintf("Chest with ID '%s' not found", req.ChestID),
		}, nil
	}

	// Step 2: Verify source matches
	if chest.Source != req.Source {
		return &models.ClaimChestResponse{
			Success: false,
			Message: "Source mismatch",
		}, nil
	}

	// Step 3: Check if secure and validate QR match
	if chest.IsSecure {
		if req.QRMatch == "" {
			return &models.ClaimChestResponse{
				Success: false,
				Message: "QR match is required for secure chests",
			}, nil
		}
		if chest.QRMatch != req.QRMatch {
			return &models.ClaimChestResponse{
				Success: false,
				Message: "Invalid QR match code",
			}, nil
		}
	}

	// Step 4: Insert into treasure_explorer table
	explorer := &models.TreasureExplorer{
		ChestID:     req.ChestID,
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.Phone,
	}
	err = s.explorerRepo.Create(ctx, explorer)
	if err != nil {
		return nil, fmt.Errorf("failed to create explorer entry: %w", err)
	}

	// Step 5: Check if already claimed
	if chest.Status == "CLAIMED" {
		return &models.ClaimChestResponse{
			Success: false,
			Message: "This chest already has an owner",
		}, nil
	}

	// Step 6: Mark chest as claimed
	chest.Status = "CLAIMED"
	err = s.chestRepo.Update(ctx, chest)
	if err != nil {
		return nil, fmt.Errorf("failed to update chest status: %w", err)
	}

	// Step 7: Insert into treasure_owner table (may fail with duplicate constraint)
	owner := &models.TreasureOwner{
		ChestID:     req.ChestID,
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.Phone,
	}
	err = s.ownerRepo.Create(ctx, owner)
	if err != nil {
		// Check if it's a duplicate chest_id error (return 409 Conflict)
		if repository.IsDuplicateChestError(err) {
			return &models.ClaimChestResponse{
				Success: false,
				Message: "This chest already has an owner",
			}, nil
		}
		return nil, fmt.Errorf("failed to create owner entry: %w", err)
	}

	return &models.ClaimChestResponse{
		Success: true,
		Message: "Chest claimed successfully!",
		ChestID: chest.ID,
		Status:  chest.Status,
	}, nil
}

// GetChestByID retrieves a chest by its ID
func (s *treasureService) GetChestByID(
	ctx context.Context,
	chestID string,
) (*models.TreasureChest, error) {
	return s.chestRepo.GetByID(ctx, chestID)
}
