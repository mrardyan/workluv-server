#!/bin/bash

# =============================================================================
# GO-SERVER - DATABASE MIGRATION SCRIPT
# =============================================================================
# 
# This script provides convenient commands for managing database migrations
# using Goose with PostgreSQL.
#
# Usage: ./scripts/migrate.sh [command] [environment] [options]
# =============================================================================

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${BLUE}=== $1 ===${NC}"
}

# Function to check if migration tool is available
check_migration_tool() {
    # Check if we're in a container (use binary) or local (use go run)
    if [ -f "./migrate" ]; then
        print_status "Using pre-built migration binary"
        return 0
    elif [ -f "./cmd/migrate/main.go" ]; then
        print_status "Using local migration tool (go run)"
        return 0
    else
        print_error "Migration tool not found. Please ensure cmd/migrate/ directory exists."
        exit 1
    fi
}

# Function to get database connection string from environment file
get_db_string() {
    local environment="$1"
    local env_file=""
    
    case "$environment" in
        "dev"|"development")
            env_file=".env.development"
            ;;
        "staging")
            env_file=".env.staging"
            ;;
        "prod"|"production")
            env_file=".env.production"
            ;;
        *)
            print_error "Unknown environment: $environment"
            show_help
            exit 1
            ;;
    esac
    
    if [ ! -f "$env_file" ]; then
        print_error "Environment file not found: $env_file"
        print_error "Please run ./scripts/setup.sh $environment first"
        exit 1
    fi
    
    # Source the environment file and construct connection string
    source "$env_file"
    
    if [ -z "$DB_HOST" ] || [ -z "$DB_PORT" ] || [ -z "$DB_NAME" ] || [ -z "$DB_USER" ] || [ -z "$DB_PASSWORD" ]; then
        print_error "Missing database configuration in $env_file"
        print_error "Required variables: DB_HOST, DB_PORT, DB_NAME, DB_USER, DB_PASSWORD"
        exit 1
    fi
    
    echo "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"
}

# Function to validate migration naming convention
validate_migration_name() {
    local name="$1"
    
    # Check if name follows snake_case and starts with action verb
    if [[ ! "$name" =~ ^[a-z]+_[a-z0-9_]+$ ]]; then
        print_error "Migration name must follow snake_case convention: $name"
        print_error "Examples: add_user_roles, create_audit_table, update_workspace_permissions"
        return 1
    fi
    
    # Check for common action verbs
    local first_word=$(echo "$name" | cut -d'_' -f1)
    local valid_verbs=("add" "create" "remove" "drop" "update" "alter" "fix" "migrate" "seed")
    local is_valid=false
    
    for verb in "${valid_verbs[@]}"; do
        if [[ "$first_word" == "$verb" ]]; then
            is_valid=true
            break
        fi
    done
    
    if [[ "$is_valid" == false ]]; then
        print_warning "Migration name should start with an action verb"
        print_warning "Recommended verbs: ${valid_verbs[*]}"
    fi
    
    return 0
}

# Function to validate migration file content
validate_migration_file() {
    local file_path="$1"
    local errors=0
    
    print_status "Validating migration file: $(basename "$file_path")"
    
    # Check if file exists
    if [[ ! -f "$file_path" ]]; then
        print_error "Migration file not found: $file_path"
        return 1
    fi
    
    # Check for required goose comments
    if ! grep -q -- "-- +goose Up" "$file_path"; then
        print_error "Missing '-- +goose Up' comment"
        errors=$((errors + 1))
    fi
    
    if ! grep -q -- "-- +goose Down" "$file_path"; then
        print_error "Missing '-- +goose Down' comment"
        errors=$((errors + 1))
    fi
    
    # Check for dangerous operations without IF EXISTS/IF NOT EXISTS
    if grep -q "DROP TABLE\|DROP INDEX\|DROP CONSTRAINT" "$file_path" && ! grep -q "IF EXISTS" "$file_path"; then
        print_error "Found DROP operations without IF EXISTS - migrations should be idempotent"
        errors=$((errors + 1))
    fi
    
    if grep -q "CREATE TABLE\|CREATE INDEX\|CREATE CONSTRAINT" "$file_path" && ! grep -q "IF NOT EXISTS" "$file_path"; then
        print_warning "CREATE operations should use IF NOT EXISTS for idempotency"
    fi
    
    # Check for foreign key constraints
    if grep -q "CREATE TABLE" "$file_path" && grep -q "REFERENCES\|FOREIGN KEY" "$file_path"; then
        print_success "Good: Found foreign key constraints"
    elif grep -q "CREATE TABLE" "$file_path"; then
        print_warning "Consider adding foreign key constraints if relationships exist"
    fi
    
    # Check for indexes on foreign keys
    if grep -q "REFERENCES" "$file_path" && ! grep -q "CREATE INDEX" "$file_path"; then
        print_warning "Consider adding indexes for foreign key columns"
    fi
    
    # Check for comments explaining the migration
    if ! grep -q "^--[[:space:]]*[A-Z]" "$file_path"; then
        print_warning "Consider adding comments explaining the purpose of this migration"
    fi
    
    # Check DOWN migration has actual rollback code
    local down_section=$(sed -n '/-- +goose Down/,$p' "$file_path")
    if [[ $(echo "$down_section" | grep -c "^[[:space:]]*[A-Z]") -eq 0 ]]; then
        print_error "DOWN migration appears to be empty - rollback procedures are required"
        errors=$((errors + 1))
    fi
    
    if [[ $errors -eq 0 ]]; then
        print_success "Migration validation passed"
        return 0
    else
        print_error "Migration validation failed with $errors error(s)"
        return 1
    fi
}

