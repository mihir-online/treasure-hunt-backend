package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mihir/treasure-hunt-backend/internal/models"
	"github.com/mihir/treasure-hunt-backend/internal/repository"
)

// treasureExplorerRepository implements the TreasureExplorerRepository interface
type treasureExplorerRepository struct {
	db *sql.DB
}

// NewTreasureExplorerRepository creates a new PostgreSQL treasure explorer repository
func NewTreasureExplorerRepository(db *sql.DB) repository.TreasureExplorerRepository {
	return &treasureExplorerRepository{
		db: db,
	}
}

// Create creates a new explorer entry (multiple explorers allowed per chest)
func (r *treasureExplorerRepository) Create(
	ctx context.Context,
	explorer *models.TreasureExplorer,
) error {
	query := `
		INSERT INTO treasure_explorer (chest_id, player_id, score)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		explorer.ChestID,
		explorer.PlayerID,
		explorer.Score,
	).Scan(&explorer.ID, &explorer.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create treasure explorer: %w", err)
	}

	return nil
}

// GetByChestID retrieves all explorers for a specific chest
func (r *treasureExplorerRepository) GetByChestID(
	ctx context.Context,
	chestID string,
) ([]*models.TreasureExplorer, error) {
	query := `
		SELECT id, chest_id, player_id, score, created_at
		FROM treasure_explorer
		WHERE chest_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, chestID)
	if err != nil {
		return nil, fmt.Errorf("failed to get explorers by chest ID: %w", err)
	}
	defer rows.Close()

	var explorers []*models.TreasureExplorer
	for rows.Next() {
		explorer := &models.TreasureExplorer{}
		err := rows.Scan(
			&explorer.ID,
			&explorer.ChestID,
			&explorer.PlayerID,
			&explorer.Score,
			&explorer.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan explorer: %w", err)
		}
		explorers = append(explorers, explorer)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating explorer rows: %w", err)
	}

	return explorers, nil
}
