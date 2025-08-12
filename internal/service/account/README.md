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
- `DELETE /accounts/:id` - Delete an account by ID

## Dependencies

- **Database**: PostgreSQL via GORM
- **HTTP Client**: For external API integrations
- **Gin**: HTTP routing and middleware

## Features

- **Account Management**: Create and delete user accounts
- **Profile Data**: Store username, email, and password
- **Data Validation**: Input validation and error handling
- **External Integration**: Ready for external API calls (email validation, etc.)

## Data Models

- **Account**: Core account entity with username, email, and password
- **Email**: Type-safe email address
- **Password**: Secure password storage (ready for hashing)
