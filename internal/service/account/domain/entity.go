package domain

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"
	"workluv/internal/shared"

	"github.com/google/uuid"
)

type Account struct {
	ID            uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email         Email       `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash  string      `gorm:"not null" json:"-"`
	FullName      string      `gorm:"size:200" json:"full_name"`
	IsActive      bool        `gorm:"default:true" json:"is_active"`
	EmailVerified bool        `gorm:"default:false" json:"email_verified"`
	CreatedAt     shared.Time `gorm:"default:EXTRACT(epoch FROM NOW())" json:"created_at"`
	UpdatedAt     shared.Time `gorm:"default:EXTRACT(epoch FROM NOW())" json:"updated_at"`
}

// TableName overrides the table name for GORM
func (Account) TableName() string {
	return "accounts"
}

type Email string

func (e Email) IsValid() bool {
	return strings.Contains(string(e), "@")
}

type Password string

func (p Password) IsValid() bool {
	return len(p) >= 8
}

func (p Password) String() string {
	return string(p)
}

// VerificationType represents the type of verification
type VerificationType string

const (
	EmailVerification VerificationType = "email"
	PhoneVerification VerificationType = "phone"
	TwoFactorSetup    VerificationType = "2fa_setup"
)

// VerificationStatus represents the status of a verification
type VerificationStatus string

const (
	VerificationPending   VerificationStatus = "pending"
	VerificationCompleted VerificationStatus = "completed"
	VerificationExpired   VerificationStatus = "expired"
)

// Verification represents a verification request (email, phone, etc.)
type Verification struct {
	ID          uuid.UUID          `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AccountID   uuid.UUID          `gorm:"type:uuid;not null;index" json:"account_id"`
	Type        VerificationType   `gorm:"type:varchar(50);not null;index" json:"type"`
	Status      VerificationStatus `gorm:"type:varchar(20);default:'pending';index" json:"status"`
	Token       string             `gorm:"type:varchar(64);not null;uniqueIndex" json:"-"`
	ExpiresAt   shared.Time        `gorm:"not null;index" json:"expires_at"`
	CompletedAt *shared.Time       `gorm:"index" json:"completed_at,omitempty"`
	CreatedAt   shared.Time        `gorm:"default:EXTRACT(epoch FROM NOW())" json:"created_at"`
	UpdatedAt   shared.Time        `gorm:"default:EXTRACT(epoch FROM NOW())" json:"updated_at"`

	// Relationships
	Account Account `gorm:"foreignKey:AccountID;references:ID" json:"account,omitempty"`
}

// TableName overrides the table name for GORM
func (Verification) TableName() string {
	return "verifications"
}

// GenerateToken generates a new verification token
func (v *Verification) GenerateToken() error {
	// Generate 32 random bytes
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return err
	}

	// Convert to hex string
	v.Token = hex.EncodeToString(bytes)

	// Set expiration based on verification type
	switch v.Type {
	case EmailVerification:
		v.ExpiresAt = shared.NewTime(time.Now().Add(24 * time.Hour)) // 24 hours for email
	case PhoneVerification:
		v.ExpiresAt = shared.NewTime(time.Now().Add(15 * time.Minute)) // 15 minutes for SMS
	case TwoFactorSetup:
		v.ExpiresAt = shared.NewTime(time.Now().Add(1 * time.Hour)) // 1 hour for 2FA setup
	default:
		v.ExpiresAt = shared.NewTime(time.Now().Add(24 * time.Hour)) // Default 24 hours
	}

	v.Status = VerificationPending
	return nil
}

// IsTokenValid checks if the token is valid and not expired
func (v *Verification) IsTokenValid(token string) bool {
	// Check if token matches
	if v.Token != token {
		return false
	}

	// Check if verification is already completed
	if v.Status == VerificationCompleted {
		return false
	}

	// Check if token has expired
	if shared.Now().After(v.ExpiresAt) {
		return false
	}

	return true
}

// MarkAsCompleted marks the verification as completed
func (v *Verification) MarkAsCompleted() {
	v.Status = VerificationCompleted
	now := shared.Now()
	v.CompletedAt = &now
}

// IsExpired checks if the verification has expired
func (v *Verification) IsExpired() bool {
	return shared.Now().After(v.ExpiresAt) && v.Status == VerificationPending
}

// NewEmailVerification creates a new email verification for an account
func NewEmailVerification(accountID uuid.UUID) (*Verification, error) {
	verification := &Verification{
		AccountID: accountID,
		Type:      EmailVerification,
		Status:    VerificationPending,
	}

	if err := verification.GenerateToken(); err != nil {
		return nil, err
	}

	return verification, nil
}
