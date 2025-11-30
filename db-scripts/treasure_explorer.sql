-- Drop table if exists (optional, comment out in production)
DROP TABLE IF EXISTS treasure_explorer;

-- Create treasure_explorer table
CREATE TABLE treasure_explorer (
    id SERIAL PRIMARY KEY,
    chest_id UUID NOT NULL,
    name VARCHAR(200) NOT NULL,
    email VARCHAR(320) NOT NULL,
    phone_number NUMERIC NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes on commonly queried columns
CREATE INDEX idx_treasure_explorer_chest_id ON treasure_explorer(chest_id);
CREATE INDEX idx_treasure_explorer_email ON treasure_explorer(email);
CREATE INDEX idx_treasure_explorer_created_at ON treasure_explorer(created_at);

-- Add comments for documentation
COMMENT ON TABLE treasure_explorer IS 'Stores explorer information for treasure chests (multiple explorers per chest allowed)';
COMMENT ON COLUMN treasure_explorer.id IS 'Serial primary key identifier';
COMMENT ON COLUMN treasure_explorer.chest_id IS 'Reference to treasure chest (multiple explorers can have same chest_id)';
COMMENT ON COLUMN treasure_explorer.name IS 'Explorer name';
COMMENT ON COLUMN treasure_explorer.email IS 'Explorer email address';
COMMENT ON COLUMN treasure_explorer.phone_number IS 'Explorer phone number';
COMMENT ON COLUMN treasure_explorer.created_at IS 'Timestamp when record was created';

