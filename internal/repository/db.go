package repository

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

// DBClient инкапсулирует подключение к БД
type DBClient struct {
	DB *sql.DB
}

// NewClient создает и инициализирует подключение
func NewClient(connString string) (*DBClient, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	client := &DBClient{DB: db}
	if err := client.ensureSchema(); err != nil {
		return nil, err
	}

	return client, nil
}

// ensureSchema проверяет и создает таблицу при необходимости
func (c *DBClient) ensureSchema() error {
	_, err := c.DB.Exec(`
		CREATE TABLE IF NOT EXISTS measurements (
			username VARCHAR(12),
			muscle VARCHAR(40),
			measure REAL,
			"date" TIMESTAMP
		)
	`)
	return err
}
