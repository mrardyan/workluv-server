-- +goose Up
-- +goose StatementBegin

-- Create sessions table for user session management
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    refresh_token VARCHAR(255) UNIQUE NOT NULL,
    expires_at BIGINT NOT NULL, -- epoch seconds using shared.Time
    revoked_at BIGINT NULL, -- null = active, value = revoked timestamp
    created_at BIGINT NOT NULL DEFAULT EXTRACT(epoch FROM NOW()),
    updated_at BIGINT NOT NULL DEFAULT EXTRACT(epoch FROM NOW())
);

-- Create index for faster lookup by account_id
CREATE INDEX idx_sessions_account_id ON sessions(account_id);

-- Create index for refresh token lookup
CREATE INDEX idx_sessions_refresh_token ON sessions(refresh_token);

-- Create index for active sessions (where revoked_at is null)
CREATE INDEX idx_sessions_active ON sessions(account_id) WHERE revoked_at IS NULL;

-- Create index for expired sessions cleanup
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop indexes first
DROP INDEX IF EXISTS idx_sessions_expires_at;
DROP INDEX IF EXISTS idx_sessions_active;
DROP INDEX IF EXISTS idx_sessions_refresh_token;
DROP INDEX IF EXISTS idx_sessions_account_id;

-- Drop sessions table
DROP TABLE IF EXISTS sessions;

-- +goose StatementEnd
