# Workspace Service

## Overview

The Workspace Service manages collaborative workspaces with user access control and member management capabilities.

## Architecture

Following Clean Architecture principles:
```
+-----------+     +-------------+     +----------+     +-------------+     +-----------+
|  Router   | --> | Controller  | --> | Service  | <-- | Repository  | <-- | Database  |
+-----------+     +-------------+     +----------+     +-------------+     +-----------+
```

## API Endpoints

- `POST /workspaces/` - Create a new workspace
- `DELETE /workspaces/:id` - Delete a workspace by ID
- `POST /workspaces/:id/members` - Invite members to a workspace
- `POST /workspaces/:id/members/remove` - Remove members from a workspace
- `PUT /workspaces/:id/members/:member_id/access` - Change member access level

## Dependencies

- **Database**: PostgreSQL via GORM
- **HTTP Client**: For external API integrations
- **Gin**: HTTP routing and middleware

## Features

- **Workspace Management**: Create and delete workspaces
- **Member Management**: Invite and remove workspace members
- **Access Control**: Manage member permissions (read/write access)
- **Owner System**: Each workspace has an owner with full control

## Data Models

- **Workspace**: Core workspace entity with name and owner
- **Owner**: User who owns the workspace
- **Member**: User with access to the workspace
- **Access**: Permission levels (workspace:read, workspace:write)