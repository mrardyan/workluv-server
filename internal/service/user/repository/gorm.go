package repository

import (
	userdomain "go-server/internal/service/user/domain"

	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewGormRepository(db *gorm.DB) userdomain.Repository {
	return &Database{DB: db}
}

func (d *Database) Create(user userdomain.User) (userdomain.User, error) {
	return user, d.DB.Create(&user).Error
}

func (d *Database) Delete(id uint) error {
	return d.DB.Delete(&userdomain.User{ID: id}).Error
}
