package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email         Email     `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash  string    `gorm:"not null" json:"-"`
	FullName      string    `gorm:"size:200" json:"full_name"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	EmailVerified bool      `gorm:"default:false" json:"email_verified"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
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
