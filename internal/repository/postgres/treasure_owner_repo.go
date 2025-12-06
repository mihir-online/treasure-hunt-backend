package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/mihir/treasure-hunt-backend/internal/models"
	"github.com/mihir/treasure-hunt-backend/internal/repository"
)

// treasureOwnerRepository implements the TreasureOwnerRepository interface
type treasureOwnerRepository struct {
	db *sql.DB
}

// NewTreasureOwnerRepository creates a new PostgreSQL treasure owner repository
func NewTreasureOwnerRepository(db *sql.DB) repository.TreasureOwnerRepository {
	return &treasureOwnerRepository{
		db: db,
	}
}

// Create creates a new owner entry (one owner per chest, unique constraint)
func (r *treasureOwnerRepository) Create(ctx context.Context, owner *models.TreasureOwner) error {
	query := `
		INSERT INTO treasure_owner (chest_id, player_id, source)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		owner.ChestID,
		owner.PlayerID,
		owner.Source,
	).Scan(&owner.ID, &owner.CreatedAt)

	if err != nil {
		// Check for unique constraint violation on chest_id
		if pqErr, ok := err.(*pq.Error); ok {
			// PostgreSQL error code 23505 is unique_violation
			if pqErr.Code == "23505" {
				return &repository.DuplicateChestError{ChestID: owner.ChestID}
			}
		}
		return fmt.Errorf("failed to create treasure owner: %w", err)
	}

	return nil
}

// GetByChestID retrieves the owner for a specific chest
func (r *treasureOwnerRepository) GetByChestID(
	ctx context.Context,
	chestID string,
) (*models.TreasureOwner, error) {
	query := `
		SELECT id, chest_id, player_id, source, created_at
		FROM treasure_owner
		WHERE chest_id = $1
	`

	owner := &models.TreasureOwner{}
	err := r.db.QueryRowContext(ctx, query, chestID).Scan(
		&owner.ID,
		&owner.ChestID,
		&owner.PlayerID,
		&owner.Source,
		&owner.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("owner for chest '%s' not found", chestID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get owner by chest ID: %w", err)
	}

	return owner, nil
}
