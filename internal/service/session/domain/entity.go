package domain

import (
	"go-server/internal/shared"

	"github.com/google/uuid"
)

// Session represents a user session stored in the database
type Session struct {
	ID           uuid.UUID    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AccountID    uuid.UUID    `gorm:"type:uuid;not null;index" json:"account_id"`
	RefreshToken string       `gorm:"size:255;uniqueIndex;not null" json:"-"`
	ExpiresAt    shared.Time  `gorm:"not null" json:"expires_at"`
	RevokedAt    *shared.Time `gorm:"default:null;index" json:"revoked_at,omitempty"`
	CreatedAt    shared.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt    shared.Time  `gorm:"not null" json:"updated_at"`
}

// TableName overrides the table name for GORM
func (Session) TableName() string {
	return "sessions"
}

// IsActive checks if the session is still active (not revoked and not expired)
func (s Session) IsActive() bool {
	now := shared.Now()
	return s.RevokedAt == nil && s.ExpiresAt.After(now)
}

// Revoke marks the session as revoked by setting RevokedAt to current time
func (s *Session) Revoke() {
	now := shared.Now()
	s.RevokedAt = &now
	s.UpdatedAt = now
}
