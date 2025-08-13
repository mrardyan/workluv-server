package infrastructure

import (
	"database/sql"
	"go-server/pkg/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectDB(cfg *config.Config) (*sql.DB, error) {
	return sql.Open("pgx", cfg.Database.URL)
}
