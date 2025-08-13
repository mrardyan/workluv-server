# Database Migrations

This directory contains database migration files for the Go server project. All migrations use [Goose](https://github.com/pressly/goose) for database versioning.

## Migration Structure

Migrations are organized by service and table, with clear naming conventions:

### Naming Convention
- **Format**: `YYYYMMDDHHMMSS_service_name_table_name.sql`
- **Service Prefix**: Each migration is prefixed with the service name (e.g., `account_service_`, `session_service_`)
- **Table Focus**: Each migration file focuses on a single table for better maintainability

### Migration Files

#### 1. Account Service
- **`20250812043912_account_service_accounts_table.sql`** - Core accounts table
- **`20250813105242_account_service_verifications_table.sql`** - Email verification and 2FA tokens

#### 2. Session Service  
- **`20250812093858_session_service_sessions_table.sql`** - User session management

#### 3. Workspace Service
- **`20250812043913_workspace_service_workspaces_table.sql`** - User workspaces and projects

## Epoch Time Standard

**All date/time fields use epoch seconds (BIGINT) to align with the `shared.Time` type.**

### Date Field Types
- **`created_at`**: BIGINT NOT NULL DEFAULT EXTRACT(epoch FROM NOW())
- **`updated_at`**: BIGINT NOT NULL DEFAULT EXTRACT(epoch FROM NOW())
- **`expires_at`**: BIGINT NOT NULL (for time-limited operations)
- **`completed_at`**: BIGINT NULL (for completion timestamps)

### Benefits of Epoch Seconds
1. **Consistency**: Aligns with Go's `shared.Time` type
2. **Performance**: Efficient storage and comparison operations
3. **Timezone Independence**: No timezone conversion issues
4. **JSON Serialization**: Natural epoch representation in APIs

### Reference Implementation
See `internal/shared/value.go` for the `Time` type implementation and conversion methods.

## Running Migrations

### Local Development
```bash
# Run all pending migrations
goose postgres "postgres://username:password@localhost/dbname?sslmode=disable" up

# Rollback last migration
goose postgres "postgres://username:password@localhost/dbname?sslmode=disable" down

# Check migration status
goose postgres "postgres://username:password@localhost/dbname?sslmode=disable" status
```

### Docker Environment
```bash
# Run migrations in Docker
./scripts/migrate.sh

# Or manually
docker-compose exec app goose postgres "postgres://username:password@db/dbname?sslmode=disable" up
```

## Migration Best Practices

### 1. Single Table Focus
- Each migration file should focus on one table
- Include all related indexes and constraints
- Use descriptive names that indicate the service and table

### 2. Epoch Time Fields
- Always use `BIGINT` for date fields
- Set appropriate defaults using `EXTRACT(epoch FROM NOW())`
- Include comments indicating epoch seconds usage

### 3. Indexes
- Create indexes for frequently queried fields
- Include indexes for date fields used in range queries
- Consider composite indexes for common query patterns

### 4. Triggers
- Use service-specific trigger functions for `updated_at` fields
- Ensure triggers update epoch seconds, not timestamps

### 5. Rollback Safety
- Always include proper `Down` migrations
- Test rollbacks in development before production
- Consider data preservation during rollbacks

## Adding New Migrations

### 1. Create Migration File
```bash
# Generate new migration
goose create add_new_table sql
```

### 2. Follow Naming Convention
- Use service prefix (e.g., `account_service_`)
- Include table name
- Use descriptive action (e.g., `add_`, `modify_`, `drop_`)

### 3. Implement Epoch Time Fields
```sql
-- ✅ CORRECT - Use epoch seconds
created_at BIGINT NOT NULL DEFAULT EXTRACT(epoch FROM NOW()),
updated_at BIGINT NOT NULL DEFAULT EXTRACT(epoch FROM NOW()),
expires_at BIGINT NOT NULL,

-- ❌ INCORRECT - Don't use timestamps
created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
```

### 4. Include Proper Indexes
```sql
-- Index for date range queries
CREATE INDEX idx_table_name_created_at ON table_name(created_at);
CREATE INDEX idx_table_name_expires_at ON table_name(expires_at);
```

### 5. Test Migration
```bash
# Test up migration
goose postgres "connection_string" up

# Test down migration  
goose postgres "connection_string" down

# Verify table structure
\d table_name
```

## Troubleshooting

### Common Issues
1. **Migration Order**: Ensure migrations run in chronological order
2. **Dependencies**: Tables with foreign keys must be created after referenced tables
3. **Epoch Conversion**: Verify all date fields use `EXTRACT(epoch FROM NOW())`
4. **Index Conflicts**: Check for duplicate index names across migrations

### Debugging
```bash
# Check migration status
goose postgres "connection_string" status

# View migration logs
goose postgres "connection_string" version

# Reset migrations (development only)
goose postgres "connection_string" reset
```

## References
- [Goose Documentation](https://github.com/pressly/goose)
- [PostgreSQL Epoch Functions](https://www.postgresql.org/docs/current/functions-datetime.html)
- [Project Time Handling](internal/shared/value.go)
