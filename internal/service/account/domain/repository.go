package domain

import "github.com/google/uuid"

type Repository interface {
	Create(account Account) (Account, error)
	Delete(accountID uuid.UUID) error
	FindByEmail(email Email) (Account, error)
	FindByID(id uuid.UUID) (Account, error)
	Update(account Account) (Account, error)

	// Verification
	CreateVerification(verification Verification) (Verification, error)
	FindByToken(token string) (Verification, error)
	UpdateVerification(verification Verification) (Verification, error)
	FindByAccountID(accountID uuid.UUID) (Verification, error)
}
