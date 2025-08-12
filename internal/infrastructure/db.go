package infrastructure

import (
	"go-server/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(cfg *config.Config) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(cfg.Database.URL), &gorm.Config{})
}
