# Database Migrations

This directory contains database migration files managed by [Goose](https://github.com/pressly/goose).

## Setup

Make sure you have goose installed:
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Usage

### Creating a new migration
```bash
# Create a new SQL migration
goose -dir migration create migration_name sql

# Create a new Go migration
goose -dir migration create migration_name go
```

### Running migrations
```bash
# Apply all pending migrations
goose -dir migration postgres "your-database-connection-string" up

# Apply specific number of migrations
goose -dir migration postgres "your-database-connection-string" up-by-one

# Rollback the last migration
goose -dir migration postgres "your-database-connection-string" down

# Check migration status
goose -dir migration postgres "your-database-connection-string" status
```

### Database Connection String Format
```
postgres://username:password@host:port/database?sslmode=disable
```

## Migration File Naming Convention

Goose uses the following naming convention:
- `YYYYMMDDHHMMSS_migration_name.sql`
- `YYYYMMDDHHMMSS_migration_name.go`

## Best Practices

1. **Always test migrations**: Test both up and down migrations in a development environment
2. **Keep migrations small**: Break large changes into smaller, manageable migrations
3. **Use transactions**: Wrap your migrations in transactions when possible
4. **Backup before production**: Always backup your production database before running migrations
5. **Review rollback scripts**: Ensure down migrations properly reverse the up migration

## Environment Variables

You can set these environment variables for easier migration management:

```bash
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="postgres://username:password@host:port/database?sslmode=disable"
export GOOSE_MIGRATION_DIR=migration
```

Then you can run commands without specifying the driver and connection string:
```bash
goose up
goose down
goose status
```
