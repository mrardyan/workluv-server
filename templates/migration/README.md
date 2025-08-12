# Database Migration Templates

This directory contains templates and rules for database migrations using Goose with PostgreSQL.

## Overview

Migration templates ensure consistent structure and best practices across all database migrations in the project. These templates integrate with the automated validation system to maintain high quality and safety standards.

## Template Files

- `migration.template.sql` - Standard SQL migration template with proper structure and examples

## Template Structure

The migration template provides:

### Required Sections
- **Header Comments**: Business context, dependencies, and impact assessment
- **UP Migration**: Forward changes with proper safety checks
- **DOWN Migration**: Complete rollback procedures

### Built-in Best Practices
- **Idempotent Operations**: All operations use IF EXISTS/IF NOT EXISTS
- **Foreign Key Constraints**: Proper referential integrity with CASCADE behavior
- **Performance Indexes**: Indexes for foreign keys and frequently queried columns
- **Safety Checks**: Transaction-safe operations with proper error handling

### Documentation Standards
- **Business Context**: Clear explanation of why the change is needed
- **Dependencies**: List of prerequisite migrations or system changes
- **Performance Impact**: Expected impact on database and application performance
- **Rollback Notes**: Special considerations for the DOWN migration

## Usage

### Creating New Migrations

Use the migration script to create new migrations from the template:

```bash
# Create new migration with automatic template application
./scripts/migrate.sh create add_user_roles dev
```

The script will:
1. Validate the migration name against naming conventions
2. Generate a timestamped migration file
3. Apply the template structure automatically
4. Guide you to customize the template for your specific needs

### Template Customization

When using the template:

1. **Replace Placeholders**: Update `[MIGRATION_NAME]` and example content
2. **Add Business Context**: Explain why this migration is needed
3. **Document Dependencies**: List any prerequisite migrations
4. **Assess Performance**: Note expected impact on database performance
5. **Design Rollback**: Ensure DOWN migration safely reverses UP changes

### Validation

All migrations created from templates are automatically validated against project rules:

```bash
# Validate specific migration
./scripts/migrate.sh validate migration_name

# Validate all migrations
./scripts/migrate.sh validate

# Lint for common issues
./scripts/migrate.sh lint
```

## Template Features

### Safety by Default
- All table operations use IF EXISTS/IF NOT EXISTS
- Foreign key constraints include proper CASCADE behavior
- Indexes are created for performance optimization
- Transactions ensure atomicity

### Comprehensive Documentation
- Clear section headers and organization
- Business context and technical rationale
- Dependencies and prerequisites
- Performance impact assessment

### Rollback Reliability
- Complete reversal of UP migration changes
- Proper order of operations (reverse of creation order)
- Safety checks for all DROP operations

## Best Practices

### When Customizing Templates
1. **Keep Safety First**: Never remove IF EXISTS/IF NOT EXISTS clauses
2. **Document Thoroughly**: Explain the business need and technical approach
3. **Test Rollbacks**: Ensure DOWN migration works correctly
4. **Consider Performance**: Add appropriate indexes for query patterns

### Common Patterns
- **Adding Tables**: Use full template structure with constraints and indexes
- **Adding Columns**: Include safe ALTER TABLE operations
- **Data Migrations**: Separate from schema changes when possible
- **Index Changes**: Always include performance impact assessment

## Integration with Development Workflow

### Pre-Commit Validation
The template works with project validation rules to ensure:
- Proper naming conventions
- Required documentation sections
- Safety and performance best practices
- Complete rollback procedures

### Cursor IDE Integration
Template rules are integrated with Cursor IDE via `.cursor/05-migration-rules.mdc`, providing:
- Automatic code generation assistance
- Real-time validation feedback
- Best practice recommendations
- Anti-pattern detection

## File Organization
```
templates/migration/
├── README.md                # This documentation
└── migration.template.sql   # Standard migration template
```

## Related Documentation

- [Migration Rules](.cursor/05-migration-rules.mdc) - Cursor IDE integration
- [Migration README](../migration/README.md) - Usage and commands
- [Migration Rules](../migration/MIGRATION_RULES.md) - Detailed governance
- [Setup Script](../scripts/setup.sh) - Environment configuration

## Security Considerations

### Template Safety
- Never include hardcoded secrets or passwords
- Use environment variables for configuration
- Follow principle of least privilege
- Include proper constraint validation

### Production Readiness
- All templates are production-ready by default
- Safety checks prevent accidental data loss
- Performance considerations are built-in
- Rollback procedures are thoroughly tested
