#!/bin/bash

# =============================================================================
# GO-SERVER - UNIFIED DEPLOYMENT SCRIPT
# =============================================================================
# 
# This script handles deployment to DigitalOcean App Platform with proper
# secrets management. It generates app specs from .env files and deploys them.
#
# Usage: ./scripts/deploy.sh [environment] [action]
# =============================================================================

set -e

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

print_secret() {
    echo -e "${PURPLE}[SECRET]${NC} $1"
}

print_header() {
    echo -e "${BLUE}=== $1 ===${NC}"
}

# Function to check if doctl is installed
check_doctl() {
    if ! command -v doctl &> /dev/null; then
        print_error "doctl is not installed. Please install it first:"
        echo ""
        echo "On macOS:"
        echo "  brew install doctl"
        echo ""
        echo "On Linux:"
        echo "  snap install doctl"
        echo "  # or download from: https://github.com/digitalocean/doctl/releases"
        echo ""
        echo "On Windows:"
        echo "  # Download from: https://github.com/digitalocean/doctl/releases"
        echo ""
        echo "After installation, run:"
        echo "  doctl auth init"
        echo ""
        exit 1
    fi
}

# Function to check if user is authenticated
check_auth() {
    if ! doctl auth list &> /dev/null; then
        print_error "Not authenticated with DigitalOcean. Please run:"
        echo ""
        echo "  doctl auth init"
        echo ""
        echo "This will open a browser window for you to authenticate."
        echo "After authentication, run this script again."
        echo ""
        exit 1
    fi
}

# Function to check if yq is available
check_yq() {
    if ! command -v yq &> /dev/null; then
        print_error "yq is not installed. Please install it first:"
        echo ""
        echo "On macOS:"
        echo "  brew install yq"
        echo ""
        echo "On Linux:"
        echo "  snap install yq"
        echo "  # or download from: https://github.com/mikefarah/yq/releases"
        echo ""
        echo "On Windows:"
        echo "  # Download from: https://github.com/mikefarah/yq/releases"
        echo ""
        exit 1
    fi
}

# Function to check all prerequisites
check_prerequisites() {
    print_header "Checking Prerequisites"
    
    print_status "Checking doctl (DigitalOcean CLI)..."
    check_doctl
    print_success "doctl is available"
    
    print_status "Checking yq (YAML processor)..."
    check_yq
    print_success "yq is available"
    
    print_status "Checking DigitalOcean authentication..."
    check_auth
    print_success "Authenticated with DigitalOcean"
    
    print_success "All prerequisites are satisfied!"
    echo ""
}

# Function to load environment variables from .env file
load_env_file() {
    local env_file="$1"
    if [ -f "$env_file" ]; then
        print_status "Loading environment variables from $env_file"
        export $(grep -v '^#' "$env_file" | xargs)
    else
        print_warning "Environment file not found: $env_file"
        return 1
    fi
}

# Function to get app ID by name
get_app_id_by_name() {
    local app_name="$1"
    doctl apps list | grep "$app_name" | awk '{print $1}' | head -1
}

