-- +goose Up
-- +goose StatementBegin

-- Create accounts table for account service
-- All date fields use epoch seconds (BIGINT) to align with shared.Time type
CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(200),
    is_active BOOLEAN DEFAULT true,
    email_verified BOOLEAN DEFAULT false,
    created_at BIGINT NOT NULL DEFAULT EXTRACT(epoch FROM NOW()), -- epoch seconds
    updated_at BIGINT NOT NULL DEFAULT EXTRACT(epoch FROM NOW())  -- epoch seconds
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_accounts_email ON accounts(email);
CREATE INDEX IF NOT EXISTS idx_accounts_active ON accounts(is_active);
CREATE INDEX IF NOT EXISTS idx_accounts_created_at ON accounts(created_at);
CREATE INDEX IF NOT EXISTS idx_accounts_updated_at ON accounts(updated_at);

-- Create updated_at trigger function for epoch-based updates
CREATE OR REPLACE FUNCTION update_accounts_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = EXTRACT(epoch FROM NOW()); -- Set to current epoch seconds
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_accounts_updated_at BEFORE UPDATE ON accounts
    FOR EACH ROW EXECUTE FUNCTION update_accounts_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop trigger first
DROP TRIGGER IF EXISTS update_accounts_updated_at ON accounts;

-- Drop function
DROP FUNCTION IF EXISTS update_accounts_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_accounts_updated_at;
DROP INDEX IF EXISTS idx_accounts_created_at;
DROP INDEX IF EXISTS idx_accounts_active;
DROP INDEX IF EXISTS idx_accounts_email;

-- Drop accounts table
DROP TABLE IF EXISTS accounts;

-- +goose StatementEnd
