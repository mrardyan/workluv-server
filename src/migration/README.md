# Database Migrations

This directory contains database migration files for the Go server application.

## Overview

Migrations are managed using [Goose](https://github.com/pressly/goose), a database migration tool for Go. Migrations are applied in order based on their timestamp prefix.

## Migration Files

Each migration file follows the naming convention:
```
YYYYMMDDHHMMSS_description.sql
```

For example:
- `20250812043912_account_service_accounts_table.sql`
- `20250812043913_workspace_service_workspaces_table.sql`
- `20250812093858_session_service_sessions_table.sql`
- `20250813105242_account_service_verifications_table.sql`

## Migration Structure

Each migration file contains:

1. **Up Migration**: The forward migration that creates/modifies database structures
2. **Down Migration**: The rollback migration that undoes the changes

Example structure:
```sql
-- +goose Up
-- This migration creates the users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
-- This migration drops the users table
DROP TABLE users;
```

## Running Migrations

### Prerequisites

1. **Install Goose**:
   ```bash
   go install github.com/pressly/goose/v3/cmd/goose@latest
   ```

2. **Set up environment**:
   ```bash
   ./src/scripts/setup.sh dev
   ```

3. **Configure database connection** in your `.env.development` file:
   ```bash
   DB_HOST=localhost
   DB_PORT=5432
   DB_NAME=your_database
   DB_USER=your_user
   DB_PASSWORD=your_password
   ```

### Basic Commands

```bash
# Apply all pending migrations
./src/scripts/migrate.sh up dev

# Apply one migration at a time
./src/scripts/migrate.sh up-by-one dev

# Rollback the last migration
./src/scripts/migrate.sh down dev

# Check migration status
./src/scripts/migrate.sh status dev

# View current version
./src/scripts/migrate.sh version dev
```

### Advanced Commands

```bash
# Create a new migration
./src/scripts/migrate.sh create migration_name dev

# Validate migrations
./src/scripts/migrate.sh validate

# Test a specific migration
./src/scripts/migrate.sh test migration_name dev

# Lint migrations for issues
./src/scripts/migrate.sh lint

# Check migration dependencies
./src/scripts/migrate.sh check-deps
```

## Migration Rules

1. **Always use the script**: Never create migration files manually
2. **Test migrations**: Always test both up and down migrations
3. **Validate before committing**: Run validation to ensure migrations follow rules
4. **Use descriptive names**: Migration names should clearly describe the change
5. **Include rollback logic**: Every migration must have a proper down migration

## Best Practices

1. **One change per migration**: Keep migrations focused and atomic
2. **Test thoroughly**: Test migrations in development before production
3. **Backup before production**: Always backup production database before applying migrations
4. **Review rollback procedures**: Ensure down migrations work correctly
5. **Document complex changes**: Add comments for complex database operations

## Troubleshooting

### Common Issues

1. **Migration fails to apply**:
   - Check database connection
   - Verify environment variables
   - Check migration file syntax

2. **Rollback fails**:
   - Ensure down migration is properly written
   - Check for data dependencies
   - Verify migration order

3. **Version conflicts**:
   - Check current migration version
   - Ensure all team members have latest migrations
   - Resolve conflicts manually if necessary

### Getting Help

- Check the [Goose documentation](https://github.com/pressly/goose)
- Review migration rules in `.cursor/05-migration-rules.mdc`
- Use validation commands to identify issues
- Check database logs for detailed error messages

## Environment-Specific Notes

### Development
- Use local PostgreSQL instance
- Apply migrations frequently during development
- Test all migrations before committing

### Staging
- Use staging database for testing
- Apply migrations before production
- Verify application functionality after migrations

### Production
- Always backup before applying migrations
- Apply migrations during maintenance windows
- Monitor application after migrations
- Have rollback plan ready

## Migration History

| Timestamp | Description | Status |
|-----------|-------------|---------|
| 20250812043912 | Create accounts table | Applied |
| 20250812043913 | Create workspaces table | Applied |
| 20250812093858 | Create sessions table | Applied |
| 20250813105242 | Create verifications table | Applied |

## Related Documentation

- [Migration Rules](../.cursor/05-migration-rules.mdc) - Detailed migration guidelines
- [Setup Script](../src/scripts/setup.sh) - Environment configuration
- [Deployment Guide](../DEPLOYMENT.md) - Production deployment instructions
