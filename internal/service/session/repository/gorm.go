package repository

import (
	"go-server/internal/service/session/domain"
	"go-server/internal/shared"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Database implements the session repository using GORM
type Database struct {
	DB *gorm.DB
}

// NewGormRepository creates a new GORM session repository
func NewGormRepository(db *gorm.DB) domain.Repository {
	return &Database{DB: db}
}

// Create creates a new session in the database
func (d *Database) Create(session domain.Session) (domain.Session, error) {
	if err := d.DB.Create(&session).Error; err != nil {
		return domain.Session{}, err
	}
	return session, nil
}

// FindByRefreshToken finds a session by its refresh token
func (d *Database) FindByRefreshToken(refreshToken string) (domain.Session, error) {
	var session domain.Session
	err := d.DB.Where("refresh_token = ?", refreshToken).First(&session).Error
	if err != nil {
		return domain.Session{}, err
	}
	return session, nil
}

// FindActiveByAccountID finds all active sessions for an account
func (d *Database) FindActiveByAccountID(accountID uuid.UUID) ([]domain.Session, error) {
	var sessions []domain.Session
	now := shared.Now().ToEpoch()

	err := d.DB.Where("account_id = ? AND revoked_at IS NULL AND expires_at > ?", accountID, now).
		Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

// Update updates an existing session
func (d *Database) Update(session domain.Session) (domain.Session, error) {
	session.UpdatedAt = shared.Now()
	if err := d.DB.Save(&session).Error; err != nil {
		return domain.Session{}, err
	}
	return session, nil
}

// RevokeByID revokes a specific session by setting revoked_at
func (d *Database) RevokeByID(sessionID uuid.UUID) error {
	now := shared.Now().ToEpoch()
	return d.DB.Model(&domain.Session{}).
		Where("id = ? AND revoked_at IS NULL", sessionID).
		Updates(map[string]interface{}{
			"revoked_at": now,
			"updated_at": now,
		}).Error
}

// RevokeAllByAccountID revokes all sessions for an account
func (d *Database) RevokeAllByAccountID(accountID uuid.UUID) error {
	now := shared.Now().ToEpoch()
	return d.DB.Model(&domain.Session{}).
		Where("account_id = ? AND revoked_at IS NULL", accountID).
		Updates(map[string]interface{}{
			"revoked_at": now,
			"updated_at": now,
		}).Error
}

// DeleteExpired removes all expired sessions from the database
func (d *Database) DeleteExpired() error {
	now := shared.Now().ToEpoch()
	return d.DB.Where("expires_at < ?", now).Delete(&domain.Session{}).Error
}
