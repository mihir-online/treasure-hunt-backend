-- Drop table if exists (optional, comment out in production)
DROP TABLE IF EXISTS players;

-- Create players table
CREATE TABLE players (
    id SERIAL PRIMARY KEY,
    email VARCHAR(320) NOT NULL UNIQUE,
    name VARCHAR(200) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes on commonly queried columns
CREATE INDEX idx_players_email ON players(email);
CREATE INDEX idx_players_id ON players(id);

-- Add comments for documentation
COMMENT ON TABLE players IS 'Stores player information for the treasure hunt application';
COMMENT ON COLUMN players.id IS 'Serial primary key identifier';
COMMENT ON COLUMN players.email IS 'Player email address (unique)';
COMMENT ON COLUMN players.name IS 'Player name';
COMMENT ON COLUMN players.phone_number IS 'Player phone number';
COMMENT ON COLUMN players.created_at IS 'Timestamp when record was created';


