# Environment Templates

This directory contains environment configuration templates for different deployment environments.

## Overview

The environment templates provide a standardized way to configure the application across different environments (development, staging, production). These templates are automatically copied to the project root during setup.

## Template Files

- `env.development.template` - Development environment configuration
- `env.staging.template` - Staging environment configuration  
- `env.production.template` - Production environment configuration

## Usage

### Automatic Setup
The setup script automatically creates environment files from these templates:

```bash
./scripts/setup.sh
```

This will create:
- `.env.development` (from `env.development.template`)
- `.env.staging` (from `env.staging.template`)
- `.env.production` (from `env.production.template`)

### Manual Setup
If you need to create environment files manually:

```bash
# Copy templates to project root
cp templates/env/env.development.template .env.development
cp templates/env/env.staging.template .env.staging
cp templates/env/env.production.template .env.production
```

## Configuration Categories

### Security Critical Variables
These must be changed in production:
- `JWT_SECRET` - Secret key for JWT token signing
- `EMAIL_SMTP_USERNAME` - SMTP username for email service
- `EMAIL_SMTP_PASSWORD` - SMTP password for email service

### Environment Configuration
- `ENVIRONMENT` - Current environment (development/staging/production)
- `CLIENT_URL` - Frontend application URL

### Server Configuration
- `SERVER_HOST` - Server host address
- `SERVER_PORT` - Server port number
- `HTTP_PORT` - HTTP port number
- `SERVER_READ_TIMEOUT` - Request read timeout
- `SERVER_WRITE_TIMEOUT` - Response write timeout
- `SERVER_IDLE_TIMEOUT` - Connection idle timeout

### Database Configuration
- `DATABASE_URL` - Database connection string (managed by DigitalOcean)
- `DB_SSL_MODE` - Database SSL mode
- `DB_MAX_OPEN_CONNS` - Maximum open database connections
- `DB_MAX_IDLE_CONNS` - Maximum idle database connections
- `DB_CONN_MAX_LIFETIME` - Database connection lifetime

### Logging Configuration
- `LOG_LEVEL` - Logging level (debug/info/warn/error)
- `LOG_FORMAT` - Log format (json/text)
- `LOG_OUTPUT` - Log output destination
- `LOG_TIME_FORMAT` - Timestamp format
- `LOG_CALLER` - Include caller information

### Monitoring Configuration
- `METRICS_ENABLED` - Enable metrics collection
- `METRICS_PORT` - Metrics server port
- `HEALTH_CHECK_PATH` - Health check endpoint
- `READINESS_PATH` - Readiness probe endpoint
- `LIVENESS_PATH` - Liveness probe endpoint
- `PROMETHEUS_PATH` - Prometheus metrics endpoint

### Feature Flags
- `FEATURE_USER_REGISTRATION` - Enable user registration
- `FEATURE_EMAIL_VERIFICATION` - Enable email verification
- `FEATURE_PASSWORD_RESET` - Enable password reset
- `FEATURE_MULTI_TENANCY` - Enable multi-tenancy
- `FEATURE_AUDIT_LOGGING` - Enable audit logging

### Email Service Configuration
- `EMAIL_SERVICE_ENABLED` - Enable email service
- `EMAIL_SERVICE_PROVIDER` - Email service provider
- `EMAIL_SMTP_HOST` - SMTP server host
- `EMAIL_SMTP_PORT` - SMTP server port
- `EMAIL_SMTP_USE_TLS` - Use TLS for SMTP
- `EMAIL_SMTP_USE_SSL` - Use SSL for SMTP
- `EMAIL_FROM_ADDRESS` - Default sender email
- `EMAIL_FROM_NAME` - Default sender name
- `EMAIL_TEMPLATE_DIR` - Email template directory

### Security Configuration
- `JWT_EXPIRATION` - JWT token expiration time
- `SESSION_EXPIRATION` - Session expiration time
- `BCRYPT_COST` - Bcrypt hashing cost

### CORS Configuration
- `CORS_ALLOWED_ORIGINS` - Allowed CORS origins

### Rate Limiting
- `RATE_LIMIT_REQUESTS` - Maximum requests per window
- `RATE_LIMIT_WINDOW` - Rate limiting time window

### Redis Configuration
- `REDIS_HOST` - Redis server host
- `REDIS_PORT` - Redis server port
- `REDIS_PASSWORD` - Redis password
- `REDIS_DB` - Redis database number
- `REDIS_TIMEOUT` - Redis connection timeout
- `REDIS_POOL_SIZE` - Redis connection pool size

### Deployment Configuration
- `DO_APP_ID` - DigitalOcean App Platform app ID
- `DO_APP_NAME` - DigitalOcean App Platform app name

## Security Best Practices

### 1. Never Commit Environment Files
Environment files (`.env.*`) are automatically added to `.gitignore` to prevent accidental commits.

### 2. Use Environment-Specific Secrets
- Development: Use development-specific secrets
- Staging: Use staging-specific secrets  
- Production: Use production-specific secrets

### 3. Secure Secret Management
For production deployments:
- Use DigitalOcean App Platform secrets
- Use environment variables for sensitive data
- Never hardcode secrets in templates

### 4. Template vs Environment Files
- **Templates**: Contain placeholder values and documentation
- **Environment Files**: Contain actual configuration values

## Environment-Specific Considerations

### Development
- Uses debug logging level
- Allows localhost CORS origins
- Uses development-specific secrets
- Enables all features for testing

### Staging
- Uses info logging level
- Allows staging domain CORS origins
- Uses staging-specific secrets
- Mirrors production configuration

### Production
- Uses info logging level
- Allows production domain CORS origins
- Uses production-specific secrets
- Optimized for performance and security

## Troubleshooting

### Template Not Found
If you get "Template file not found" error:
1. Check that template files exist in `templates/env/`
2. Verify file permissions
3. Ensure setup script is run from project root

### Environment File Not Created
If environment files aren't created:
1. Run `./scripts/setup.sh` again
2. Check for existing `.env.*` files
3. Verify template file paths

### Configuration Not Applied
If configuration changes aren't applied:
1. Restart the application
2. Check environment variable loading
3. Verify configuration file paths

## Related Documentation

- [Configuration Management](../03-development/configuration/)
- [Environment Setup](../03-development/environment/)
- [Deployment Guide](../04-deployment/)
- [Security Guidelines](../06-security/)