# Function to generate app spec from .env file
generate_app_spec() {
    local environment="$1"
    local env_file=".env.${environment}"
    local spec_file=".do/app-${environment}.yaml"
    local template_file="templates/do/app-${environment}.template.yaml"
    
    # Handle different naming conventions
    case "$environment" in
        "dev")
            env_file=".env.development"
            template_file="templates/do/app-development.template.yaml"
            spec_file=".do/app-development.yaml"
            ;;
        "staging")
            env_file=".env.staging"
            template_file="templates/do/app-staging.template.yaml"
            spec_file=".do/app-staging.yaml"
            ;;
        "prod")
            env_file=".env.production"
            template_file="templates/do/app-production.template.yaml"
            spec_file=".do/app-production.yaml"
            ;;
    esac
    
    print_header "Generating app spec for $environment environment"
    
    # Check if .env file exists
    if [ ! -f "$env_file" ]; then
        print_error "Environment file not found: $env_file"
        print_error "Please create it from the template:"
        echo "cp templates/env/env.${environment}.template $env_file"
        exit 1
    fi
    
    # Check if template exists
    if [ ! -f "$template_file" ]; then
        print_error "Template file not found: $template_file"
        print_error "Please create a template file first"
        exit 1
    fi
    
    # Create output directory if it doesn't exist
    mkdir -p .do
    
    # Copy template to spec file
    cp "$template_file" "$spec_file"
    
    print_status "Generating app spec with environment variables..."
    
    # Create a temporary file to build the environment variables
    local temp_env_file=$(mktemp)
    local env_vars_added=false
    
    # Read .env file and collect environment variables
    while IFS= read -r line; do
        # Skip comments and empty lines
        if [[ "$line" =~ ^[[:space:]]*# ]] || [[ -z "$line" ]]; then
            continue
        fi
        
        # Extract key and value
        if [[ "$line" =~ ^([^=]+)=(.*)$ ]]; then
            local key="${BASH_REMATCH[1]}"
            local value="${BASH_REMATCH[2]}"
            
            # Skip empty values
            if [ -z "$value" ]; then
                continue
            fi
            
            # Determine scope based on the variable type for Go server
            local scope="RUN_TIME"
            case "$key" in
                # Build-time variables (can be in app spec)
                ENVIRONMENT|GO_ENV|GIN_MODE|HTTP_PORT|SERVER_HOST|SERVER_PORT|SERVER_READ_TIMEOUT|SERVER_WRITE_TIMEOUT|SERVER_IDLE_TIMEOUT|SERVER_MAX_HEADER_BYTES|CLIENT_URL|DB_SSL_MODE|DB_MAX_OPEN_CONNS|DB_MAX_IDLE_CONNS|DB_CONN_MAX_LIFETIME|LOG_LEVEL|LOG_FORMAT|LOG_OUTPUT|LOG_TIME_FORMAT|LOG_CALLER|METRICS_ENABLED|METRICS_PORT|HEALTH_CHECK_PATH|READINESS_PATH|LIVENESS_PATH|PROMETHEUS_PATH|FEATURE_USER_REGISTRATION|FEATURE_EMAIL_VERIFICATION|FEATURE_PASSWORD_RESET|FEATURE_MULTI_TENANCY|FEATURE_AUDIT_LOGGING|EMAIL_SERVICE_ENABLED|EMAIL_SERVICE_PROVIDER|EMAIL_SMTP_HOST|EMAIL_SMTP_PORT|EMAIL_SMTP_USE_TLS|EMAIL_SMTP_USE_SSL|EMAIL_FROM_ADDRESS|EMAIL_FROM_NAME|EMAIL_TEMPLATE_DIR|JWT_EXPIRATION|SESSION_EXPIRATION|BCRYPT_COST|CORS_ALLOWED_ORIGINS|RATE_LIMIT_REQUESTS|RATE_LIMIT_WINDOW)
                    scope="RUN_AND_BUILD_TIME"
                    ;;
                # Runtime-only variables (secrets and sensitive data)
                JWT_SECRET|JWT_REFRESH_SECRET|EMAIL_SMTP_USERNAME|EMAIL_SMTP_PASSWORD|DATABASE_URL|REDIS_HOST|REDIS_PORT|REDIS_PASSWORD|REDIS_DB|REDIS_TIMEOUT|REDIS_POOL_SIZE|DO_APP_ID|DO_APP_NAME|TRUSTED_PROXIES)
                    scope="RUN_TIME"
                    ;;
                # Default to runtime
                *)
                    scope="RUN_TIME"
                    ;;
            esac
            
            print_secret "Adding $key (scope: $scope)"
            
            # Add to temporary file with correct indentation
            echo "      - key: $key" >> "$temp_env_file"
            echo "        scope: $scope" >> "$temp_env_file"
            echo "        value: \"$value\"" >> "$temp_env_file"
            env_vars_added=true
        fi
    done < "$env_file"
    
    # Replace the envs section in the spec file
    if [ "$env_vars_added" = true ]; then
        # Create a new spec file with environment variables
        local new_spec_file=$(mktemp)
        
        # Use sed to replace the envs section with our generated environment variables
        # First, create a temporary file with our envs content
        local temp_envs_file=$(mktemp)
        echo "    envs:" > "$temp_envs_file"
        cat "$temp_env_file" >> "$temp_envs_file"
        
        # Use sed to replace the envs section in the template
        # This will replace from "envs:" to the next non-indented line
        sed '/^    envs:/,/^[^ ]/ {
            /^    envs:/ {
                r '"$temp_envs_file"'
                d
            }
            /^[^ ]/ !d
        }' "$template_file" > "$new_spec_file"
        
        # Replace the original file
        mv "$new_spec_file" "$spec_file"
        
        # Clean up
        rm -f "$temp_env_file" "$temp_envs_file"
    fi
    
    if [ "$env_vars_added" = true ]; then
        print_success "App spec generated successfully: $spec_file"
        return 0
    else
        print_warning "No environment variables were added"
        return 1
    fi
}

