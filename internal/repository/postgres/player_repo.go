package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/mihir/treasure-hunt-backend/internal/models"
	"github.com/mihir/treasure-hunt-backend/internal/repository"
)

// playerRepository implements the PlayerRepository interface
type playerRepository struct {
	db *sql.DB
}

// NewPlayerRepository creates a new PostgreSQL player repository
func NewPlayerRepository(db *sql.DB) repository.PlayerRepository {
	return &playerRepository{
		db: db,
	}
}

// GetByEmail retrieves a player by email
func (r *playerRepository) GetByEmail(ctx context.Context, email string) (*models.Player, error) {
	query := `
		SELECT id, email, name, phone_number, created_at
		FROM players
		WHERE email = $1
	`

	player := &models.Player{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&player.ID,
		&player.Email,
		&player.Name,
		&player.PhoneNumber,
		&player.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Return nil if player doesn't exist
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get player by email: %w", err)
	}

	return player, nil
}

// Create creates a new player
func (r *playerRepository) Create(ctx context.Context, player *models.Player) error {
	query := `
		INSERT INTO players (email, name, phone_number)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		player.Email,
		player.Name,
		player.PhoneNumber,
	).Scan(&player.ID, &player.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create player: %w", err)
	}

	return nil
}

// GetOrCreate retrieves a player by email or creates one if it doesn't exist
// Uses INSERT ... ON CONFLICT for atomic upsert to handle concurrent requests
func (r *playerRepository) GetOrCreate(
	ctx context.Context,
	email, name, phoneNumber string,
) (*models.Player, error) {
	// Use INSERT ... ON CONFLICT for atomic upsert
	// This handles the race condition where multiple concurrent requests
	// try to create the same player
	query := `
		INSERT INTO players (email, name, phone_number)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) DO UPDATE SET
			name = EXCLUDED.name,
			phone_number = EXCLUDED.phone_number
		RETURNING id, email, name, phone_number, created_at
	`

	player := &models.Player{}
	err := r.db.QueryRowContext(
		ctx,
		query,
		email,
		name,
		phoneNumber,
	).Scan(
		&player.ID,
		&player.Email,
		&player.Name,
		&player.PhoneNumber,
		&player.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get or create player: %w", err)
	}

	return player, nil
}
