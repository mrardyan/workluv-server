package domain

type Repository interface {
	Create(account Account) (Account, error)
	Delete(accountID uint) error
}
