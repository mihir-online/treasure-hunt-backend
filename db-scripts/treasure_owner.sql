-- Drop table if exists (optional, comment out in production)
DROP TABLE IF EXISTS treasure_owner;

-- Create treasure_owner table
CREATE TABLE treasure_owner (
    id SERIAL PRIMARY KEY,
    chest_id UUID NOT NULL UNIQUE,
    name VARCHAR(200) NOT NULL,
    email VARCHAR(320) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes on commonly queried columns
CREATE INDEX idx_treasure_owner_chest_id ON treasure_owner(chest_id);
CREATE INDEX idx_treasure_owner_email ON treasure_owner(email);
CREATE INDEX idx_treasure_owner_created_at ON treasure_owner(created_at);

-- Add comments for documentation
COMMENT ON TABLE treasure_owner IS 'Stores owner information for treasure chests';
COMMENT ON COLUMN treasure_owner.id IS 'Serial primary key identifier';
COMMENT ON COLUMN treasure_owner.chest_id IS 'Unique reference to treasure chest (one owner per chest)';
COMMENT ON COLUMN treasure_owner.name IS 'Owner name';
COMMENT ON COLUMN treasure_owner.email IS 'Owner email address';
COMMENT ON COLUMN treasure_owner.phone_number IS 'Owner phone number';
COMMENT ON COLUMN treasure_owner.created_at IS 'Timestamp when record was created';

