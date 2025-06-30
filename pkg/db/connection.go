package db

import (
	"github.com/jmoiron/sqlx"
)

func ConnectToDb(dataSourceName string) (*sqlx.DB, error) {
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
