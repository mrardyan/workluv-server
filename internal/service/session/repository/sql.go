package repository

import (
	"context"
	"database/sql"
	"fmt"
	"go-server/internal/service/session/domain"
	"go-server/internal/shared"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Database implements the session repository using standard SQL
type Database struct {
	DB *sql.DB
}

// NewRepository creates a new SQL session repository
func NewRepository(db *sql.DB) domain.Repository {
	return &Database{DB: db}
}

// Create creates a new session in the database
func (d *Database) Create(ctx context.Context, session domain.Session) (domain.Session, error) {
	query := `
		INSERT INTO sessions (id, account_id, refresh_token, expires_at, revoked_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, account_id, refresh_token, expires_at, revoked_at, created_at, updated_at`

	var revokedAtValue interface{}
	if session.RevokedAt != nil {
		revokedAtValue = session.RevokedAt.ToEpoch()
	}

	row := d.DB.QueryRowContext(ctx, query,
		session.ID,
		session.AccountID,
		session.RefreshToken,
		session.ExpiresAt.ToEpoch(),
		revokedAtValue,
		session.CreatedAt.ToEpoch(),
		session.UpdatedAt.ToEpoch(),
	)

	var createdSession domain.Session
	var revokedAt sql.NullInt64
	var expiresAt, createdAt, updatedAt int64
	err := row.Scan(
		&createdSession.ID,
		&createdSession.AccountID,
		&createdSession.RefreshToken,
		&expiresAt,
		&revokedAt,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return domain.Session{}, fmt.Errorf("failed to create session: %w", err)
	}

	// Convert epoch times to shared.Time
	createdSession.ExpiresAt = shared.NewTimeFromEpoch(expiresAt)
	createdSession.CreatedAt = shared.NewTimeFromEpoch(createdAt)
	createdSession.UpdatedAt = shared.NewTimeFromEpoch(updatedAt)

	if revokedAt.Valid {
		revokedTime := shared.NewTimeFromEpoch(revokedAt.Int64)
		createdSession.RevokedAt = &revokedTime
	}

	return createdSession, nil
}

// FindByRefreshToken finds a session by its refresh token
func (d *Database) FindByRefreshToken(ctx context.Context, refreshToken string) (domain.Session, error) {
	query := `
		SELECT id, account_id, refresh_token, expires_at, revoked_at, created_at, updated_at
		FROM sessions 
		WHERE refresh_token = $1`

	row := d.DB.QueryRowContext(ctx, query, refreshToken)

	var session domain.Session
	var revokedAt sql.NullInt64
	var expiresAt, createdAt, updatedAt int64
	err := row.Scan(
		&session.ID,
		&session.AccountID,
		&session.RefreshToken,
		&expiresAt,
		&revokedAt,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Session{}, fmt.Errorf("session not found")
		}
		return domain.Session{}, fmt.Errorf("failed to find session by refresh token: %w", err)
	}

	// Convert epoch times to shared.Time
	session.ExpiresAt = shared.NewTimeFromEpoch(expiresAt)
	session.CreatedAt = shared.NewTimeFromEpoch(createdAt)
	session.UpdatedAt = shared.NewTimeFromEpoch(updatedAt)

	if revokedAt.Valid {
		revokedTime := shared.NewTimeFromEpoch(revokedAt.Int64)
		session.RevokedAt = &revokedTime
	}

	return session, nil
}