# Function to test migration up and down
test_migration() {
    local environment="$1"
    local migration_file="$2"
    
    print_header "Testing Migration: $(basename "$migration_file")"
    
    local db_string=$(get_db_string "$environment")
    local migration_dir="migration"
    
    # Get current version
    local current_version=$(goose -dir "$migration_dir" postgres "$db_string" version 2>/dev/null | tail -1)
    
    print_status "Current database version: $current_version"
    print_status "Testing UP migration..."
    
    # Test up migration
    if goose -dir "$migration_dir" postgres "$db_string" up-by-one; then
        print_success "UP migration successful"
        
        print_status "Testing DOWN migration..."
        # Test down migration
        if goose -dir "$migration_dir" postgres "$db_string" down; then
            print_success "DOWN migration successful"
            print_success "Migration test completed successfully"
            return 0
        else
            print_error "DOWN migration failed"
            return 1
        fi
    else
        print_error "UP migration failed"
        return 1
    fi
}

# Function to check migration dependencies
check_migration_dependencies() {
    local migration_dir="migration"
    
    print_header "Checking Migration Dependencies"
    
    # Get list of migration files
    local migration_files=($(ls "$migration_dir"/*.sql 2>/dev/null | sort))
    
    if [[ ${#migration_files[@]} -eq 0 ]]; then
        print_warning "No migration files found"
        return 0
    fi
    
    print_status "Found ${#migration_files[@]} migration files"
    
    # Check for potential dependency issues
    local tables_created=()
    local tables_dropped=()
    local dependencies_ok=true
    
    for file in "${migration_files[@]}"; do
        local filename=$(basename "$file")
        print_status "Analyzing: $filename"
        
        # Extract table names from CREATE TABLE statements
        while read -r line; do
            if [[ "$line" =~ CREATE[[:space:]]+TABLE[[:space:]]+.*[[:space:]]([a-zA-Z_][a-zA-Z0-9_]*) ]]; then
                local table_name="${BASH_REMATCH[1]}"
                tables_created+=("$table_name")
            fi
        done < "$file"
        
        # Extract table names from DROP TABLE statements
        while read -r line; do
            if [[ "$line" =~ DROP[[:space:]]+TABLE[[:space:]]+.*[[:space:]]([a-zA-Z_][a-zA-Z0-9_]*) ]]; then
                local table_name="${BASH_REMATCH[1]}"
                tables_dropped+=("$table_name")
            fi
        done < "$file"
        
        # Check for REFERENCES to tables that might not exist yet
        while read -r line; do
            if [[ "$line" =~ REFERENCES[[:space:]]+([a-zA-Z_][a-zA-Z0-9_]*) ]]; then
                local referenced_table="${BASH_REMATCH[1]}"
                local found=false
                for created_table in "${tables_created[@]}"; do
                    if [[ "$created_table" == "$referenced_table" ]]; then
                        found=true
                        break
                    fi
                done
                if [[ "$found" == false ]]; then
                    print_warning "In $filename: References table '$referenced_table' which may not exist yet"
                    dependencies_ok=false
                fi
            fi
        done < "$file"
    done
    
    if [[ "$dependencies_ok" == true ]]; then
        print_success "No dependency issues found"
    else
        print_warning "Potential dependency issues detected - review migration order"
    fi
    
    return 0
}

# Function to lint migrations for common issues
lint_migrations() {
    local migration_dir="migration"
    
    print_header "Linting Migration Files"
    
    local migration_files=($(ls "$migration_dir"/*.sql 2>/dev/null | sort))
    local total_issues=0
    
    for file in "${migration_files[@]}"; do
        local filename=$(basename "$file")
        local issues=0
        
        print_status "Linting: $filename"
        
        # Check for common anti-patterns
        if grep -q "SELECT.*FROM.*WHERE.*LIKE.*%" "$file"; then
            print_warning "  Found LIKE queries with leading wildcards - consider performance impact"
            issues=$((issues + 1))
        fi
        
        if grep -q "VARCHAR(255)" "$file"; then
            print_warning "  Using VARCHAR(255) - consider if a smaller size is appropriate"
            issues=$((issues + 1))
        fi
        
        if grep -q "TEXT" "$file" && grep -q "CREATE TABLE" "$file"; then
            print_warning "  Using TEXT columns - consider if VARCHAR(n) would be more appropriate"
            issues=$((issues + 1))
        fi
        
        if grep -q "UPDATE.*SET.*WHERE" "$file" && ! grep -q "WHERE.*PRIMARY KEY\|WHERE.*UNIQUE" "$file"; then
            print_warning "  UPDATE statements should use specific WHERE clauses"
            issues=$((issues + 1))
        fi
        
        if grep -q "DELETE FROM" "$file" && ! grep -q "WHERE" "$file"; then
            print_error "  DELETE without WHERE clause found - this is dangerous!"
            issues=$((issues + 1))
        fi
        
        # Check for missing semicolons
        local sql_lines=$(grep -v "^--" "$file" | grep -v "^$")
        if echo "$sql_lines" | grep -q "[^;]$"; then
            print_warning "  Some SQL statements may be missing semicolons"
            issues=$((issues + 1))
        fi
        
        if [[ $issues -eq 0 ]]; then
            print_success "  No issues found"
        else
            print_warning "  Found $issues issue(s)"
            total_issues=$((total_issues + issues))
        fi
    done
    
    if [[ $total_issues -eq 0 ]]; then
        print_success "All migrations passed linting"
    else
        print_warning "Found $total_issues total issue(s) across all migrations"
    fi
    
    return 0
}

# Function to run migration command
run_migration() {
    local command="$1"
    local environment="$2"
    local migration_name="$3"
    local migration_dir="src/migration"
    
    # Commands that don't need database connection
    local no_db_commands=("lint" "check-deps" "create" "create-go" "validate")
    local needs_db=true
    
    for no_db_cmd in "${no_db_commands[@]}"; do
        if [[ "$command" == "$no_db_cmd" ]]; then
            needs_db=false
            break
        fi
    done
    
    # Get database connection only if needed
    local db_string=""
    if [ "$needs_db" == true ]; then
        db_string=$(get_db_string "$environment")
        print_status "Running goose $command for $environment environment..."
        print_status "Database: $DB_HOST:$DB_PORT/$DB_NAME"
    else
        print_status "Running $command..."
    fi
    
    case "$command" in
        "create")
            if [ -z "$migration_name" ]; then
                print_error "Migration name is required for create command"
                exit 1
            fi
            # Validate migration name
            if ! validate_migration_name "$migration_name"; then
                exit 1
            fi
            # Use our migration tool instead of goose
            if [ -f "./migrate" ]; then
                ./migrate create "$migration_name"
            elif [ -f "./cmd/migrate/main.go" ]; then
                go run ./cmd/migrate/main.go create "$migration_name"
            else
                print_error "Migration tool not found"
                exit 1
            fi
            # Validate the created file
            local latest_file=$(ls -t "$migration_dir"/*.sql | head -1)
            print_status "Created migration file: $(basename "$latest_file")"
            print_status "Please edit the migration file and follow the rules in .cursor/05-migration-rules.mdc"
            ;;
        "create-go")
            if [ -z "$migration_name" ]; then
                print_error "Migration name is required for create-go command"
                exit 1
            fi
            # Validate migration name
            if ! validate_migration_name "$migration_name"; then
                exit 1
            fi
            # Use our migration tool instead of goose
            if [ -f "./migrate" ]; then
                ./migrate create "$migration_name"
            elif [ -f "./cmd/migrate/main.go" ]; then
                go run ./cmd/migrate/main.go create "$migration_name"
            else
                print_error "Migration tool not found"
                exit 1
            fi
            ;;
        "validate")
            if [ -n "$migration_name" ]; then
                # Validate specific migration file
                local file_pattern="$migration_dir/*$migration_name*.sql"
                local matching_files=($(ls $file_pattern 2>/dev/null))
                if [[ ${#matching_files[@]} -eq 0 ]]; then
                    print_error "No migration file found matching: $migration_name"
                    exit 1
                elif [[ ${#matching_files[@]} -gt 1 ]]; then
                    print_error "Multiple migration files found matching: $migration_name"
                    for file in "${matching_files[@]}"; do
                        echo "  - $(basename "$file")"
                    done
                    exit 1
                else
                    validate_migration_file "${matching_files[0]}"
                fi
            else
                # Validate all migration files
                print_header "Validating All Migrations"
                local migration_files=($(ls "$migration_dir"/*.sql 2>/dev/null | sort))
                local total_errors=0
                for file in "${migration_files[@]}"; do
                    if ! validate_migration_file "$file"; then
                        total_errors=$((total_errors + 1))
                    fi
                done
                if [[ $total_errors -eq 0 ]]; then
                    print_success "All migrations passed validation"
                else
                    print_error "$total_errors migration(s) failed validation"
                    exit 1
                fi
            fi
            ;;
        "test")
            if [ -z "$migration_name" ]; then
                print_error "Migration name is required for test command"
                exit 1
            fi
            local file_pattern="$migration_dir/*$migration_name*.sql"
            local matching_files=($(ls $file_pattern 2>/dev/null))
            if [[ ${#matching_files[@]} -eq 0 ]]; then
                print_error "No migration file found matching: $migration_name"
                exit 1
            elif [[ ${#matching_files[@]} -gt 1 ]]; then
                print_error "Multiple migration files found matching: $migration_name"
                exit 1
            else
                test_migration "$environment" "${matching_files[0]}"
            fi
            ;;
        "check-deps")
            check_migration_dependencies
            ;;
        "lint")
            lint_migrations
            ;;
        "up")
            if [ -f "./migrate" ]; then
                ./migrate up
            elif [ -f "./cmd/migrate/main.go" ]; then
                go run ./cmd/migrate/main.go up
            else
                print_error "Migration tool not found"
                exit 1
            fi
            ;;
        "up-by-one")
            if [ -f "./migrate" ]; then
                ./migrate up
            elif [ -f "./cmd/migrate/main.go" ]; then
                go run ./cmd/migrate/main.go up
            else
                print_error "Migration tool not found"
                exit 1
            fi
            ;;
        "down")
            if [ -f "./migrate" ]; then
                ./migrate down
            elif [ -f "./cmd/migrate/main.go" ]; then
                go run ./cmd/migrate/main.go down
            else
                print_error "Migration tool not found"
                exit 1
            fi
            ;;
        "reset")
            print_warning "This will rollback ALL migrations. Are you sure? (y/N)"
            read -r response
            if [[ "$response" =~ ^[Yy]$ ]]; then
                if [ -f "./migrate" ]; then
                    ./migrate reset
                elif [ -f "./cmd/migrate/main.go" ]; then
                    go run ./cmd/migrate/main.go reset
                else
                    print_error "Migration tool not found"
                    exit 1
                fi
            else
                print_status "Reset cancelled"
                exit 0
            fi
            ;;
        "status")
            if [ -f "./migrate" ]; then
                ./migrate status
            elif [ -f "./cmd/migrate/main.go" ]; then
                go run ./cmd/migrate/main.go status
            else
                print_error "Migration tool not found"
                exit 1
            fi
            ;;
        "version")
            if [ -f "./migrate" ]; then
                ./migrate version
            elif [ -f "./cmd/migrate/main.go" ]; then
                go run ./cmd/migrate/main.go version
            else
                exit 1
            fi
            ;;
        *)
            print_error "Unknown command: $command"
            show_help
            exit 1
            ;;
    esac
}

# Function to show help
show_help() {
    echo "Go-Server - Database Migration Script"
    echo ""
    echo "Usage: $0 [command] [environment] [migration_name]"
    echo ""
    echo "Migration Commands:"
    echo "  create NAME ENV        - Create a new SQL migration file"
    echo "  create-go NAME ENV     - Create a new Go migration file"
    echo "  up ENV                 - Apply all pending migrations"
    echo "  up-by-one ENV          - Apply one pending migration"
    echo "  down ENV               - Rollback the last migration"
    echo "  reset ENV              - Rollback ALL migrations (destructive!)"
    echo "  status ENV             - Show migration status"
    echo "  version ENV            - Show current migration version"
    echo ""
    echo "Validation Commands:"
    echo "  validate [NAME]        - Validate migration file(s) against rules"
    echo "  test NAME ENV          - Test migration up and down (requires DB)"
    echo "  check-deps             - Check for migration dependency issues"
    echo "  lint                   - Lint migrations for common issues"
    echo ""
    echo "Environments:"
    echo "  dev/development        - Development environment"
    echo "  staging               - Staging environment"
    echo "  prod/production       - Production environment"
    echo ""
    echo "Examples:"
    echo "  # Basic migration operations"
    echo "  $0 create add_users_table dev          # Create new migration"
    echo "  $0 up dev                              # Apply all pending migrations"
    echo "  $0 status dev                          # Check migration status"
    echo "  $0 down prod                           # Rollback last migration in prod"
    echo ""
    echo "  # Validation and testing"
    echo "  $0 validate                            # Validate all migrations"
    echo "  $0 validate add_users_table            # Validate specific migration"
    echo "  $0 test add_users_table dev            # Test migration up and down"
    echo "  $0 lint                                # Lint all migrations for issues"
    echo "  $0 check-deps                          # Check migration dependencies"
    echo ""
    echo "Prerequisites:"
echo "  - Migration tool (./migrate binary or ./cmd/migrate/main.go)"
echo "  - Environment file (.env.development, .env.staging, or .env.production)"
echo "  - Database configuration in environment file"
    echo ""
}

# Main script logic
main() {
    local command="$1"
    local arg2="$2"
    local arg3="$3"
    local migration_name=""
    local environment=""
    
    # Handle different argument patterns based on command
    case "$command" in
        "create"|"create-go")
            migration_name="$arg2"
            environment="$arg3"
            ;;
        "test")
            migration_name="$arg2"
            environment="$arg3"
            ;;
        "validate")
            if [ -n "$arg2" ] && [[ "$arg2" != "dev" ]] && [[ "$arg2" != "development" ]] && [[ "$arg2" != "staging" ]] && [[ "$arg2" != "prod" ]] && [[ "$arg2" != "production" ]]; then
                migration_name="$arg2"
                environment="$arg3"
            else
                environment="$arg2"
            fi
            ;;
        *)
            environment="$arg2"
            migration_name="$arg3"
            ;;
    esac
    
    # Check if help is requested
    if [ "$command" = "-h" ] || [ "$command" = "--help" ] || [ "$command" = "help" ] || [ -z "$command" ]; then
        show_help
        exit 0
    fi
    
    # Check prerequisites
    check_migration_tool
    
    # Commands that don't require environment
    local no_env_commands=("lint" "check-deps")
    local needs_env=true
    
    for no_env_cmd in "${no_env_commands[@]}"; do
        if [[ "$command" == "$no_env_cmd" ]]; then
            needs_env=false
            break
        fi
    done
    
    # Commands that can work without environment (validate can work on files only)
    if [[ "$command" == "validate" ]]; then
        needs_env=false
    fi
    
    # Validate arguments
    if [ "$needs_env" == true ] && [ -z "$environment" ]; then
        print_error "Environment is required for command: $command"
        show_help
        exit 1
    fi
    
    # For create commands, migration name is required
    if [[ "$command" == "create"* ]] && [ -z "$migration_name" ]; then
        print_error "Migration name is required for create commands"
        show_help
        exit 1
    fi
    
    # For test command, migration name is required
    if [[ "$command" == "test" ]] && [ -z "$migration_name" ]; then
        print_error "Migration name is required for test command"
        show_help
        exit 1
    fi
    
    # Run the command
    run_migration "$command" "$environment" "$migration_name"
    
    if [ $? -eq 0 ]; then
        print_success "Migration command completed successfully!"
    else
        print_error "Migration command failed!"
        exit 1
    fi
}

# Run main function with all arguments
main "$@"