# Function to deploy app
deploy_app() {
    local environment="$1"
    local action="$2"
    local spec_file=""
    local app_name=""
    local env_file=""
    
    # Determine spec file, app name, and env file based on environment
    case "$environment" in
        "dev")
            spec_file=".do/app-development.yaml"
            app_name="workluv-dev"
            env_file=".env.development"
            ;;
        "staging")
            spec_file=".do/app-staging.yaml"
            app_name="go-server-staging"
            env_file=".env.staging"
            ;;
        "prod")
            spec_file=".do/app-production.yaml"
            app_name="go-server-prod"
            env_file=".env.production"
            ;;
        *)
            print_error "Unknown environment: $environment"
            show_help
            exit 1
            ;;
    esac
    
    print_header "Deploying to $environment environment"
    print_status "Using spec file: $spec_file"
    print_status "App name: $app_name"
    print_status "Environment file: $env_file"
    
    # Generate app spec if it doesn't exist or if explicitly requested
    if [ ! -f "$spec_file" ] || [ "$action" = "generate" ]; then
        print_status "Generating app spec from .env file..."
        generate_app_spec "$environment"
    fi
    
    case "$action" in
        "deploy"|"update")
            # Find existing app
            local app_id=$(get_app_id_by_name "$app_name")
            if [ -n "$app_id" ]; then
                print_status "Updating existing app: $app_id"
                doctl apps update "$app_id" --spec "$spec_file"
                print_success "App updated successfully!"
                
                # Verify deployment
                verify_deployment "$app_id" "$app_name"
            else
                print_error "No existing app found with name: $app_name"
                print_error "Use 'create' action to create a new app:"
                echo "$0 $environment create"
                exit 1
            fi
            ;;
        "create")
            print_status "Creating new app..."
            local app_id=$(doctl apps create --spec "$spec_file" --format ID --no-header)
            print_success "App created successfully: $app_id"
            
            # Verify deployment
            verify_deployment "$app_id" "$app_name"
            ;;
        "setup")
            # Setup environment files
            setup_env_files "$environment"
            ;;
        "generate")
            print_success "App spec generated successfully!"
            ;;
        "logs")
            local app_id=$(get_app_id_by_name "$app_name")
            if [ -n "$app_id" ]; then
                print_status "Viewing logs for app: $app_id"
                doctl apps logs "$app_id"
            else
                print_error "No app found with name: $app_name"
                exit 1
            fi
            ;;
        "status")
            local app_id=$(get_app_id_by_name "$app_name")
            if [ -n "$app_id" ]; then
                print_status "Checking status for app: $app_id"
                doctl apps get "$app_id"
            else
                print_error "No app found with name: $app_name"
                exit 1
            fi
            ;;
        *)
            print_error "Unknown action: $action"
            show_help
            exit 1
            ;;
    esac
}

# Function to show help
show_help() {
    echo "Go-Server - Unified Deployment Script"
    echo ""
    echo "Usage: $0 [environment] [action]"
    echo ""
    echo "Environments:"
    echo "  dev      - Development environment"
    echo "  staging  - Staging environment"
    echo "  prod     - Production environment"
    echo ""
    echo "Actions:"
    echo "  deploy   - Deploy/update app with generated spec"
    echo "  create   - Create new app with generated spec"
    echo "  generate - Generate app spec from .env file only"
    echo "  setup    - Setup environment files from templates"
    echo "  logs     - View app logs"
    echo "  status   - Check app status"
    echo "  list     - List all apps"
    echo "  help     - Show environment setup instructions"
    echo ""
    echo "Examples:"
    echo "  $0 dev setup      # Setup development environment files"
    echo "  $0 dev deploy     # Deploy to development"
    echo "  $0 prod create    # Create production app"
    echo "  $0 staging generate # Generate staging spec only"
    echo "  $0 dev logs       # View development logs"
    echo ""
    echo "Environment Files:"
    echo "  The script will look for .env files in the root directory:"
    echo "  - .env.development for development"
    echo "  - .env.staging for staging"
    echo "  - .env.production for production"
    echo ""
    echo "Note: Make sure your .env files contain the necessary secrets!"
    echo ""
}