// FindActiveByAccountID finds all active sessions for an account
func (d *Database) FindActiveByAccountID(ctx context.Context, accountID uuid.UUID) ([]domain.Session, error) {
	now := shared.Now().ToEpoch()
	query := `
		SELECT id, account_id, refresh_token, expires_at, revoked_at, created_at, updated_at
		FROM sessions 
		WHERE account_id = $1 AND revoked_at IS NULL AND expires_at > $2`

	rows, err := d.DB.QueryContext(ctx, query, accountID, now)
	if err != nil {
		return nil, fmt.Errorf("failed to find active sessions: %w", err)
	}
	defer rows.Close()

	var sessions []domain.Session
	for rows.Next() {
		var session domain.Session
		var revokedAt sql.NullInt64
		var expiresAt, createdAt, updatedAt int64
		err := rows.Scan(
			&session.ID,
			&session.AccountID,
			&session.RefreshToken,
			&expiresAt,
			&revokedAt,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}

		// Convert epoch times to shared.Time
		session.ExpiresAt = shared.NewTimeFromEpoch(expiresAt)
		session.CreatedAt = shared.NewTimeFromEpoch(createdAt)
		session.UpdatedAt = shared.NewTimeFromEpoch(updatedAt)

		if revokedAt.Valid {
			revokedTime := shared.NewTimeFromEpoch(revokedAt.Int64)
			session.RevokedAt = &revokedTime
		}

		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sessions: %w", err)
	}

	return sessions, nil
}

// Update updates an existing session
func (d *Database) Update(ctx context.Context, session domain.Session) (domain.Session, error) {
	session.UpdatedAt = shared.Now()
	query := `
		UPDATE sessions 
		SET refresh_token = $2, expires_at = $3, revoked_at = $4, updated_at = $5
		WHERE id = $1
		RETURNING id, account_id, refresh_token, expires_at, revoked_at, created_at, updated_at`

	var revokedAtValue interface{}
	if session.RevokedAt != nil {
		revokedAtValue = session.RevokedAt.ToEpoch()
	}

	row := d.DB.QueryRowContext(ctx, query,
		session.ID,
		session.RefreshToken,
		session.ExpiresAt.ToEpoch(),
		revokedAtValue,
		session.UpdatedAt.ToEpoch(),
	)

	var updatedSession domain.Session
	var revokedAt sql.NullInt64
	var expiresAt, createdAt, updatedAt int64
	err := row.Scan(
		&updatedSession.ID,
		&updatedSession.AccountID,
		&updatedSession.RefreshToken,
		&expiresAt,
		&revokedAt,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Session{}, fmt.Errorf("session not found")
		}
		return domain.Session{}, fmt.Errorf("failed to update session: %w", err)
	}

	// Convert epoch times to shared.Time
	updatedSession.ExpiresAt = shared.NewTimeFromEpoch(expiresAt)
	updatedSession.CreatedAt = shared.NewTimeFromEpoch(createdAt)
	updatedSession.UpdatedAt = shared.NewTimeFromEpoch(updatedAt)

	if revokedAt.Valid {
		revokedTime := shared.NewTimeFromEpoch(revokedAt.Int64)
		updatedSession.RevokedAt = &revokedTime
	}

	return updatedSession, nil
}

// RevokeByID revokes a specific session by setting revoked_at
func (d *Database) RevokeByID(ctx context.Context, sessionID uuid.UUID) error {
	now := shared.Now().ToEpoch()
	query := `
		UPDATE sessions 
		SET revoked_at = $2, updated_at = $2
		WHERE id = $1 AND revoked_at IS NULL`

	result, err := d.DB.ExecContext(ctx, query, sessionID, now)
	if err != nil {
		return fmt.Errorf("failed to revoke session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found or already revoked")
	}

	return nil
}

// RevokeAllByAccountID revokes all sessions for an account
func (d *Database) RevokeAllByAccountID(ctx context.Context, accountID uuid.UUID) error {
	now := shared.Now().ToEpoch()
	query := `
		UPDATE sessions 
		SET revoked_at = $2, updated_at = $2
		WHERE account_id = $1 AND revoked_at IS NULL`

	_, err := d.DB.ExecContext(ctx, query, accountID, now)
	if err != nil {
		return fmt.Errorf("failed to revoke all sessions: %w", err)
	}

	return nil
}

// DeleteExpired removes all expired sessions from the database
func (d *Database) DeleteExpired(ctx context.Context) error {
	now := shared.Now().ToEpoch()
	query := `DELETE FROM sessions WHERE expires_at < $1`

	_, err := d.DB.ExecContext(ctx, query, now)
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}

	return nil
}
