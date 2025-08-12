-- +goose Up
-- +goose StatementBegin

-- =============================================================================
-- MIGRATION: [MIGRATION_NAME]
-- =============================================================================
--
-- Description: [Brief description of what this migration does]
--
-- Business Context: [Why is this change needed? Link to ticket/issue]
--
-- Dependencies: [List any other migrations or systems this depends on]
--
-- Performance Impact: [Expected impact on database performance]
--
-- Rollback Notes: [Any special considerations for rolling back]
--
-- =============================================================================

-- Example: Create a new table with proper constraints and indexes
CREATE TABLE IF NOT EXISTS example_table (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'active' NOT NULL,
    parent_id UUID REFERENCES parent_table(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Add constraints
    CONSTRAINT chk_example_status CHECK (status IN ('active', 'inactive', 'deleted'))
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_example_table_name ON example_table(name);
CREATE INDEX IF NOT EXISTS idx_example_table_status ON example_table(status);
CREATE INDEX IF NOT EXISTS idx_example_table_parent_id ON example_table(parent_id);
CREATE INDEX IF NOT EXISTS idx_example_table_created_at ON example_table(created_at DESC);

-- Add trigger for auto-updating updated_at (if update_updated_at_column function exists)
-- CREATE TRIGGER update_example_table_updated_at BEFORE UPDATE ON example_table
--     FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Example: Seed reference data (if needed)
-- INSERT INTO example_table (id, name, description, status) VALUES
--     ('550e8400-e29b-41d4-a716-446655440001', 'Default Entry', 'System default entry', 'active')
-- ON CONFLICT (id) DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- =============================================================================
-- ROLLBACK: [MIGRATION_NAME]
-- =============================================================================
--
-- This section rolls back all changes made in the UP migration
-- Order matters: reverse the order of creation (triggers, indexes, tables)
--
-- =============================================================================

-- Drop trigger first (if created)
-- DROP TRIGGER IF EXISTS update_example_table_updated_at ON example_table;

-- Drop indexes
DROP INDEX IF EXISTS idx_example_table_created_at;
DROP INDEX IF EXISTS idx_example_table_parent_id;
DROP INDEX IF EXISTS idx_example_table_status;
DROP INDEX IF EXISTS idx_example_table_name;

-- Drop table (this will also drop constraints)
DROP TABLE IF EXISTS example_table;

-- Note: If you inserted reference data, you may want to remove it
-- DELETE FROM example_table WHERE id = '550e8400-e29b-41d4-a716-446655440001';

-- +goose StatementEnd