# Function to list all apps
list_apps() {
    print_header "All Apps"
    doctl apps list
}

# Function to setup environment files
setup_env_files() {
    local environment="$1"
    local template_file=""
    local env_file=""
    
    case "$environment" in
        "dev")
            template_file="templates/env/env.development.template"
            env_file=".env.development"
            ;;
        "staging")
            template_file="templates/env/env.staging.template"
            env_file=".env.staging"
            ;;
        "prod")
            template_file="templates/env/env.production.template"
            env_file=".env.production"
            ;;
    esac
    
    print_header "Setting up environment files for $environment"
    
    # Check if template exists
    if [ ! -f "$template_file" ]; then
        print_error "Template file not found: $template_file"
        exit 1
    fi
    
    # Copy template to env file if it doesn't exist
    if [ ! -f "$env_file" ]; then
        print_status "Creating $env_file from template..."
        cp "$template_file" "$env_file"
        print_success "Created $env_file"
        print_warning "Please edit $env_file with your actual values before deploying"
    else
        print_status "$env_file already exists"
    fi
}

# Function to show environment setup instructions
show_env_setup() {
    print_header "Environment Setup Instructions"
    echo ""
    echo "1. Copy the environment template file:"
    echo "   cp templates/env/env.development.template .env.development"
    echo ""
    echo "2. Edit the environment file with your actual values:"
    echo "   nano .env.development"
    echo ""
    echo "3. Make sure to set these critical secrets:"
    echo "   - JWT_SECRET"
    echo "   - EMAIL_SMTP_USERNAME"
    echo "   - EMAIL_SMTP_PASSWORD"
    echo ""
    echo "4. Deploy your application:"
    echo "   ./scripts/deploy.sh dev deploy"
    echo ""
}

# Function to verify deployment
verify_deployment() {
    local app_id="$1"
    local app_name="$2"
    
    print_header "Verifying deployment for $app_name"
    
    # Wait a moment for deployment to start
    sleep 5
    
    # Check app status
    print_status "Checking app status..."
    local status=$(doctl apps get "$app_id" --format Status --no-header 2>/dev/null || echo "UNKNOWN")
    
    if [ "$status" = "RUNNING" ]; then
        print_success "App is running successfully!"
    elif [ "$status" = "BUILDING" ] || [ "$status" = "DEPLOYING" ]; then
        print_warning "App is still deploying. Status: $status"
        print_status "You can check status with: ./scripts/deploy.sh $environment status"
        print_status "View logs with: ./scripts/deploy.sh $environment logs"
    else
        print_warning "App status: $status"
        print_status "Check logs for details: ./scripts/deploy.sh $environment logs"
    fi
    
    # Show app URL if available
    local url=$(doctl apps get "$app_id" --format DefaultIngress --no-header 2>/dev/null)
    if [ -n "$url" ] && [ "$url" != "null" ]; then
        print_success "App URL: $url"
    fi
}

# Main script logic
main() {
    # Check prerequisites
    check_prerequisites
    
    # Parse command line arguments
    local environment="$1"
    local action="$2"
    
    # Show help if no arguments or help requested
    if [ -z "$environment" ] || [ "$environment" = "help" ] || [ "$environment" = "--help" ] || [ "$environment" = "-h" ]; then
        show_help
        exit 0
    fi
    
    # Handle special commands
    case "$environment" in
        "list")
            list_apps
            exit 0
            ;;
        "setup")
            show_env_setup
            exit 0
            ;;
    esac
    
    # Validate action
    if [ -z "$action" ]; then
        print_error "Action is required. Usage: $0 [environment] [action]"
        show_help
        exit 1
    fi
    
    # Deploy app
    deploy_app "$environment" "$action"
}

# Run main function with all arguments
main "$@" 