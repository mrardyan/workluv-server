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
