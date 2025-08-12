#!/bin/bash

# =============================================================================
# GO-SERVER - COMPREHENSIVE ENVIRONMENT SETUP SCRIPT
# =============================================================================
# 
# This script helps set up environment files for deployment and configures
# GitHub repository integration for DigitalOcean App Platform.
#
# Usage: ./scripts/setup.sh [environment] [--github USERNAME REPO]
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

print_github() {
    echo -e "${PURPLE}[GITHUB]${NC} $1"
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

# Function to check if yq is installed
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

# Function to check if goose is installed
check_goose() {
    if ! command -v goose &> /dev/null; then
        print_error "goose is not installed. Please install it first:"
        echo ""
        echo "Install via Go:"
        echo "  go install github.com/pressly/goose/v3/cmd/goose@latest"
        echo ""
        echo "On macOS:"
        echo "  brew install goose"
        echo ""
        echo "On Linux:"
        echo "  # Download from: https://github.com/pressly/goose/releases"
        echo ""
        echo "On Windows:"
        echo "  # Download from: https://github.com/pressly/goose/releases"
        echo ""
        echo "After installation, make sure your GOPATH/bin is in your PATH"
        echo ""
        exit 1
    fi
}

# Function to check if user is authenticated with DigitalOcean
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

# Function to check prerequisites
check_prerequisites() {
    print_header "Checking Prerequisites"
    
    print_status "Checking doctl (DigitalOcean CLI)..."
    check_doctl
    print_success "doctl is available"
    
    print_status "Checking yq (YAML processor)..."
    check_yq
    print_success "yq is available"
    
    print_status "Checking goose (database migration tool)..."
    check_goose
    print_success "goose is available"
    
    print_status "Checking DigitalOcean authentication..."
    check_auth
    print_success "Authenticated with DigitalOcean"
    
    print_success "All prerequisites are satisfied!"
    echo ""
}

# Function to update DigitalOcean app templates with GitHub info
update_github_templates() {
    local github_username="$1"
    local repository_name="$2"
    local repo_path="${github_username}/${repository_name}"
    
    print_header "Updating DigitalOcean App Templates"
    
    # Update development template
    if [ -f "templates/do/app-development.template.yaml" ]; then
        print_github "Updating development template..."
        sed -i.bak "s/YOUR_GITHUB_USERNAME\/YOUR_REPO_NAME/$repo_path/g" "templates/do/app-development.template.yaml"
        rm -f "templates/do/app-development.template.yaml.bak"
        print_success "Updated development template"
    fi
    
    # Update staging template
    if [ -f "templates/do/app-staging.template.yaml" ]; then
        print_github "Updating staging template..."
        sed -i.bak "s/YOUR_GITHUB_USERNAME\/YOUR_REPO_NAME/$repo_path/g" "templates/do/app-staging.template.yaml"
        rm -f "templates/do/app-staging.template.yaml.bak"
        print_success "Updated staging template"
    fi
    
    # Update production template
    if [ -f "templates/do/app-production.template.yaml" ]; then
        print_github "Updating production template..."
        sed -i.bak "s/YOUR_GITHUB_USERNAME\/YOUR_REPO_NAME/$repo_path/g" "templates/do/app-production.template.yaml"
        rm -f "templates/do/app-production.template.yaml.bak"
        print_success "Updated production template"
    fi
}

# Function to setup git repository
setup_git_repository() {
    local github_username="$1"
    local repository_name="$2"
    local repo_path="${github_username}/${repository_name}"
    
    print_header "Setting up Git Repository"
    
    # Check if git is already initialized
    if [ -d ".git" ]; then
        print_status "Git repository already exists"
        
        # Check if remote origin is set
        if git remote get-url origin &>/dev/null; then
            local current_remote=$(git remote get-url origin)
            print_status "Current remote origin: $current_remote"
            
            if [[ "$current_remote" != *"$repo_path"* ]]; then
                print_warning "Remote origin doesn't match expected repository"
                print_status "Updating remote origin..."
                git remote set-url origin "https://github.com/$repo_path.git"
                print_success "Updated remote origin"
            else
                print_success "Remote origin is already correct"
            fi
        else
            print_status "Adding remote origin..."
            git remote add origin "https://github.com/$repo_path.git"
            print_success "Added remote origin"
        fi
    else
        print_status "Initializing git repository..."
        git init
        git add .
        git commit -m "Initial commit"
        
        print_status "Adding remote origin..."
        git remote add origin "https://github.com/$repo_path.git"
        print_success "Initialized git repository"
    fi
}

# Function to setup environment files
setup_environment() {
    local environment="$1"
    local template_file=""
    local env_file=""
    
    case "$environment" in
        "dev"|"development")
            template_file="templates/env/env.development.template"
            env_file=".env.development"
            ;;
        "staging")
            template_file="templates/env/env.staging.template"
            env_file=".env.staging"
            ;;
        "prod"|"production")
            template_file="templates/env/env.production.template"
            env_file=".env.production"
            ;;
        *)
            print_error "Unknown environment: $environment"
            show_help
            exit 1
            ;;
    esac
    
    print_header "Setting up $environment environment"
    
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
        print_warning "If you want to reset it, delete the file and run this script again"
    fi
}

# Function to setup all environments
setup_all_environments() {
    print_header "Setting up all environments"
    
    setup_environment "dev"
    setup_environment "staging"
    setup_environment "prod"
    
    print_success "All environments have been set up!"
}

