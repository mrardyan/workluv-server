# DigitalOcean App Platform Templates

This directory contains DigitalOcean App Platform configuration templates for different deployment environments. These templates define the infrastructure and configuration needed to deploy the Golang Architecture project on DigitalOcean's managed platform.

## Overview

The templates provide a complete infrastructure-as-code solution for deploying the Golang Architecture project across multiple environments:

- **Development** (`app-development.template.yaml`) - For development and testing
- **Staging** (`app-staging.template.yaml`) - For pre-production testing
- **Production** (`app-production.template.yaml`) - For production deployment

## Template Structure

Each template follows the DigitalOcean App Platform specification and includes:

### Core Components

1. **Application Service**
   - Containerized Go application
   - GitHub integration with automatic deployments
   - Health checks and monitoring
   - Environment-specific configuration

2. **Database**
   - PostgreSQL 15 managed database
   - Production-ready cluster configuration
   - Environment-specific database names

3. **Domain Configuration**
   - Custom domain setup
   - SSL certificate management
   - Environment-specific subdomains

4. **Monitoring & Alerts**
   - Deployment failure alerts
   - Domain failure alerts
   - Health check monitoring

## Environment Configurations

### Development Environment
- **Name**: `golang-arch-dev`
- **Domain**: `api.dev.golang-arch.app`
- **Database**: `golang_arch_dev`
- **Resources**: 1 instance, 0.5GB RAM
- **Auto-deploy**: Enabled (deploys on every push)

### Staging Environment
- **Name**: `golang-arch-staging`
- **Domain**: `api.staging.golang-arch.app`
- **Database**: `golang_arch_staging`
- **Resources**: 1 instance, 1GB RAM
- **Auto-deploy**: Enabled (deploys on every push)

### Production Environment
- **Name**: `golang-arch-prod`
- **Domain**: `api.golang-arch.app`
- **Database**: `golang_arch_prod`
- **Resources**: 2 instances, 2GB RAM each
- **Auto-deploy**: Disabled (manual deployment required)

## Key Features

### Infrastructure
- **Region**: Singapore (sgp) for optimal latency
- **Buildpack**: Ubuntu 22.04 stack
- **Health Checks**: HTTP endpoint at `/health/ping`
- **Load Balancing**: Automatic traffic distribution

### Security
- **SSL/TLS**: Automatic certificate management
- **Database**: Production-grade PostgreSQL clusters
- **Network**: Isolated VPC with secure connections

### Monitoring
- **Health Checks**: Automatic application health monitoring
- **Alerts**: Deployment and domain failure notifications
- **Logging**: Centralized log management

## Usage

### Prerequisites

1. **DigitalOcean Account**: Active DigitalOcean account with App Platform access
2. **GitHub Repository**: Code must be in a GitHub repository
3. **Domain**: Custom domain for production (optional for dev/staging)
4. **Database Clusters**: Pre-created PostgreSQL clusters (recommended for production)

### Deployment Steps

1. **Clone the Repository**
   ```bash
   git clone https://github.com/mrardyan/golang-arch.git
   cd golang-arch
   ```

2. **Configure Environment Variables**
   - Create environment-specific `.env` files
   - Update template files with required environment variables

3. **Deploy to DigitalOcean**
   ```bash
   # Using doctl CLI
   doctl apps create --spec templates/do/app-development.template.yaml
   
   # Or via DigitalOcean Console
   # Upload the template file through the web interface
   ```

### Environment Variables

The templates include placeholders for environment variables that should be configured:

```yaml
envs:
  # Environment variables will be added here by the generator script
  # from .env.{environment} file
```

Common environment variables to configure:
- `DATABASE_URL`: PostgreSQL connection string
- `REDIS_URL`: Redis connection string
- `JWT_SECRET`: JWT signing secret
- `LOG_LEVEL`: Application logging level
- `ENVIRONMENT`: Environment name (dev/staging/prod)

## Template Customization

