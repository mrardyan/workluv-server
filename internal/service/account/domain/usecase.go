package domain

import (
	"context"
	"fmt"

	"go-server/internal/infrastructure"
	"go-server/pkg/config"

	"github.com/google/uuid"
)

type UseCase struct {
	Repo         Repository
	EmailService infrastructure.EmailService
	Config       *config.Config
}

func NewUseCase(repo Repository, emailService infrastructure.EmailService, cfg *config.Config) *UseCase {
	return &UseCase{
		Repo:         repo,
		EmailService: emailService,
		Config:       cfg,
	}
}

func (uc *UseCase) CreateAccount(ctx context.Context, account Account) (Account, error) {
	// Create account in database first
	createdAccount, err := uc.Repo.Create(ctx, account)
	if err != nil {
		return Account{}, err
	}

	// Check if there's already a verification for this account
	existingVerification, err := uc.Repo.FindByAccountID(ctx, createdAccount.ID)
	if err != nil {
		// No existing verification, create a new one
		verification, err := NewEmailVerification(createdAccount.ID)
		if err != nil {
			return Account{}, fmt.Errorf("failed to create email verification: %w", err)
		}

		// Save verification to database
		_, err = uc.Repo.CreateVerification(ctx, *verification)
		if err != nil {
			return Account{}, fmt.Errorf("failed to save verification: %w", err)
		}

		// Send verification email
		if uc.EmailService.IsEnabled() {
			if err := infrastructure.SendVerificationEmail(
				ctx,
				uc.EmailService,
				uc.Config,
				string(createdAccount.Email),
				createdAccount.FullName,
				verification.Token,
			); err != nil {
				// Log error but don't fail account creation
				fmt.Printf("Failed to send verification email: %v\n", err)
			}
		}
	} else {
		// Existing verification found, update it with new token
		if err := existingVerification.GenerateToken(); err != nil {
			return Account{}, fmt.Errorf("failed to generate verification token: %w", err)
		}

		_, err = uc.Repo.UpdateVerification(ctx, existingVerification)
		if err != nil {
			return Account{}, fmt.Errorf("failed to update verification: %w", err)
		}

		// Send verification email with the updated token
		if uc.EmailService.IsEnabled() {
			if err := infrastructure.SendVerificationEmail(
				ctx,
				uc.EmailService,
				uc.Config,
				string(createdAccount.Email),
				createdAccount.FullName,
				existingVerification.Token,
			); err != nil {
				// Log error but don't fail account creation
				fmt.Printf("Failed to send verification email: %v\n", err)
			}
		}
	}

	return createdAccount, nil
}

func (uc *UseCase) VerifyEmail(ctx context.Context, token string) error {
	// Find verification by token
	verification, err := uc.Repo.FindByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("invalid verification token")
	}

	// Check if token is valid and not expired
	if !verification.IsTokenValid(token) {
		return fmt.Errorf("verification token is invalid or expired")
	}

	// Mark verification as completed
	verification.MarkAsCompleted()

	// Update verification in database
	_, err = uc.Repo.UpdateVerification(ctx, verification)
	if err != nil {
		return fmt.Errorf("failed to update verification: %w", err)
	}

	// Update account to mark email as verified
	account, err := uc.Repo.FindByID(ctx, verification.AccountID)
	if err != nil {
		return fmt.Errorf("failed to find account: %w", err)
	}

	account.EmailVerified = true
	_, err = uc.Repo.Update(ctx, account)
	if err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	return nil
}

func (uc *UseCase) ResendVerificationEmail(ctx context.Context, email Email) error {
	// Find account by email
	account, err := uc.Repo.FindByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("account not found")
	}

	// Check if already verified
	if account.EmailVerified {
		return fmt.Errorf("email is already verified")
	}

	// Check if there's already a verification for this account
	existingVerification, err := uc.Repo.FindByAccountID(ctx, account.ID)
	if err == nil {
		// Update existing verification with new token
		if err := existingVerification.GenerateToken(); err != nil {
			return fmt.Errorf("failed to generate verification token: %w", err)
		}

		_, err = uc.Repo.UpdateVerification(ctx, existingVerification)
		if err != nil {
			return fmt.Errorf("failed to update verification: %w", err)
		}

		// Send verification email with the updated token
		if uc.EmailService.IsEnabled() {
			if err := infrastructure.SendVerificationEmail(
				ctx,
				uc.EmailService,
				uc.Config,
				string(account.Email),
				account.FullName,
				existingVerification.Token,
			); err != nil {
				return fmt.Errorf("failed to send verification email: %w", err)
			}
		}

		return nil
	}

	// No existing verification found, create a new one
	verification, err := NewEmailVerification(account.ID)
	if err != nil {
		return fmt.Errorf("failed to create email verification: %w", err)
	}

	// Save verification to database
	_, err = uc.Repo.CreateVerification(ctx, *verification)
	if err != nil {
		return fmt.Errorf("failed to save verification: %w", err)
	}

	// Send verification email
	if uc.EmailService.IsEnabled() {
		if err := infrastructure.SendVerificationEmail(
			ctx,
			uc.EmailService,
			uc.Config,
			string(account.Email),
			account.FullName,
			verification.Token,
		); err != nil {
			return fmt.Errorf("failed to send verification email: %w", err)
		}
	}

	return nil
}

func (uc *UseCase) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	return uc.Repo.Delete(ctx, id)
}
