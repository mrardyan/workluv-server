-- +goose Up
-- +goose StatementBegin

-- Create workspaces table for workspace service
-- All date fields use epoch seconds (BIGINT) to align with shared.Time type
CREATE TABLE IF NOT EXISTS workspaces (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    owner_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    is_active BOOLEAN DEFAULT true,
    created_at BIGINT NOT NULL DEFAULT EXTRACT(epoch FROM NOW()), -- epoch seconds
    updated_at BIGINT NOT NULL DEFAULT EXTRACT(epoch FROM NOW())  -- epoch seconds
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_workspaces_owner ON workspaces(owner_id);
CREATE INDEX IF NOT EXISTS idx_workspaces_active ON workspaces(is_active);
CREATE INDEX IF NOT EXISTS idx_workspaces_created_at ON workspaces(created_at);
CREATE INDEX IF NOT EXISTS idx_workspaces_updated_at ON workspaces(updated_at);

-- Create updated_at trigger function for epoch-based updates
CREATE OR REPLACE FUNCTION update_workspaces_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = EXTRACT(epoch FROM NOW()); -- Set to current epoch seconds
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_workspaces_updated_at BEFORE UPDATE ON workspaces
    FOR EACH ROW EXECUTE FUNCTION update_workspaces_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop trigger first
DROP TRIGGER IF EXISTS update_workspaces_updated_at ON workspaces;

-- Drop function
DROP FUNCTION IF EXISTS update_workspaces_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_workspaces_updated_at;
DROP INDEX IF EXISTS idx_workspaces_created_at;
DROP INDEX IF EXISTS idx_workspaces_active;
DROP INDEX IF EXISTS idx_workspaces_owner;

-- Drop workspaces table
DROP TABLE IF EXISTS workspaces;

-- +goose StatementEnd
