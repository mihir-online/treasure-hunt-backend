-- Drop table if exists (optional, comment out in production)
DROP TABLE IF EXISTS treasure_chest;

-- Drop sequence if exists (optional, comment out in production)
DROP SEQUENCE IF EXISTS qr_code_sequence;

-- Create sequence for QR code numbering
CREATE SEQUENCE qr_code_sequence START 1;

-- Create function to generate random 8 character string
CREATE OR REPLACE FUNCTION generate_random_string(length INTEGER)
RETURNS TEXT AS $$
DECLARE
    chars TEXT := 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
    result TEXT := '';
    i INTEGER := 0;
BEGIN
    FOR i IN 1..length LOOP
        result := result || substr(chars, floor(random() * length(chars) + 1)::integer, 1);
    END LOOP;
    RETURN result;
END;
$$ LANGUAGE plpgsql;

-- Create function to generate QR match code
CREATE OR REPLACE FUNCTION generate_qr_match()
RETURNS TEXT AS $$
BEGIN
    RETURN nextval('qr_code_sequence')::TEXT || '-' || generate_random_string(8);
END;
$$ LANGUAGE plpgsql;

-- Create treasure_chest table
CREATE TABLE treasure_chest (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source VARCHAR(50) NOT NULL,
    qr_match VARCHAR(20) NOT NULL UNIQUE DEFAULT generate_qr_match(),
    is_secure BOOLEAN NOT NULL DEFAULT FALSE,
    status VARCHAR(20) NOT NULL DEFAULT 'UNCLAIMED',
    first_score INTEGER NOT NULL DEFAULT 100,
    subsequent_score INTEGER NOT NULL DEFAULT 0,
    created_by VARCHAR(200) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes on commonly queried columns
CREATE INDEX idx_treasure_chest_status ON treasure_chest(status);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_treasure_chest_updated_at
    BEFORE UPDATE ON treasure_chest
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Add comments for documentation
COMMENT ON TABLE treasure_chest IS 'Stores treasure chest information for the treasure hunt application';
COMMENT ON COLUMN treasure_chest.id IS 'Unique identifier for the treasure chest';
COMMENT ON COLUMN treasure_chest.source IS 'Source of the treasure chest';
COMMENT ON COLUMN treasure_chest.qr_match IS 'QR code match identifier - auto-generated as {sequence}-{random_8_chars}';
COMMENT ON COLUMN treasure_chest.is_secure IS 'Indicates if the treasure chest is secure';
COMMENT ON COLUMN treasure_chest.status IS 'Current status of the treasure chest';
COMMENT ON COLUMN treasure_chest.first_score IS 'Score awarded to the first finder (owner)';
COMMENT ON COLUMN treasure_chest.subsequent_score IS 'Score awarded to subsequent finders (explorers)';
COMMENT ON COLUMN treasure_chest.created_by IS 'User who created this record';
COMMENT ON COLUMN treasure_chest.created_at IS 'Timestamp when record was created';
COMMENT ON COLUMN treasure_chest.updated_at IS 'Timestamp when record was last updated';

