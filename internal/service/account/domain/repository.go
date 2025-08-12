package domain

import "github.com/google/uuid"

type Repository interface {
	Create(account Account) (Account, error)
	Delete(accountID uuid.UUID) error
}
