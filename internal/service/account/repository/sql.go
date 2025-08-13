package repository

import (
	"context"
	"database/sql"
	"fmt"
	accountdomain "go-server/internal/service/account/domain"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) accountdomain.Repository {
	return &Database{DB: db}
}

func (d *Database) Create(ctx context.Context, account accountdomain.Account) (accountdomain.Account, error) {
	query := `
		INSERT INTO accounts (id, email, password_hash, full_name, email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, email, password_hash, full_name, email_verified, created_at, updated_at`

	row := d.DB.QueryRowContext(ctx, query,
		account.ID,
		string(account.Email),
		account.PasswordHash,
		account.FullName,
		account.EmailVerified,
		account.CreatedAt,
		account.UpdatedAt,
	)

	var createdAccount accountdomain.Account
	var email string
	err := row.Scan(
		&createdAccount.ID,
		&email,
		&createdAccount.PasswordHash,
		&createdAccount.FullName,
		&createdAccount.EmailVerified,
		&createdAccount.CreatedAt,
		&createdAccount.UpdatedAt,
	)
	if err != nil {
		return accountdomain.Account{}, fmt.Errorf("failed to create account: %w", err)
	}

	createdAccount.Email = accountdomain.Email(email)
	return createdAccount, nil
}

func (d *Database) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM accounts WHERE id = $1`
	result, err := d.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("account not found")
	}

	return nil
}

func (d *Database) FindByEmail(ctx context.Context, email accountdomain.Email) (accountdomain.Account, error) {
	query := `
		SELECT id, email, password_hash, full_name, email_verified, created_at, updated_at
		FROM accounts 
		WHERE email = $1`

	row := d.DB.QueryRowContext(ctx, query, string(email))

	var account accountdomain.Account
	var emailStr string
	err := row.Scan(
		&account.ID,
		&emailStr,
		&account.PasswordHash,
		&account.FullName,
		&account.EmailVerified,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return accountdomain.Account{}, fmt.Errorf("account not found")
		}
		return accountdomain.Account{}, fmt.Errorf("failed to find account by email: %w", err)
	}

	account.Email = accountdomain.Email(emailStr)
	return account, nil
}

func (d *Database) FindByID(ctx context.Context, id uuid.UUID) (accountdomain.Account, error) {
	query := `
		SELECT id, email, password_hash, full_name, email_verified, created_at, updated_at
		FROM accounts 
		WHERE id = $1`

	row := d.DB.QueryRowContext(ctx, query, id)

	var account accountdomain.Account
	var email string
	err := row.Scan(
		&account.ID,
		&email,
		&account.PasswordHash,
		&account.FullName,
		&account.EmailVerified,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return accountdomain.Account{}, fmt.Errorf("account not found")
		}
		return accountdomain.Account{}, fmt.Errorf("failed to find account by ID: %w", err)
	}

	account.Email = accountdomain.Email(email)
	return account, nil
}

func (d *Database) Update(ctx context.Context, account accountdomain.Account) (accountdomain.Account, error) {
	query := `
		UPDATE accounts 
		SET email = $2, password_hash = $3, full_name = $4, email_verified = $5, updated_at = $6
		WHERE id = $1
		RETURNING id, email, password_hash, full_name, email_verified, created_at, updated_at`

	account.UpdatedAt = time.Now()

	row := d.DB.QueryRowContext(ctx, query,
		account.ID,
		string(account.Email),
		account.PasswordHash,
		account.FullName,
		account.EmailVerified,
		account.UpdatedAt,
	)

	var updatedAccount accountdomain.Account
	var email string
	err := row.Scan(
		&updatedAccount.ID,
		&email,
		&updatedAccount.PasswordHash,
		&updatedAccount.FullName,
		&updatedAccount.EmailVerified,
		&updatedAccount.CreatedAt,
		&updatedAccount.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return accountdomain.Account{}, fmt.Errorf("account not found")
		}
		return accountdomain.Account{}, fmt.Errorf("failed to update account: %w", err)
	}

	updatedAccount.Email = accountdomain.Email(email)
	return updatedAccount, nil
}

// VerificationRepository methods

func (d *Database) CreateVerification(ctx context.Context, verification accountdomain.Verification) (accountdomain.Verification, error) {
	query := `
		INSERT INTO verifications (id, account_id, type, status, token, expires_at, completed_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, account_id, type, status, token, expires_at, completed_at, created_at, updated_at`

	row := d.DB.QueryRowContext(ctx, query,
		verification.ID,
		verification.AccountID,
		string(verification.Type),
		string(verification.Status),
		verification.Token,
		verification.ExpiresAt,
		verification.CompletedAt,
		verification.CreatedAt,
		verification.UpdatedAt,
	)

	var createdVerification accountdomain.Verification
	var verificationType, status string
	err := row.Scan(
		&createdVerification.ID,
		&createdVerification.AccountID,
		&verificationType,
		&status,
		&createdVerification.Token,
		&createdVerification.ExpiresAt,
		&createdVerification.CompletedAt,
		&createdVerification.CreatedAt,
		&createdVerification.UpdatedAt,
	)
	if err != nil {
		return accountdomain.Verification{}, fmt.Errorf("failed to create verification: %w", err)
	}

	createdVerification.Type = accountdomain.VerificationType(verificationType)
	createdVerification.Status = accountdomain.VerificationStatus(status)
	return createdVerification, nil
}

func (d *Database) FindByToken(ctx context.Context, token string) (accountdomain.Verification, error) {
	query := `
		SELECT id, account_id, type, status, token, expires_at, completed_at, created_at, updated_at
		FROM verifications 
		WHERE token = $1`

	row := d.DB.QueryRowContext(ctx, query, token)

	var verification accountdomain.Verification
	var verificationType, status string
	err := row.Scan(
		&verification.ID,
		&verification.AccountID,
		&verificationType,
		&status,
		&verification.Token,
		&verification.ExpiresAt,
		&verification.CompletedAt,
		&verification.CreatedAt,
		&verification.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return accountdomain.Verification{}, fmt.Errorf("verification not found")
		}
		return accountdomain.Verification{}, fmt.Errorf("failed to find verification by token: %w", err)
	}

	verification.Type = accountdomain.VerificationType(verificationType)
	verification.Status = accountdomain.VerificationStatus(status)
	return verification, nil
}

func (d *Database) UpdateVerification(ctx context.Context, verification accountdomain.Verification) (accountdomain.Verification, error) {
	query := `
		UPDATE verifications 
		SET type = $2, status = $3, token = $4, expires_at = $5, completed_at = $6, updated_at = $7
		WHERE id = $1
		RETURNING id, account_id, type, status, token, expires_at, completed_at, created_at, updated_at`

	verification.UpdatedAt = time.Now()

	row := d.DB.QueryRowContext(ctx, query,
		verification.ID,
		string(verification.Type),
		string(verification.Status),
		verification.Token,
		verification.ExpiresAt,
		verification.CompletedAt,
		verification.UpdatedAt,
	)

	var updatedVerification accountdomain.Verification
	var verificationType, status string
	err := row.Scan(
		&updatedVerification.ID,
		&updatedVerification.AccountID,
		&verificationType,
		&status,
		&updatedVerification.Token,
		&updatedVerification.ExpiresAt,
		&updatedVerification.CompletedAt,
		&updatedVerification.CreatedAt,
		&updatedVerification.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return accountdomain.Verification{}, fmt.Errorf("verification not found")
		}
		return accountdomain.Verification{}, fmt.Errorf("failed to update verification: %w", err)
	}

	updatedVerification.Type = accountdomain.VerificationType(verificationType)
	updatedVerification.Status = accountdomain.VerificationStatus(status)
	return updatedVerification, nil
}

func (d *Database) FindByAccountID(ctx context.Context, accountID uuid.UUID) (accountdomain.Verification, error) {
	query := `
		SELECT id, account_id, type, status, token, expires_at, completed_at, created_at, updated_at
		FROM verifications 
		WHERE account_id = $1
		ORDER BY created_at DESC
		LIMIT 1`

	row := d.DB.QueryRowContext(ctx, query, accountID)

	var verification accountdomain.Verification
	var verificationType, status string
	err := row.Scan(
		&verification.ID,
		&verification.AccountID,
		&verificationType,
		&status,
		&verification.Token,
		&verification.ExpiresAt,
		&verification.CompletedAt,
		&verification.CreatedAt,
		&verification.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return accountdomain.Verification{}, fmt.Errorf("verification not found")
		}
		return accountdomain.Verification{}, fmt.Errorf("failed to find verification by account ID: %w", err)
	}

	verification.Type = accountdomain.VerificationType(verificationType)
	verification.Status = accountdomain.VerificationStatus(status)
	return verification, nil
}
