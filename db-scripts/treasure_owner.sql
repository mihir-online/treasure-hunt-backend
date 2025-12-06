-- Drop table if exists (optional, comment out in production)
DROP TABLE IF EXISTS treasure_owner;

-- Create treasure_owner table
CREATE TABLE treasure_owner (
    id SERIAL PRIMARY KEY,
    chest_id UUID NOT NULL UNIQUE,
    player_id UUID NOT NULL,
    source VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


-- Add comments for documentation
COMMENT ON TABLE treasure_owner IS 'Stores owner information for treasure chests';
COMMENT ON COLUMN treasure_owner.id IS 'Serial primary key identifier';
COMMENT ON COLUMN treasure_owner.chest_id IS 'Unique reference to treasure chest (one owner per chest)';
COMMENT ON COLUMN treasure_owner.player_id IS 'UUID reference to player who owns this chest';
COMMENT ON COLUMN treasure_owner.source IS 'Source/category of the treasure hunt';
COMMENT ON COLUMN treasure_owner.created_at IS 'Timestamp when record was created';

