package store

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", "host=localhost port=5433 user=pg password=password dbname=tests sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("Error ao conectar com o banco postgres: %w", err)
	}
	return db, nil
}
