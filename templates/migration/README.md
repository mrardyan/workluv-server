# Migration Template

This template provides a standardized structure for creating new database migrations.

## Usage

### 1. Create New Migration

Always use the migration script to create new migrations:

```bash
./src/scripts/migrate.sh create migration_name dev
```

**Examples:**
```bash
./src/scripts/migrate.sh create add_user_roles dev
./src/scripts/migrate.sh create create_audit_table dev
./src/scripts/migrate.sh create update_workspace_permissions dev
```

### 2. Migration Naming Convention

- **Format**: `action_table_name` (snake_case)
- **Action verbs**: `add`, `create`, `remove`, `drop`, `update`, `alter`, `fix`, `migrate`, `seed`
- **Examples**:
  - `add_user_roles` - Add new roles to users table
  - `create_audit_table` - Create new audit logging table
  - `update_workspace_permissions` - Modify workspace permission structure

### 3. File Structure

The script will create a file with the current timestamp:

```sql
-- +goose Up
-- +goose StatementBegin
-- Your UP migration SQL here
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Your DOWN migration SQL here
-- +goose StatementEnd
```

## Migration Rules

### 1. Always Use Scripts

**✅ CORRECT:**
```bash
./src/scripts/migrate.sh create migration_name dev
```

**❌ INCORRECT:**
- Creating files manually
- Copying existing migrations
- Renaming migration files

### 2. Required Elements

Every migration must include:

1. **Goose Comments**: `-- +goose Up` and `-- +goose Down`
2. **Up Migration**: Forward migration logic
3. **Down Migration**: Rollback logic
4. **Comments**: Explain what the migration does

### 3. SQL Best Practices

**✅ GOOD:**
```sql
-- Create users table with proper constraints
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at BIGINT NOT NULL DEFAULT EXTRACT(epoch FROM NOW()),
    updated_at BIGINT NOT NULL DEFAULT EXTRACT(epoch FROM NOW())
);

-- Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
```

**❌ AVOID:**
```sql
-- Missing IF NOT EXISTS
CREATE TABLE users (...);

-- Missing indexes on foreign keys
CREATE TABLE posts (
    user_id INTEGER REFERENCES users(id)
    -- Missing: CREATE INDEX idx_posts_user_id ON posts(user_id);
);
```

### 4. Rollback Safety

**✅ SAFE:**
```sql
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users CASCADE;
-- +goose StatementEnd
```

**❌ DANGEROUS:**
```sql
-- +goose Down
-- +goose StatementBegin
DROP TABLE users; -- No IF EXISTS, could fail
-- +goose StatementEnd
```

## Validation

### 1. Validate Before Committing

```bash
# Validate all migrations
./src/scripts/migrate.sh validate

# Validate specific migration
./src/scripts/migrate.sh validate migration_name
```

### 2. Lint for Issues

```bash
# Check for common problems
./src/scripts/migrate.sh lint
```

### 3. Test Migrations

```bash
# Test specific migration
./src/scripts/migrate.sh test migration_name dev
```

## Common Patterns

### 1. Adding Columns

```sql
-- +goose Up
-- +goose StatementBegin
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS phone VARCHAR(20),
ADD COLUMN IF NOT EXISTS verified_at BIGINT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users 
DROP COLUMN IF EXISTS phone,
DROP COLUMN IF EXISTS verified_at;
-- +goose StatementEnd
```

### 2. Creating Tables

```sql
-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    bio TEXT,
    avatar_url VARCHAR(500),
    created_at BIGINT NOT NULL DEFAULT EXTRACT(epoch FROM NOW()),
    updated_at BIGINT NOT NULL DEFAULT EXTRACT(epoch FROM NOW())
);

-- Add indexes
CREATE INDEX IF NOT EXISTS idx_user_profiles_user_id ON user_profiles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_profiles_created_at ON user_profiles(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_profiles CASCADE;
-- +goose StatementEnd
```

### 3. Adding Constraints

```sql
-- +goose Up
-- +goose StatementBegin
-- Add unique constraint
ALTER TABLE users 
ADD CONSTRAINT IF NOT EXISTS uk_users_email UNIQUE (email);

-- Add check constraint
ALTER TABLE users 
ADD CONSTRAINT IF NOT EXISTS chk_users_age 
CHECK (age >= 13 AND age <= 120);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove constraints
ALTER TABLE users DROP CONSTRAINT IF EXISTS uk_users_email;
ALTER TABLE users DROP CONSTRAINT IF EXISTS chk_users_age;
-- +goose StatementEnd
```

## Troubleshooting

### 1. Migration Fails

**Check:**
- SQL syntax errors
- Missing dependencies
- Database connection
- Environment variables

**Debug:**
```bash
# Check migration status
./src/scripts/migrate.sh status dev

# View database logs
docker-compose logs postgres
```

### 2. Rollback Issues

**Common causes:**
- Missing `IF EXISTS` clauses
- Data dependencies
- Constraint violations

**Solution:**
- Always use `IF EXISTS` for drops
- Test rollbacks in development
- Handle data migration carefully

### 3. Version Conflicts

**Symptoms:**
- Migration version mismatch
- Applied migrations not in sync

**Solution:**
```bash
# Check current version
./src/scripts/migrate.sh version dev

# Reset if needed (development only)
./src/scripts/migrate.sh reset dev
```

## Environment Setup

### 1. Development

```bash
# Setup environment
./src/scripts/setup.sh dev

# Edit .env.development with database details
nano .env.development

# Test connection
./src/scripts/migrate.sh status dev
```

### 2. Staging/Production

```bash
# Setup environment
./src/scripts/setup.sh staging
./src/scripts/setup.sh prod

# Configure database credentials
nano .env.staging
nano .env.production
```

## Related Scripts

- [Setup Script](../src/scripts/setup.sh) - Environment configuration
- [Deploy Script](../src/scripts/deploy.sh) - Production deployment
- [Migration Script](../src/scripts/migrate.sh) - Migration management

## Additional Resources

- [Goose Documentation](https://github.com/pressly/goose)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Migration Rules](../../.cursor/05-migration-rules.mdc)
