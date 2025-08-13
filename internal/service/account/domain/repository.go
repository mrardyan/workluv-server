package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, account Account) (Account, error)
	Delete(ctx context.Context, accountID uuid.UUID) error
	FindByEmail(ctx context.Context, email Email) (Account, error)
	FindByID(ctx context.Context, id uuid.UUID) (Account, error)
	Update(ctx context.Context, account Account) (Account, error)

	// Verification
	CreateVerification(ctx context.Context, verification Verification) (Verification, error)
	FindByToken(ctx context.Context, token string) (Verification, error)
	UpdateVerification(ctx context.Context, verification Verification) (Verification, error)
	FindByAccountID(ctx context.Context, accountID uuid.UUID) (Verification, error)
}
