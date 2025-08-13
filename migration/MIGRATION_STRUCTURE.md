# Migration Structure Changes

## Overview
This document outlines the migration structure changes made to separate tables by service and ensure all date fields use epoch seconds.

## Changes Made

### 1. Separated Monolithic Migration
**Before**: Single `20250812043912_initial_schema.sql` file containing multiple tables
**After**: Separate migration files for each table with clear service prefixes

### 2. Epoch Time Standardization
**Before**: Mixed `TIMESTAMP WITH TIME ZONE` and `BIGINT` date fields
**After**: All date fields consistently use `BIGINT` (epoch seconds)

## New Migration Structure

### Migration Files

| Timestamp | Service | Table | Description |
|-----------|---------|-------|-------------|
| `20250812043912` | `account_service` | `accounts` | Core user accounts |
| `20250812043913` | `workspace_service` | `workspaces` | User workspaces |
| `20250812093858` | `session_service` | `sessions` | User sessions |
| `20250813105242` | `account_service` | `verifications` | Email/2FA verification |

### File Naming Convention
```
YYYYMMDDHHMMSS_service_name_table_name.sql
```

**Examples:**
- `20250812043912_account_service_accounts_table.sql`
- `20250812093858_session_service_sessions_table.sql`

## Date Field Changes

### Before (TIMESTAMP)
```sql
created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
```

### After (Epoch Seconds)
```sql
created_at BIGINT NOT NULL DEFAULT EXTRACT(epoch FROM NOW()),
updated_at BIGINT NOT NULL DEFAULT EXTRACT(epoch FROM NOW()),
expires_at BIGINT NOT NULL, -- epoch seconds using shared.Time
```

## Benefits of New Structure

### 1. **Maintainability**
- Single table per migration file
- Clear service ownership
- Easier to track changes per table

### 2. **Consistency**
- All date fields use epoch seconds
- Aligns with Go `shared.Time` type
- No timezone conversion issues

### 3. **Performance**
- Epoch seconds are efficient for storage
- Better index performance for date queries
- Simplified date comparisons

### 4. **Developer Experience**
- Clear file naming indicates purpose
- Service-specific migrations are easier to find
- Better separation of concerns

## Migration Order

The migrations must run in this order due to foreign key dependencies:

1. **`20250812043912_account_service_accounts_table.sql`** - Base table
2. **`20250812043913_workspace_service_workspaces_table.sql`** - References accounts
3. **`20250812093858_session_service_sessions_table.sql`** - References accounts  
4. **`20250813105242_account_service_verifications_table.sql`** - References accounts

## Trigger Functions

Each service now has its own trigger function for `updated_at` fields:

```sql
-- Account service
CREATE OR REPLACE FUNCTION update_accounts_updated_at_column()

-- Workspace service  
CREATE OR REPLACE FUNCTION update_workspaces_updated_at_column()

-- Session service
CREATE OR REPLACE FUNCTION update_sessions_updated_at_column()

-- Verification service
CREATE OR REPLACE FUNCTION update_verifications_updated_at_column()
```

## Index Strategy

### Date Field Indexes
All date fields now have dedicated indexes for efficient range queries:

```sql
-- Created/Updated indexes for all tables
CREATE INDEX idx_table_name_created_at ON table_name(created_at);
CREATE INDEX idx_table_name_updated_at ON table_name(updated_at);

-- Expiration indexes where applicable
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
CREATE INDEX idx_verifications_expires_at ON verifications(expires_at);
```

### Service-Specific Indexes
Each service maintains its own set of indexes optimized for its use cases.

## Testing the New Structure

### 1. Run Migrations
```bash
# Test up migrations
goose postgres "connection_string" up

# Verify table structure
\d accounts
\d workspaces  
\d sessions
\d verifications
```

### 2. Verify Date Fields
```sql
-- Check that all date fields are BIGINT
SELECT column_name, data_type 
FROM information_schema.columns 
WHERE table_name IN ('accounts', 'workspaces', 'sessions', 'verifications')
  AND column_name LIKE '%_at';
```

### 3. Test Rollbacks
```bash
# Test down migrations
goose postgres "connection_string" down
```

## Future Migrations

When adding new tables or modifying existing ones:

1. **Use service prefix** in filename
2. **Focus on single table** per migration
3. **Use epoch seconds** for all date fields
4. **Include proper indexes** for date fields
5. **Test both up and down** migrations

## References
- [Epoch Time Usage Rules](../.cursor/rules/epoch-time-usage.md)
- [Shared Time Implementation](../internal/shared/value.go)
- [Migration README](./README.md)
