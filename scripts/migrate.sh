#!/bin/bash

# Database Migration Script for Workluv

set -e

echo "🗄️  Database Migration Tool"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [command] [options]"
    echo ""
    echo "Commands:"
    echo "  up         Run all pending migrations"
    echo "  down       Revert the last migration"
    echo "  down-all   Revert all migrations"
    echo "  force      Force migration version"
    echo "  status     Show migration status"
    echo "  create     Create new migration files"
    echo "  reset      Reset database (down-all + up)"
    echo ""
    echo "Options:"
    echo "  -env ENV     Environment (development, staging, production)"
    echo "  -name NAME   Migration name for create command"
    echo "  -version N   Version for force command"
    echo ""
    echo "Examples:"
    echo "  $0 up"
    echo "  $0 create -name add_new_feature"
    echo "  $0 force -version 5"
    echo "  $0 status"
    echo "  $0 reset"
}

# Default values
ENVIRONMENT="development"
ACTION=""
NAME=""
VERSION=""

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        up|down|down-all|force|status|create|reset)
            ACTION="$1"
            shift
            ;;
        -env)
            ENVIRONMENT="$2"
            shift 2
            ;;
        -name)
            NAME="$2"
            shift 2
            ;;
        -version)
            VERSION="$2"
            shift 2
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Check if action is provided
if [ -z "$ACTION" ]; then
    print_error "No action specified"
    show_usage
    exit 1
fi

# Function to check if running in container
is_container() {
    # Check if we're in a container environment
    # Method 1: Check for /.dockerenv file (Docker)
    if [ -f /.dockerenv ]; then
        return 0
    fi
    
    # Method 2: Check cgroup for container indicators
    if [ -f /proc/1/cgroup ]; then
        if grep -q "docker\|lxc\|containerd\|kubepods" /proc/1/cgroup; then
            return 0
        fi
    fi
    
    # Method 3: Check for container environment variables
    if [ -n "$KUBERNETES_SERVICE_HOST" ] || [ -n "$KUBERNETES_PORT" ]; then
        return 0
    fi
    
    # Method 4: Check for DigitalOcean App Platform environment
    if [ -n "$DIGITALOCEAN_APP_ID" ] || [ -n "$DIGITALOCEAN_APP_PLATFORM" ]; then
        return 0
    fi
    
    # Method 5: Check for common container environment variables
    if [ -n "$CONTAINER" ] || [ -n "$DOCKER_CONTAINER" ]; then
        return 0
    fi
    
    # Method 6: Check for DigitalOcean App Platform specific indicators
    if [ -n "$DIGITALOCEAN_APP_SPEC" ] || [ -n "$DIGITALOCEAN_APP_NAME" ]; then
        return 0
    fi
    
    # Method 7: Check for App Platform environment variables
    if [ -n "$APP_PLATFORM" ] || [ -n "$APP_ENVIRONMENT" ]; then
        return 0
    fi
    
    # Method 8: Check for any DigitalOcean-related environment variables
    if env | grep -q "DIGITALOCEAN\|DO_"; then
        return 0
    fi
    
    # Method 9: Force container mode via environment variable
    if [ -n "$FORCE_CONTAINER_MODE" ]; then
        return 0
    fi
    
    return 1
}

# Function to run migration command
run_migration() {
    local cmd="$1"
    print_status "Running: $cmd"
    
    if is_container; then
        print_status "Detected container environment, using pre-built binary"
        if ./migrate $cmd; then
            print_success "Migration command completed successfully"
        else
            print_error "Migration command failed"
            exit 1
        fi
    else
        print_status "Using go run for local development"
        if go run cmd/migrate/main.go $cmd; then
            print_success "Migration command completed successfully"
        else
            print_error "Migration command failed"
            exit 1
        fi
    fi
}

# Execute action
case $ACTION in
    "up")
        run_migration "up"
        ;;
    "down")
        run_migration "down"
        ;;
    "down-all")
        run_migration "down-to 0"
        ;;
    "force")
        if [ -z "$VERSION" ]; then
            print_error "Version is required for force action"
            exit 1
        fi
        run_migration "up-to $VERSION"
        ;;
    "status")
        run_migration "status"
        ;;
    "create")
        if [ -z "$NAME" ]; then
            print_error "Name is required for create action"
            exit 1
        fi
        run_migration "create $NAME"
        ;;
    "reset")
        print_warning "This will reset the database (revert all + run all migrations)"
        read -p "Are you sure? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            print_status "Resetting database..."
            run_migration "down-to 0"
            run_migration "up"
            print_success "Database reset completed"
        else
            print_status "Reset cancelled"
        fi
        ;;
    *)
        print_error "Unknown action: $ACTION"
        show_usage
        exit 1
        ;;
esac
