-- Drop table if exists (optional, comment out in production)
DROP TABLE IF EXISTS treasure_explorer;

-- Create treasure_explorer table
CREATE TABLE treasure_explorer (
    id SERIAL PRIMARY KEY,
    chest_id UUID NOT NULL,
    player_id INTEGER NOT NULL,
    score INTEGER NOT NULL,
    source VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (chest_id, player_id)
);

-- Create indexes on commonly queried columns
CREATE INDEX idx_treasure_explorer_chest_id ON treasure_explorer(chest_id);
CREATE INDEX idx_treasure_explorer_player_id ON treasure_explorer(player_id);
CREATE INDEX idx_treasure_explorer_source ON treasure_explorer(source);
CREATE INDEX idx_treasure_explorer_created_at ON treasure_explorer(created_at);

-- Composite index for optimized leaderboard queries
-- Covers: source filter, player_id grouping, score & created_at aggregation
CREATE INDEX idx_explorer_source_player_score_time 
ON treasure_explorer(source, player_id, score, created_at);

-- Add comments for documentation
COMMENT ON TABLE treasure_explorer IS 'Stores explorer information for treasure chests (multiple explorers per chest allowed)';
COMMENT ON COLUMN treasure_explorer.id IS 'Serial primary key identifier';
COMMENT ON COLUMN treasure_explorer.chest_id IS 'Reference to treasure chest (multiple explorers can have same chest_id)';
COMMENT ON COLUMN treasure_explorer.player_id IS 'Reference to player who explored this chest';
COMMENT ON COLUMN treasure_explorer.score IS 'Score earned by the player for this chest';
COMMENT ON COLUMN treasure_explorer.source IS 'Source/category of the treasure hunt';
COMMENT ON COLUMN treasure_explorer.created_at IS 'Timestamp when record was created';

