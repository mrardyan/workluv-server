package repository

import (
	accountdomain "go-server/internal/service/account/domain"

	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewGormRepository(db *gorm.DB) accountdomain.Repository {
	return &Database{DB: db}
}

func (d *Database) Create(account accountdomain.Account) (accountdomain.Account, error) {
	return account, d.DB.Create(&account).Error
}

func (d *Database) Delete(id uint) error {
	return d.DB.Delete(&accountdomain.Account{ID: id}).Error
}
