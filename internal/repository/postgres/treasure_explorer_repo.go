package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
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
		INSERT INTO treasure_explorer (chest_id, player_id, score, source)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		explorer.ChestID,
		explorer.PlayerID,
		explorer.Score,
		explorer.Source,
	).Scan(&explorer.ID, &explorer.CreatedAt)

	if err != nil {
		// Check for unique constraint violation on (chest_id, player_id)
		if pqErr, ok := err.(*pq.Error); ok {
			// PostgreSQL error code 23505 is unique_violation
			if pqErr.Code == "23505" {
				return &repository.DuplicateExplorerError{
					ChestID:  explorer.ChestID,
					PlayerID: explorer.PlayerID,
				}
			}
		}
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
		SELECT id, chest_id, player_id, score, source, created_at
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
			&explorer.Source,
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

// GetLeaderboard retrieves the leaderboard for a specific source
// Returns leaderboard entries with rank, name, and total points
func (r *treasureExplorerRepository) GetLeaderboard(
	ctx context.Context,
	source string,
	limit, offset int,
) ([]models.LeaderboardEntry, int, error) {
	// First, get total count of players for this source
	countQuery := `
		SELECT COUNT(DISTINCT player_id)
		FROM treasure_explorer
		WHERE source = $1
	`

	var total int
	err := r.db.QueryRowContext(ctx, countQuery, source).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count players: %w", err)
	}

	// Optimized leaderboard query (no window function)
	// Tiebreaker: MAX(created_at) DESC - later time means slower, ranks lower
	query := `
		SELECT 
			p.name,
			SUM(te.score) as total_points,
			MAX(te.created_at) as latest_time
		FROM treasure_explorer te
		JOIN players p ON te.player_id = p.id
		WHERE te.source = $1
		GROUP BY te.player_id, p.name
		ORDER BY total_points DESC, latest_time DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, source, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get leaderboard: %w", err)
	}
	defer rows.Close()

	var entries []models.LeaderboardEntry
	rank := offset + 1 // Calculate rank from offset + position

	for rows.Next() {
		var entry models.LeaderboardEntry
		var latestTime string // We scan but don't return it

		err := rows.Scan(
			&entry.Name,
			&entry.Points,
			&latestTime,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan leaderboard entry: %w", err)
		}

		entry.Rank = rank
		rank++

		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating leaderboard rows: %w", err)
	}

	return entries, total, nil
}
