# Environment Configuration Templates

This directory contains environment variable templates for different deployment environments.

## Overview

The environment templates provide a standardized way to configure your Go server application across different deployment stages. Each template includes all necessary configuration variables with appropriate defaults and documentation.

## Available Templates

- **`env.development.template`** - Development environment configuration
- **`env.staging.template`** - Staging environment configuration  
- **`env.production.template`** - Production environment configuration

## Quick Setup

### 1. Setup Environment Files

Use the setup script to create environment files from templates:

```bash
# Setup development environment
./src/scripts/setup.sh dev

# Setup staging environment
./src/scripts/setup.sh staging

# Setup production environment
./src/scripts/setup.sh prod

# Setup all environments at once
./src/scripts/setup.sh all
```

### 2. Configure Your Values

Edit the generated `.env.[environment]` files with your actual values:

```bash
# Edit development environment
nano .env.development

# Edit staging environment
nano .env.staging

# Edit production environment
nano .env.production
```

## Environment Variable Categories

### Server Configuration (Build-Time)

These variables are used during the build process and can be included in DigitalOcean app specs:

```bash
# Server Settings
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s
SERVER_IDLE_TIMEOUT=60s
SERVER_MAX_HEADER_BYTES=1048576

# Application Settings
ENVIRONMENT=development
GO_ENV=development
GIN_MODE=debug
HTTP_PORT=8080

# Feature Flags
FEATURE_USER_REGISTRATION=true
FEATURE_EMAIL_VERIFICATION=true
FEATURE_PASSWORD_RESET=true
FEATURE_MULTI_TENANCY=true
FEATURE_AUDIT_LOGGING=true

# Logging
LOG_LEVEL=debug
LOG_FORMAT=json
LOG_OUTPUT=stdout
LOG_TIME_FORMAT=2006-01-02T15:04:05Z07:00
LOG_CALLER=true

# Monitoring
METRICS_ENABLED=true
METRICS_PORT=9090
PROMETHEUS_PATH=/metrics
```

### Database Configuration (Runtime-Only)

Database connection details are managed by DigitalOcean and injected at runtime:

```bash
# DATABASE_URL is managed by DigitalOcean and injected at runtime
DATABASE_URL=${db.DATABASE_URL}

# Database Settings
DB_SSL_MODE=require
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=300s
```

### External Service Secrets (Runtime-Only)

Sensitive configuration that should never be committed to version control:

```bash
# JWT Secrets (Runtime-Only)
JWT_SECRET=your-super-secret-jwt-key-here
JWT_REFRESH_SECRET=your-super-secret-refresh-key-here
JWT_EXPIRATION=15m
SESSION_EXPIRATION=7d

# Email Configuration (Runtime-Only)
EMAIL_SERVICE_ENABLED=true
EMAIL_SERVICE_PROVIDER=smtp
EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
EMAIL_SMTP_USE_TLS=true
EMAIL_SMTP_USE_SSL=false
EMAIL_SMTP_USERNAME=your-email@gmail.com
EMAIL_SMTP_PASSWORD=your-app-password
EMAIL_FROM_ADDRESS=noreply@yourdomain.com
EMAIL_FROM_NAME=Your App Name
EMAIL_TEMPLATE_DIR=src/email

# Security
BCRYPT_COST=12
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://yourdomain.com
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m
```

### DigitalOcean Configuration (Runtime-Only)

DigitalOcean-specific settings for App Platform:

```bash
# DigitalOcean App Platform
DO_APP_ID=your-app-id
DO_APP_NAME=your-app-name
TRUSTED_PROXIES=10.0.0.0/8,172.16.0.0/12,192.168.0.0/16
```

## Environment-Specific Configurations

### Development

- **Logging**: Debug level with detailed output
- **Features**: All features enabled for testing
- **Security**: Relaxed settings for development
- **Database**: Local or development database

### Staging

- **Logging**: Info level with structured output
- **Features**: Production-like feature configuration
- **Security**: Production-like security settings
- **Database**: Staging database cluster

### Production

- **Logging**: Warning level with minimal output
- **Features**: Production feature configuration
- **Security**: Strict security settings
- **Database**: Production database cluster

## Security Considerations

### Never Commit These Files

- `.env.development`
- `.env.staging`
- `.env.production`

### Use DigitalOcean Secrets

For production deployments, use DigitalOcean's secure environment variable storage:

```bash
# In DigitalOcean App Platform
JWT_SECRET=${JWT_SECRET}  # Set as secret in DO dashboard
EMAIL_SMTP_PASSWORD=${EMAIL_SMTP_PASSWORD}  # Set as secret in DO dashboard
```

### Rotate Secrets Regularly

- JWT secrets
- Database passwords
- SMTP credentials
- API keys

## Validation

### 1. Check Required Variables

Ensure all required variables are set:

```bash
# Check for missing required variables
./src/scripts/setup.sh validate
```

### 2. Test Configuration

Test your configuration before deploying:

```bash
# Test development configuration
./src/scripts/deploy.sh dev generate

# Test staging configuration
./src/scripts/deploy.sh staging generate

# Test production configuration
./src/scripts/deploy.sh prod generate
```

## Troubleshooting

### Common Issues

1. **Missing Environment File**
   ```bash
   # Create from template
   ./src/scripts/setup.sh dev
   ```

2. **Invalid Variable Values**
   - Check variable format and syntax
   - Ensure proper escaping for special characters
   - Verify environment-specific requirements

3. **Configuration Conflicts**
   - Check for duplicate variable definitions
   - Verify environment variable precedence
   - Ensure proper scope classification

### Debug Mode

Enable debug output to troubleshoot configuration issues:

```bash
export DEBUG=true
./src/scripts/setup.sh dev
```

## Best Practices

### 1. Environment Separation

- Keep development, staging, and production completely separate
- Use different secrets for each environment
- Never share production credentials with development

### 2. Configuration Management

- Use templates for consistency
- Document all configuration changes
- Validate configuration before deployment

### 3. Security

- Never commit secrets to version control
- Use DigitalOcean's secure storage for production
- Rotate secrets regularly
- Follow principle of least privilege

## Related Documentation

- [Setup Script](../src/scripts/setup.sh) - Environment setup automation
- [Deploy Script](../src/scripts/deploy.sh) - Deployment configuration
- [Migration Script](../src/scripts/migrate.sh) - Database configuration
- [Deployment Guide](../../DEPLOYMENT.md) - Production deployment