# Function to show GitHub next steps
show_github_next_steps() {
    local github_username="$1"
    local repository_name="$2"
    
    print_header "GitHub Setup Next Steps"
    echo ""
    echo "1. Create your GitHub repository:"
    echo "   - Go to https://github.com/new"
    echo "   - Repository name: $repository_name"
    echo "   - Make it public or private as needed"
    echo "   - Don't initialize with README, .gitignore, or license"
    echo ""
    echo "2. Push your code to GitHub:"
    echo "   git push -u origin main"
    echo ""
    echo "3. Deploy to development:"
    echo "   ./scripts/deploy.sh dev create"
    echo ""
}

# Function to show environment setup next steps
show_env_next_steps() {
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
        "all")
            env_file="all environment files"
            ;;
    esac
    
    print_header "Next Steps"
    echo ""
    echo "1. Edit your environment file(s):"
    if [ "$environment" = "all" ]; then
        echo "   nano .env.development"
        echo "   nano .env.staging"
        echo "   nano .env.production"
    else
        echo "   nano $env_file"
    fi
    echo ""
    echo "2. Make sure to set these critical secrets:"
    echo "   - JWT_SECRET"
    echo "   - JWT_REFRESH_SECRET"
    echo "   - EMAIL_SMTP_USERNAME"
    echo "   - EMAIL_SMTP_PASSWORD"
    echo ""
    echo "3. Deploy your application:"
    if [ "$environment" = "all" ]; then
        echo "   ./scripts/deploy.sh dev create    # Start with development"
    else
        echo "   ./scripts/deploy.sh $environment deploy"
    fi
    echo ""
}

# Function to show help
show_help() {
    echo "Go-Server - Comprehensive Environment Setup Script"
    echo ""
    echo "Usage: $0 [environment] [--github USERNAME REPO]"
    echo ""
    echo "Environments:"
    echo "  dev/development - Development environment"
    echo "  staging        - Staging environment"
    echo "  prod/production - Production environment"
    echo "  all            - All environments"
    echo ""
    echo "GitHub Integration:"
    echo "  --github USERNAME REPO  - Configure GitHub repository integration"
    echo ""
    echo "Examples:"
    echo "  $0 dev                    # Setup development environment only"
    echo "  $0 all                    # Setup all environments"
    echo "  $0 dev --github ardyan go-server  # Setup dev + GitHub config"
    echo "  $0 all --github ardyan go-server  # Setup all + GitHub config"
    echo ""
    echo "This script will:"
    echo "  1. Check all prerequisites (doctl, yq, authentication)"
    echo "  2. Copy the appropriate environment template(s)"
    echo "  3. Create .env.[environment] file(s)"
    echo "  4. Configure GitHub integration (if --github specified)"
    echo "  5. Provide guidance for next steps"
    echo ""
    echo "Prerequisites:"
    echo "  - doctl (DigitalOcean CLI)"
    echo "  - yq (YAML processor)"
    echo "  - goose (database migration tool)"
    echo "  - DigitalOcean authentication"
    echo ""
}

# Main script logic
main() {
    local environment=""
    local github_username=""
    local repository_name=""
    local setup_github=false
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --github)
                if [ -n "$2" ] && [ -n "$3" ]; then
                    github_username="$2"
                    repository_name="$3"
                    setup_github=true
                    shift 3
                else
                    print_error "--github requires both username and repository name"
                    show_help
                    exit 1
                fi
                ;;
            -h|--help|help)
                show_help
                exit 0
                ;;
            dev|development|staging|prod|production|all)
                if [ -z "$environment" ]; then
                    environment="$1"
                else
                    print_error "Only one environment can be specified"
                    show_help
                    exit 1
                fi
                shift
                ;;
            *)
                print_error "Unknown argument: $1"
                show_help
                exit 1
                ;;
        esac
    done
    
    # Default to 'all' if no environment specified
    if [ -z "$environment" ]; then
        environment="all"
    fi
    
    # Check prerequisites first
    check_prerequisites
    
    # Setup GitHub integration if requested
    if [ "$setup_github" = true ]; then
        # Validate GitHub username format
        if [[ ! "$github_username" =~ ^[a-zA-Z0-9-]+$ ]]; then
            print_error "Invalid GitHub username format: $github_username"
            print_error "GitHub usernames can only contain letters, numbers, and hyphens"
            exit 1
        fi
        
        # Validate repository name format
        if [[ ! "$repository_name" =~ ^[a-zA-Z0-9_-]+$ ]]; then
            print_error "Invalid repository name format: $repository_name"
            print_error "Repository names can only contain letters, numbers, hyphens, and underscores"
            exit 1
        fi
        
        print_header "GitHub Integration Setup"
        echo ""
        echo "GitHub Username: $github_username"
        echo "Repository Name: $repository_name"
        echo "Full Path: $github_username/$repository_name"
        echo ""
        
        # Update templates and setup git
        update_github_templates "$github_username" "$repository_name"
        setup_git_repository "$github_username" "$repository_name"
        show_github_next_steps "$github_username" "$repository_name"
    fi
    
    # Setup environment(s)
    if [ "$environment" = "all" ]; then
        setup_all_environments
    else
        setup_environment "$environment"
    fi
    
    # Show next steps
    show_env_next_steps "$environment"
    
    print_success "Setup completed successfully!"
}

# Run main function with all arguments
main "$@" 