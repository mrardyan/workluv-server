# Session Service

This service handles user authentication sessions with JWT tokens and database storage.

## Overview

The session service provides secure authentication functionality with the following features:
- User login with email/password authentication
- JWT-based access tokens (stateless)
- Refresh tokens stored in database (stateful)
- Session revocation (logout single device)
- Mass session revocation (logout all devices)
- Session expiration handling
- Soft delete with `revoked_at` timestamps

## Architecture

Following the clean architecture pattern:

```
internal/service/session/
├── domain/           # Business logic layer
│   ├── entity.go     # Session entity with shared.Time
│   ├── repository.go # Repository interface
│   └── usecase.go    # Business logic (minimal 4 methods)
├── repository/       # Data access layer
│   └── gorm.go      # GORM implementation
├── delivery/        # Presentation layer
│   ├── controller.go # HTTP handlers
│   └── router.go    # Route definitions
├── dto/             # Data transfer objects
│   └── dto.go       # Request/response DTOs
└── app.go           # Service registration
```

## API Endpoints

### POST `/auth/login`
Authenticates user and creates session.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": 1691234567,
  "account": {
    "id": "uuid",
    "email": "user@example.com",
    "full_name": "User Name",
    "is_active": true,
    "email_verified": true
  }
}
```

### POST `/auth/logout`
Revokes current session.

**Request:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### POST `/auth/logout-all`
Revokes all sessions for user.

**Request:**
```json
{
  "account_id": "uuid"
}
```

### POST `/auth/refresh`
Generates new access token.

**Request:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": 1691234567
}
```

## Database Schema

The `sessions` table uses epoch timestamps (BIGINT) compatible with `shared.Time`:

```sql
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    refresh_token VARCHAR(255) UNIQUE NOT NULL,
    expires_at BIGINT NOT NULL,     -- epoch seconds
    revoked_at BIGINT NULL,         -- null = active, value = revoked
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);
```

## JWT Configuration

Configure via environment variables using duration format:

```env
JWT_SECRET=your-secret-key
JWT_REFRESH_SECRET=your-refresh-secret-key
JWT_EXPIRATION=24h              # Access token expiration (e.g., "24h", "30m", "2h30m")
SESSION_EXPIRATION=168h         # Refresh token expiration (e.g., "168h" = 7 days)
```

Duration formats supported: `ns`, `us`, `ms`, `s`, `m`, `h` (e.g., "1h30m", "24h", "7200s")

## Security Features

1. **Separate secrets** for access and refresh tokens
2. **Soft deletion** with `revoked_at` for audit trails
3. **Session validation** checks both JWT validity and database state
4. **Mass logout** capability for security incidents
5. **Expired session cleanup** support
6. **Password verification** using bcrypt

## Usage

The service is automatically registered in `main.go`:

```go
session.RegisterSessionService(router, db, cfg)
```

## Dependencies

- JWT: `github.com/golang-jwt/jwt/v5`
- Password hashing: `golang.org/x/crypto/bcrypt`
- Database: GORM with PostgreSQL
- Time handling: `internal/shared.Time`

## Integration with Account Service

The session service extends the account service by adding the `FindByEmail` method to support authentication. This minimal integration maintains service boundaries while enabling login functionality.
