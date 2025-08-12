# Go-Server Deployment Guide

This guide explains how to deploy your Go server application to DigitalOcean App Platform using the provided deployment scripts.

## Overview

The deployment system consists of:
- **Deployment Scripts**: Automated deployment to DigitalOcean App Platform
- **Environment Templates**: Pre-configured environment files for different stages
- **App Specs**: DigitalOcean App Platform configuration files
- **Multi-Environment Support**: Development, Staging, and Production deployments

## Prerequisites

### 1. Install Required Tools

```bash
# Install doctl (DigitalOcean CLI)
brew install doctl

# Install yq (YAML processor)
brew install yq
```

### 2. Authenticate with DigitalOcean

```bash
# Initialize authentication
doctl auth init

# Verify authentication
doctl auth list
```

### 3. Project Structure

Ensure your project has the following structure:
```
go-server/
├── scripts/
│   ├── deploy.sh      # Main deployment script
│   └── setup.sh       # Environment setup script
├── templates/
│   ├── do/            # DigitalOcean app templates
│   │   ├── app-development.template.yaml
│   │   ├── app-staging.template.yaml
│   │   └── app-production.template.yaml
│   └── env/           # Environment templates
│       ├── env.development.template
│       ├── env.staging.template
│       └── env.production.template
├── Dockerfile         # Container configuration
└── go.mod            # Go module file
```

## Quick Start

### 1. Test Prerequisites

First, check if you have all required tools installed:

```bash
./scripts/test-prerequisites.sh
```

This will check for:
- `doctl` (DigitalOcean CLI)
- `yq` (YAML processor)
- DigitalOcean authentication
- Git repository status
- Environment files

### 2. Setup Environment Files

```bash
# Setup development environment
./scripts/setup.sh dev

# Setup staging environment
./scripts/setup.sh staging

# Setup production environment
./scripts/setup.sh prod
```

**Note**: The setup script will automatically check all prerequisites before proceeding.

### 3. Configure Environment Variables

Edit the generated `.env.[environment]` files with your actual values:

```bash
# Edit development environment
nano .env.development

# Edit staging environment
nano .env.staging

# Edit production environment
nano .env.production
```

**Critical Secrets to Configure:**
- `JWT_SECRET` - JWT signing secret
- `JWT_REFRESH_SECRET` - JWT refresh secret
- `EMAIL_SMTP_USERNAME` - SMTP username
- `EMAIL_SMTP_PASSWORD` - SMTP password

### 4. Deploy Your Application

```bash
# Deploy to development
./scripts/deploy.sh dev deploy

# Deploy to staging
./scripts/deploy.sh staging deploy

# Deploy to production
./scripts/deploy.sh prod deploy
```

## Deployment Scripts

### Available Scripts

The deployment system includes several scripts to help you manage your deployment:

#### **Prerequisites Test Script** (`test-prerequisites.sh`)
Tests all prerequisites and provides installation instructions for missing tools.

```bash
./scripts/test-prerequisites.sh
```

**What it checks:**
- `doctl` (DigitalOcean CLI) availability
- `yq` (YAML processor) availability  
- DigitalOcean authentication status
- Git repository setup
- Environment files existence

#### **Environment Setup Script** (`setup.sh`)
Sets up environment files from templates with automatic prerequisite checking.

```bash
./scripts/setup.sh [environment]
```

**Features:**
- Automatic prerequisite validation
- Environment-specific file creation
- Step-by-step guidance

#### **GitHub Configuration Script** (`configure-github.sh`)
Configures your GitHub repository information in all DigitalOcean templates.

```bash
./scripts/configure-github.sh [github_username] [repository_name]
```

**What it does:**
- Updates all app templates with your GitHub info
- Initializes git repository if needed
- Sets up git remote origin
- Provides next steps guidance

#### **Main Deployment Script** (`deploy.sh`)
Handles all deployment operations to DigitalOcean App Platform.

```bash
./scripts/deploy.sh [environment] [action]
```

### Main Deployment Script (`deploy.sh`)

The main deployment script handles:
- Environment-specific deployments
- App spec generation from environment files
- DigitalOcean App Platform operations
- Health checks and verification

**Usage:**
```bash
./scripts/deploy.sh [environment] [action]
```

**Environments:**
- `dev` - Development environment
- `staging` - Staging environment
- `prod` - Production environment

**Actions:**
- `deploy` - Deploy/update existing app
- `create` - Create new app
- `generate` - Generate app spec only
- `logs` - View app logs
- `status` - Check app status
- `setup` - Setup environment files

