package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kimvlry/simple-order-service/internal/config"
	_ "github.com/lib/pq"
)

func ConnectToDb(cfg config.DatabaseConfig) (*sqlx.DB, error) {
	dataSourceName := getPostgresDSN(cfg)
	db, err := sqlx.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getPostgresDSN(cfg config.DatabaseConfig) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)
}
