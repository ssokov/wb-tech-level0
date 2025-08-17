package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnectToDb() (*sqlx.DB, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = user
	}
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	// Формируем URL-DSN; если пароль пуст — не добавляем ":password"
	var dsn string
	if password == "" {
		dsn = fmt.Sprintf("postgresql://%s@%s:%s/%s?sslmode=%s", user, host, port, dbname, sslmode)
	} else {
		dsn = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbname, sslmode)
	}

	// Логируем безопасно (без пароля)
	log.Printf("Connecting to DB: %s@%s:%s/%s (sslmode=%s)\n", user, host, port, dbname, sslmode)

	// Используем sqlx.Connect (open + ping) и PingContext с таймаутом
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sqlx.Connect: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("db.PingContext: %w", err)
	}

	// Настройки пула (подберите под нагрузку)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Successfully connected to database")
	return db, nil
}