### Modifying Resource Allocation

To adjust CPU and memory allocation:

```yaml
instance_size_slug: apps-s-2vcpu-2gb  # 2 vCPU, 2GB RAM
instance_count: 2                      # Number of instances
```

Available instance sizes:
- `apps-s-1vcpu-0.5gb` - 1 vCPU, 0.5GB RAM
- `apps-s-1vcpu-1gb` - 1 vCPU, 1GB RAM
- `apps-s-2vcpu-2gb` - 2 vCPU, 2GB RAM
- `apps-s-4vcpu-4gb` - 4 vCPU, 4GB RAM

### Database Configuration

To modify database settings:

```yaml
databases:
  - name: db
    engine: PG
    version: "15"
    production: true
    cluster_name: your-cluster-name
    db_name: your_database_name
    db_user: your_database_user
```

### Domain Configuration

To change domain settings:

```yaml
domains:
  - domain: api.your-domain.com
    type: PRIMARY
    zone: your-domain.com
```

## Best Practices

### Security
1. **Environment Variables**: Use DigitalOcean's encrypted environment variables for secrets
2. **Database Access**: Use connection pooling and prepared statements
3. **Network Security**: Configure proper firewall rules and VPC settings
4. **SSL/TLS**: Always use HTTPS in production

### Performance
1. **Resource Allocation**: Monitor usage and adjust instance sizes accordingly
2. **Database Optimization**: Use read replicas for read-heavy workloads
3. **Caching**: Implement Redis caching for frequently accessed data
4. **CDN**: Consider using DigitalOcean's CDN for static assets

### Monitoring
1. **Health Checks**: Ensure `/health/ping` endpoint responds quickly
2. **Logging**: Use structured logging for better observability
3. **Metrics**: Monitor application and database performance
4. **Alerts**: Configure appropriate alert thresholds

## Troubleshooting

### Common Issues

1. **Deployment Failures**
   - Check Dockerfile syntax and build context
   - Verify GitHub repository access
   - Review build logs for errors

2. **Database Connection Issues**
   - Verify database cluster is running
   - Check connection string format
   - Ensure network connectivity

3. **Domain Issues**
   - Verify DNS configuration
   - Check SSL certificate status
   - Review domain verification process

### Debugging Commands

```bash
# Check app status
doctl apps list

# View app logs
doctl apps logs <app-id>

# Get app details
doctl apps get <app-id>

# Update app configuration
doctl apps update <app-id> --spec template.yaml
```

## Cost Optimization

### Development Environment
- Use minimal resources (0.5GB RAM)
- Enable auto-deploy for rapid iteration
- Consider stopping when not in use

### Staging Environment
- Use moderate resources (1GB RAM)
- Enable auto-deploy for testing
- Monitor usage patterns

### Production Environment
- Use appropriate resources for expected load
- Disable auto-deploy for controlled releases
- Implement proper monitoring and alerting

## Maintenance

### Regular Tasks
1. **Security Updates**: Keep base images and dependencies updated
2. **Database Maintenance**: Regular backups and optimization
3. **Monitoring**: Review logs and metrics regularly
4. **Cost Review**: Monitor resource usage and optimize

### Backup Strategy
1. **Database Backups**: Use DigitalOcean's automated backups
2. **Code Backups**: Maintain GitHub repository backups
3. **Configuration Backups**: Version control all template files

## Support

For issues related to:
- **DigitalOcean App Platform**: Contact DigitalOcean support
- **Application Code**: Check GitHub issues and documentation
- **Template Configuration**: Review this README and DigitalOcean documentation

## Related Documentation

- [DigitalOcean App Platform Documentation](https://docs.digitalocean.com/products/app-platform/)
- [DigitalOcean CLI (doctl) Documentation](https://docs.digitalocean.com/reference/doctl/)
- [Project Architecture Documentation](../docs/02-architecture/)
- [Deployment Documentation](../docs/04-deployment/)
