package repository

import (
	accountdomain "go-server/internal/service/account/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewGormRepository(db *gorm.DB) accountdomain.Repository {
	return &Database{DB: db}
}

func (d *Database) Create(account accountdomain.Account) (accountdomain.Account, error) {
	if err := d.DB.Create(&account).Error; err != nil {
		return accountdomain.Account{}, err
	}
	return account, nil
}

func (d *Database) Delete(id uuid.UUID) error {
	return d.DB.Delete(&accountdomain.Account{}, id).Error
}

func (d *Database) FindByEmail(email accountdomain.Email) (accountdomain.Account, error) {
	var account accountdomain.Account
	err := d.DB.Where("email = ?", string(email)).First(&account).Error
	if err != nil {
		return accountdomain.Account{}, err
	}
	return account, nil
}

func (d *Database) FindByID(id uuid.UUID) (accountdomain.Account, error) {
	var account accountdomain.Account
	err := d.DB.Where("id = ?", id).First(&account).Error
	if err != nil {
		return accountdomain.Account{}, err
	}
	return account, nil
}

func (d *Database) Update(account accountdomain.Account) (accountdomain.Account, error) {
	if err := d.DB.Save(&account).Error; err != nil {
		return accountdomain.Account{}, err
	}
	return account, nil
}

// VerificationRepository methods

func (d *Database) CreateVerification(verification accountdomain.Verification) (accountdomain.Verification, error) {
	if err := d.DB.Create(&verification).Error; err != nil {
		return accountdomain.Verification{}, err
	}
	return verification, nil
}

func (d *Database) FindByToken(token string) (accountdomain.Verification, error) {
	var verification accountdomain.Verification
	err := d.DB.Where("token = ?", token).First(&verification).Error
	if err != nil {
		return accountdomain.Verification{}, err
	}
	return verification, nil
}

func (d *Database) UpdateVerification(verification accountdomain.Verification) (accountdomain.Verification, error) {
	if err := d.DB.Save(&verification).Error; err != nil {
		return accountdomain.Verification{}, err
	}
	return verification, nil
}

func (d *Database) FindByAccountID(accountID uuid.UUID) (accountdomain.Verification, error) {
	var verification accountdomain.Verification
	err := d.DB.Where("account_id = ?", accountID).First(&verification).Error
	if err != nil {
		return accountdomain.Verification{}, err
	}
	return verification, nil
}
