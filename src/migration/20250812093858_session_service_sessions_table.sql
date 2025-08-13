-- +goose Up
-- +goose StatementBegin

-- Create sessions table for session service
-- All date fields use epoch seconds (BIGINT) to align with shared.Time type
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    refresh_token VARCHAR(255) UNIQUE NOT NULL,
    expires_at BIGINT NOT NULL, -- epoch seconds using shared.Time
    revoked_at BIGINT NULL, -- null = active, value = revoked timestamp (epoch seconds)
    created_at BIGINT NOT NULL DEFAULT EXTRACT(epoch FROM NOW()), -- epoch seconds
    updated_at BIGINT NOT NULL DEFAULT EXTRACT(epoch FROM NOW())  -- epoch seconds
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_sessions_account_id ON sessions(account_id);
CREATE INDEX IF NOT EXISTS idx_sessions_refresh_token ON sessions(refresh_token);
CREATE INDEX IF NOT EXISTS idx_sessions_active ON sessions(account_id) WHERE revoked_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
CREATE INDEX IF NOT EXISTS idx_sessions_created_at ON sessions(created_at);
CREATE INDEX IF NOT EXISTS idx_sessions_updated_at ON sessions(updated_at);

-- Create updated_at trigger function for epoch-based updates
CREATE OR REPLACE FUNCTION update_sessions_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = EXTRACT(epoch FROM NOW()); -- Set to current epoch seconds
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_sessions_updated_at BEFORE UPDATE ON sessions
    FOR EACH ROW EXECUTE FUNCTION update_sessions_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop trigger first
DROP TRIGGER IF EXISTS update_sessions_updated_at ON sessions;

-- Drop function
DROP FUNCTION IF EXISTS update_sessions_updated_at_column();

-- Drop indexes first
DROP INDEX IF EXISTS idx_sessions_updated_at;
DROP INDEX IF EXISTS idx_sessions_created_at;
DROP INDEX IF EXISTS idx_sessions_expires_at;
DROP INDEX IF EXISTS idx_sessions_active;
DROP INDEX IF EXISTS idx_sessions_refresh_token;
DROP INDEX IF EXISTS idx_sessions_account_id;

-- Drop sessions table
DROP TABLE IF EXISTS sessions;

-- +goose StatementEnd
