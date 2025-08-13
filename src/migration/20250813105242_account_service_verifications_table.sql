-- +goose Up
-- +goose StatementBegin

-- Create verifications table for account service
-- All date fields use epoch seconds (BIGINT) to align with shared.Time type
CREATE TABLE IF NOT EXISTS verifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    token VARCHAR(64) NOT NULL,
    expires_at BIGINT NOT NULL, -- epoch seconds using shared.Time
    completed_at BIGINT NULL, -- epoch seconds, null = not completed
    created_at BIGINT NOT NULL DEFAULT EXTRACT(epoch FROM NOW()), -- epoch seconds
    updated_at BIGINT NOT NULL DEFAULT EXTRACT(epoch FROM NOW())  -- epoch seconds
);

-- Create indexes for better performance
CREATE UNIQUE INDEX IF NOT EXISTS idx_verifications_token ON verifications(token);
CREATE INDEX IF NOT EXISTS idx_verifications_account_id ON verifications(account_id);
CREATE INDEX IF NOT EXISTS idx_verifications_type ON verifications(type);
CREATE INDEX IF NOT EXISTS idx_verifications_status ON verifications(status);
CREATE INDEX IF NOT EXISTS idx_verifications_expires_at ON verifications(expires_at);
CREATE INDEX IF NOT EXISTS idx_verifications_completed_at ON verifications(completed_at);
CREATE INDEX IF NOT EXISTS idx_verifications_created_at ON verifications(created_at);
CREATE INDEX IF NOT EXISTS idx_verifications_updated_at ON verifications(updated_at);

-- Create updated_at trigger function for epoch-based updates
CREATE OR REPLACE FUNCTION update_verifications_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = EXTRACT(epoch FROM NOW()); -- Set to current epoch seconds
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_verifications_updated_at BEFORE UPDATE ON verifications
    FOR EACH ROW EXECUTE FUNCTION update_verifications_updated_at_column();

-- Remove email verification columns from accounts table (if they exist)
-- These columns were added in the previous migration, so we need to remove them
ALTER TABLE accounts 
DROP COLUMN IF EXISTS email_verification_token,
DROP COLUMN IF EXISTS email_token_expires_at;

-- Drop the related indexes if they exist
DROP INDEX IF EXISTS idx_accounts_verification_token;
DROP INDEX IF EXISTS idx_accounts_token_expires_at;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Re-add email verification columns to accounts table
ALTER TABLE accounts 
ADD COLUMN email_verification_token VARCHAR(64),
ADD COLUMN email_token_expires_at BIGINT; -- epoch seconds

-- Recreate indexes
CREATE INDEX IF NOT EXISTS idx_accounts_verification_token ON accounts(email_verification_token);
CREATE INDEX IF NOT EXISTS idx_accounts_token_expires_at ON accounts(email_token_expires_at);

-- Drop trigger first
DROP TRIGGER IF EXISTS update_verifications_updated_at ON verifications;

-- Drop function
DROP FUNCTION IF EXISTS update_verifications_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_verifications_updated_at;
DROP INDEX IF EXISTS idx_verifications_created_at;
DROP INDEX IF EXISTS idx_verifications_completed_at;
DROP INDEX IF EXISTS idx_verifications_expires_at;
DROP INDEX IF EXISTS idx_verifications_status;
DROP INDEX IF EXISTS idx_verifications_type;
DROP INDEX IF EXISTS idx_verifications_account_id;
DROP INDEX IF EXISTS idx_verifications_token;

-- Drop verifications table
DROP TABLE IF EXISTS verifications;

-- +goose StatementEnd
