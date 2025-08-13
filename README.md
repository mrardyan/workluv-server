# Go Server

A Go-based web server with user management functionality and external API integration capabilities.

## Overview

This project implements a RESTful API server using:
- **Gin** for HTTP routing and middleware
- **GORM** for database operations
- **PostgreSQL** as the database
- **Clean Architecture** pattern with domain, repository, and delivery layers
- **Configurable HTTP clients** for external API integrations

## Getting Started

### Prerequisites
- Go 1.23 or higher
- Docker and Docker Compose
- PostgreSQL 15

### Running with Docker (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd go-server
   ```

2. **Start the services**
   ```bash
   docker-compose up -d
   ```

3. **Access the application**
   - API Server: http://localhost:8080
   - PostgreSQL: localhost:5432
   - pgAdmin: http://localhost:5050 (admin@example.com / admin)

### Running Locally

1. **Install dependencies**
   ```bash
   go mod download
   ```

2. **Set environment variables**
   ```bash
   export DATABASE_HOST=localhost
   export DATABASE_PORT=5432
   export DATABASE_NAME=golang_arch
   export DATABASE_USER=postgres
   export DATABASE_PASSWORD=password
   export DATABASE_SSL_MODE=disable
   export SERVER_PORT=8080
   export LOG_LEVEL=debug
   ```

3. **Start PostgreSQL**
   ```bash
   docker-compose up -d postgres
   ```

4. **Build and run**
   ```bash
   go build -o go-server ./cmd/server
   ./go-server
   ```

## Configuration

The application uses a centralized configuration system that reads from environment variables. All configuration is automatically loaded when the application starts.

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_HOST` | PostgreSQL host | `localhost` |
| `DATABASE_PORT` | PostgreSQL port | `5432` |
| `DATABASE_NAME` | Database name | `golang_arch` |
| `DATABASE_USER` | Database user | `postgres` |
| `DATABASE_PASSWORD` | Database password | `password` |
| `DATABASE_SSL_MODE` | SSL mode | `disable` |
| `SERVER_PORT` | HTTP server port | `8080` |
| `LOG_LEVEL` | Logging level | `info` |

### Configuration Structure

The configuration is organized into logical groups:
- **Database**: PostgreSQL connection settings
- **Server**: HTTP server configuration
- **Log**: Logging configuration

## API Endpoints

### User Service
- `POST /users/` - Create a new user
- `DELETE /users/:id` - Delete a user by ID

### Workspace Service
- `POST /workspaces/` - Create a new workspace
- `DELETE /workspaces/:id` - Delete a workspace by ID
- `POST /workspaces/:id/members` - Invite members to a workspace
- `POST /workspaces/:id/members/remove` - Remove members from a workspace
- `PUT /workspaces/:id/members/:member_id/access` - Change member access level

## Project Structure

```
go-server/
├── src/
│   ├── scripts/          # Deployment and utility scripts
│   │   ├── deploy.sh     # DigitalOcean deployment script
│   │   ├── setup.sh      # Environment setup script
│   │   └── migrate.sh    # Database migration script
│   └── email/            # Email templates
├── infrastructure/         # Infrastructure configuration
│   ├── db.go             # Database connection
│   ├── network.go         # HTTP client configuration and utilities
│   └── router.go          # Router setup
├── internal/              # Internal packages
│   ├── user/             # User module
│   │   ├── app.go        # User service registration
│   │   ├── delivery/     # HTTP delivery layer
│   │   ├── domain/       # Business logic and entities
│   │   ├── dto/          # Data transfer objects
│   │   └── repository/   # Data access layer
│   └── workspace/        # Workspace module
│       ├── app.go        # Workspace service registration
│       ├── delivery/     # HTTP delivery layer
│       ├── domain/       # Business logic and entities
│       ├── dto/          # Data transfer objects
│       └── repository/   # Data access layer
└── main.go               # Application entry point
```

## Architecture

The project follows Clean Architecture principles:
- **Domain Layer**: Contains business entities and interfaces
- **Repository Layer**: Handles data persistence
- **Service Layer**: Implements business logic
- **Delivery Layer**: Manages HTTP requests and responses

## Network Configuration

The application includes a configurable HTTP client system for external API integrations:

### HTTP Client Features
- **Configurable timeouts**: Default 30 seconds with customizable settings
- **Connection pooling**: Optimized for production workloads
- **Utility functions**: Pre-built methods for GET, POST, DELETE requests
- **Context support**: Full context.Context integration for cancellation and timeouts

### Usage Examples

```go
// Get default HTTP client
httpClient := infrastructure.NewDefaultHTTPClient()

// Custom configuration
config := &infrastructure.NetworkConfig{
    HTTPTimeout: 60 * time.Second,
    MaxIdleConns: 200,
    IdleConnTimeout: 120 * time.Second,
}
httpClient := infrastructure.NewHTTPClient(config)

// Make HTTP requests
ctx := context.Background()
response, err := infrastructure.Get(ctx, httpClient, "https://api.example.com/data", nil)
if err != nil {
    log.Printf("Request failed: %v", err)
    return
}

// POST with JSON body
body := map[string]interface{}{"key": "value"}
response, err = infrastructure.Post(ctx, httpClient, "https://api.example.com/create", body, nil)
```

### Common Use Cases
- **External API calls**: Third-party services, payment gateways
- **Microservice communication**: Inter-service HTTP calls
- **Webhook notifications**: Sending HTTP requests to external systems
- **File uploads**: HTTP client for file transfer services
- **Email validation**: External email verification services

## Deployment

### Quick Deployment

1. **Setup environment**
   ```bash
   ./src/scripts/setup.sh dev
   ```

2. **Deploy to DigitalOcean**
   ```bash
   ./src/scripts/deploy.sh dev create
   ```

### Available Scripts

- **`./src/scripts/setup.sh`** - Environment setup and configuration
- **`./src/scripts/deploy.sh`** - DigitalOcean App Platform deployment
- **`./src/scripts/migrate.sh`** - Database migration management

For detailed deployment instructions, see [DEPLOYMENT.md](DEPLOYMENT.md).