**Examples:**
```bash
# Deploy to development
./scripts/deploy.sh dev deploy

# Create new production app
./scripts/deploy.sh prod create

# Generate staging spec only
./scripts/deploy.sh staging generate

# View development logs
./scripts/deploy.sh dev logs

# Check production status
./scripts/deploy.sh prod status
```

### Environment Setup Script (`setup.sh`)

The setup script helps create environment files from templates.

**Usage:**
```bash
./scripts/setup.sh [environment]
```

**Examples:**
```bash
# Setup development environment
./scripts/setup.sh dev

# Setup staging environment
./scripts/setup.sh staging

# Setup production environment
./scripts/setup.sh prod
```

## Environment Configuration

### Environment Variables

The deployment system automatically classifies environment variables:

**Build-Time Variables (RUN_AND_BUILD_TIME):**
- Server configuration (ports, timeouts)
- Logging configuration
- Feature flags
- Non-sensitive database settings
- Monitoring configuration

**Runtime Variables (RUN_TIME):**
- JWT secrets
- Database credentials
- SMTP credentials
- Redis configuration
- External service keys

### Environment-Specific Settings

#### Development
- Single instance (0.5GB RAM)
- Debug logging
- Non-production database
- Development domain

#### Staging
- Single instance (1GB RAM)
- Info logging
- Non-production database
- Staging domain

#### Production
- Multiple instances (2GB RAM each)
- Warning logging
- Production database
- Production domain

## DigitalOcean App Platform

### App Configuration

Each environment creates a DigitalOcean App with:
- **Source**: GitHub repository with auto-deploy
- **Build**: Docker-based containerization
- **Database**: Managed PostgreSQL cluster
- **Domains**: Custom domain configuration
- **Health Checks**: Application health monitoring
- **Scaling**: Environment-specific resource allocation

### Database Management

- **Development**: Non-production cluster
- **Staging**: Non-production cluster
- **Production**: Production cluster with high availability

### Domain Configuration

- **Development**: `api.dev.go-server.app`
- **Staging**: `api.staging.go-server.app`
- **Production**: `api.go-server.app`

## Monitoring and Health Checks

### Health Endpoints

The deployment expects these health endpoints:
- **Health Check**: `/health`
- **Readiness**: `/health/ready`
- **Liveness**: `/health/live`
- **Metrics**: `/metrics`

### Alerts

Configured alerts for:
- Deployment failures
- Domain failures
- Health check failures

## Troubleshooting

### Common Issues

1. **Authentication Errors**
   ```bash
   # Re-authenticate with DigitalOcean
   doctl auth init
   ```

2. **Missing Environment Files**
   ```bash
   # Create environment files
   ./scripts/setup.sh [environment]
   ```

3. **Template Not Found**
   - Ensure templates are in the correct directory structure
   - Check file permissions

4. **Deployment Failures**
   ```bash
   # Check app status
   ./scripts/deploy.sh [environment] status
   
   # View logs
   ./scripts/deploy.sh [environment] logs
   ```

### Debug Mode

Enable debug output by setting:
```bash
export DEBUG=true
./scripts/deploy.sh [environment] [action]
```

## Security Considerations

### Environment Variables
- Never commit `.env` files to version control
- Use DigitalOcean's secure environment variable storage
- Rotate secrets regularly
- Use different secrets for each environment

### Database Security
- Use SSL connections in production
- Implement connection pooling
- Monitor database access logs
- Regular security updates

### Application Security
- Implement proper CORS policies
- Use HTTPS in production
- Implement rate limiting
- Regular security audits

## Best Practices

### Deployment
1. **Test in Development First**: Always test changes in development
2. **Use Staging**: Deploy to staging before production
3. **Monitor Deployments**: Watch logs and health checks
4. **Rollback Plan**: Have a rollback strategy ready

### Configuration
1. **Environment Separation**: Keep environments completely separate
2. **Secret Management**: Use DigitalOcean's secure storage
3. **Configuration Validation**: Validate configuration before deployment
4. **Documentation**: Document all configuration changes

### Monitoring
1. **Health Checks**: Implement comprehensive health checks
2. **Logging**: Use structured logging with appropriate levels
3. **Metrics**: Collect and monitor application metrics
4. **Alerts**: Set up proper alerting for critical issues

## Support

For deployment issues:
1. Check the DigitalOcean App Platform logs
2. Verify environment variable configuration
3. Ensure all prerequisites are met
4. Check the troubleshooting section above

## Additional Resources

- [DigitalOcean App Platform Documentation](https://docs.digitalocean.com/products/app-platform/)
- [doctl CLI Documentation](https://docs.digitalocean.com/reference/doctl/)
- [Go Deployment Best Practices](https://golang.org/doc/deployment.html)
