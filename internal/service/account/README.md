# Account Service

## Overview

The Account Service manages user accounts with authentication and profile management capabilities.

## Architecture

Following Clean Architecture principles:
```
+-----------+     +-------------+     +----------+     +-------------+     +-----------+
|  Router   | --> | Controller  | --> | Service  | <-- | Repository  | <-- | Database  |
+-----------+     +-------------+     +----------+     +-------------+     +-----------+
```

## API Endpoints

- `POST /accounts/` - Create a new account
- `DELETE /accounts/:id` - Delete an account by UUID

## Request/Response Format

### Create Account
**Request:**
```json
{
  "full_name": "John Doe",
  "email": "john.doe@example.com",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "john.doe@example.com",
  "full_name": "John Doe",
  "is_active": true,
  "email_verified": false,
  "created_at": "2023-01-01T12:00:00Z",
  "updated_at": "2023-01-01T12:00:00Z"
}
```

## Dependencies

- **Database**: PostgreSQL via GORM
- **HTTP Client**: For external API integrations
- **Gin**: HTTP routing and middleware
- **UUID**: Google UUID package for unique identifiers
- **Bcrypt**: Password hashing with golang.org/x/crypto/bcrypt

## Features

- **Account Management**: Create and delete user accounts
- **Profile Data**: Store full name, email, and hashed password
- **Data Validation**: Input validation and error handling
- **Password Security**: Bcrypt password hashing
- **UUID Support**: Uses UUID for account identifiers
- **Email Validation**: Built-in email format validation
- **External Integration**: Ready for external API calls

## Data Models

- **Account**: Core account entity with UUID, full name, email, and hashed password
- **Email**: Type-safe email address with validation
- **Password**: Secure password validation before hashing

## Database Schema

```sql
CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(200),
    is_active BOOLEAN DEFAULT true,
    email_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```
