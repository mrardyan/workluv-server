package domain

import (
	"github.com/google/uuid"
)

// Repository defines the interface for session data access
type Repository interface {
	// Create creates a new session in the database
	Create(session Session) (Session, error)

	// FindByRefreshToken finds a session by its refresh token
	FindByRefreshToken(refreshToken string) (Session, error)

	// FindActiveByAccountID finds all active sessions for an account
	FindActiveByAccountID(accountID uuid.UUID) ([]Session, error)

	// Update updates an existing session
	Update(session Session) (Session, error)

	// RevokeByID revokes a specific session by setting revoked_at
	RevokeByID(sessionID uuid.UUID) error

	// RevokeAllByAccountID revokes all sessions for an account
	RevokeAllByAccountID(accountID uuid.UUID) error

	// DeleteExpired removes all expired sessions from the database
	DeleteExpired() error
}
