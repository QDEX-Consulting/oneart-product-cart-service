package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // Postgres driver
)

func NewDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open error: %v", err)
	}

	// Ping to verify the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping error: %v", err)
	}

	return db, nil
}
