# Database Migration Tool

## Overview

This is a database migration tool built with goose that provides runtime migration capabilities for the Go server application. It's designed to work seamlessly with DigitalOcean App Platform and uses the same configuration system as the main server.

## Features

- **Goose Integration**: Uses the powerful goose migration framework
- **Configuration Consistency**: Uses the same config system as the main server
- **Environment-based**: Configurable via environment variables
- **Runtime Execution**: Can be run at deployment time or as a separate job
- **PostgreSQL Support**: Optimized for PostgreSQL databases

## Usage

### Local Development

```bash
# Apply all migrations
go run ./cmd/migrate up

# Check migration status
go run ./cmd/migrate status

# Roll back one migration
go run ./cmd/migrate down

# Create a new migration
go run ./cmd/migrate create add_new_table

# Migrate to specific version
go run ./cmd/migrate up-to 20250812043912
```

### Build as Binary

```bash
# Build the migration tool
go build -o migrate ./cmd/migrate

# Use the binary
./migrate up
./migrate status
```

## DigitalOcean App Platform Integration

### Method 1: Pre-Deploy Hook (Recommended)

Add a pre-deploy command in your app spec to run migrations before starting the server:

```yaml
name: my-go-app
services:
- name: api
  source_dir: /
  github:
    repo: your-repo/go-server
    branch: main
  run_command: go run ./cmd/server/main.go
  environment_slug: go
  instance_count: 1
  instance_size_slug: basic-xxs
  http_port: 8080
  # Pre-deploy hook to run migrations
  pre_deploy:
    command: "go run ./cmd/migrate up"
  envs:
  - key: DATABASE_URL
    scope: RUN_AND_BUILD_TIME
    type: SECRET
    value: ${db.DATABASE_URL}
  - key: DATABASE_HOST
    scope: RUN_AND_BUILD_TIME
    value: ${db.HOSTNAME}
  - key: DATABASE_PORT
    scope: RUN_AND_BUILD_TIME
    value: ${db.PORT}
  - key: DATABASE_NAME
    scope: RUN_AND_BUILD_TIME
    value: ${db.DATABASE}
  - key: DATABASE_USER
    scope: RUN_AND_BUILD_TIME
    value: ${db.USERNAME}
  - key: DATABASE_PASSWORD
    scope: RUN_AND_BUILD_TIME
    type: SECRET
    value: ${db.PASSWORD}
  - key: DATABASE_SSL_MODE
    scope: RUN_AND_BUILD_TIME
    value: require
```

### Method 2: Separate Job Component

Create a separate job component that runs migrations:

```yaml
jobs:
- name: migrate
  source_dir: /
  github:
    repo: your-repo/go-server
    branch: main
  run_command: go run ./cmd/migrate up
  environment_slug: go
  instance_count: 1
  instance_size_slug: basic-xxs
  kind: PRE_DEPLOY
  envs:
  - key: DATABASE_HOST
    scope: RUN_TIME
    value: ${db.HOSTNAME}
  - key: DATABASE_PORT
    scope: RUN_TIME
    value: ${db.PORT}
  - key: DATABASE_NAME
    scope: RUN_TIME
    value: ${db.DATABASE}
  - key: DATABASE_USER
    scope: RUN_TIME
    value: ${db.USERNAME}
  - key: DATABASE_PASSWORD
    scope: RUN_TIME
    type: SECRET
    value: ${db.PASSWORD}
  - key: DATABASE_SSL_MODE
    scope: RUN_TIME
    value: require
```

### Method 3: Manual Execution

You can also run migrations manually by connecting to your app console:

```bash
# Connect to your app console
doctl apps create-deployment <app-id> --wait

# Or use the DigitalOcean console to access your running container
# Then run:
go run ./cmd/migrate up
```

## Environment Configuration

The migration tool uses the same environment variables as the main server:

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_HOST` | PostgreSQL host | `localhost` |
| `DATABASE_PORT` | PostgreSQL port | `5432` |
| `DATABASE_NAME` | Database name | `golang_arch` |
| `DATABASE_USER` | Database user | `postgres` |
| `DATABASE_PASSWORD` | Database password | `password` |
| `DATABASE_SSL_MODE` | SSL mode | `disable` |

## Commands

| Command | Description | Example |
|---------|-------------|---------|
| `up` | Apply all available migrations | `migrate up` |
| `up-to VERSION` | Migrate up to a specific version | `migrate up-to 20250812043912` |
| `down` | Roll back the latest migration | `migrate down` |
| `down-to VERSION` | Roll back to a specific version | `migrate down-to 20250812043912` |
| `redo` | Re-run the latest migration | `migrate redo` |
| `reset` | Roll back all migrations | `migrate reset` |
| `status` | Show migration status | `migrate status` |
| `version` | Show current migration version | `migrate version` |
| `create NAME [TYPE]` | Create a new migration file | `migrate create add_users sql` |
| `fix` | Fix migration file sequence | `migrate fix` |

## Migration Files

Migration files are stored in the `migration/` directory and follow goose naming convention:

- `YYYYMMDDHHMMSS_description.sql` for SQL migrations
- `YYYYMMDDHHMMSS_description.go` for Go migrations

Example migration file structure:
```
migration/
├── 20250812043912_initial_schema.sql
└── 20250813120000_add_user_indexes.sql
```

## Best Practices

1. **Always test migrations locally first**
2. **Use transactions in your SQL migrations** (goose handles this automatically for SQL)
3. **Create reversible migrations** (implement both Up and Down)
4. **Run migrations before deploying new code**
5. **Monitor migration execution** in production logs
6. **Backup your database** before running migrations in production

## Troubleshooting

### Connection Issues

If you see connection refused errors:
- Verify your database environment variables
- Ensure the database service is running
- Check firewall/security group settings

### Migration Conflicts

If you see version conflicts:
- Check migration status with `migrate status`
- Ensure all team members have pulled latest migrations
- Use `migrate fix` to fix sequence issues

### Permission Issues

If migrations fail with permission errors:
- Verify the database user has appropriate privileges
- Ensure the user can create/modify tables and schemas

## Examples

### Creating a New Migration

```bash
# Create a new SQL migration
go run ./cmd/migrate create add_user_profiles

# This creates: migration/YYYYMMDDHHMMSS_add_user_profiles.sql
```

### Applying Migrations in Production

```bash
# Check current status
go run ./cmd/migrate status

# Apply all pending migrations
go run ./cmd/migrate up

# Verify the result
go run ./cmd/migrate version
```

This migration tool provides a robust, production-ready solution for managing database schema changes in your Go application deployed on DigitalOcean App Platform.
