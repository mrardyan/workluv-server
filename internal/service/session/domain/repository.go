package domain

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for session data access
type Repository interface {
	// Create creates a new session in the database
	Create(ctx context.Context, session Session) (Session, error)

	// FindByRefreshToken finds a session by its refresh token
	FindByRefreshToken(ctx context.Context, refreshToken string) (Session, error)

	// FindActiveByAccountID finds all active sessions for an account
	FindActiveByAccountID(ctx context.Context, accountID uuid.UUID) ([]Session, error)

	// Update updates an existing session
	Update(ctx context.Context, session Session) (Session, error)

	// RevokeByID revokes a specific session by setting revoked_at
	RevokeByID(ctx context.Context, sessionID uuid.UUID) error

	// RevokeAllByAccountID revokes all sessions for an account
	RevokeAllByAccountID(ctx context.Context, accountID uuid.UUID) error

	// DeleteExpired removes all expired sessions from the database
	DeleteExpired(ctx context.Context) error
}
