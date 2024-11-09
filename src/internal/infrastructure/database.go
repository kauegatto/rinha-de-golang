package infrastructure

import (
	"fmt"
	"os"

	"github.com/jackc/pgx"
)

func NewDBConnection() (*pgx.ConnPool, error) {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	connConfig, err := pgx.ParseURI(connStr)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse uri to pgx ConnectionConfig")
	}

	poolConfig := &pgx.ConnPoolConfig{
		ConnConfig:     connConfig,
		MaxConnections: 80,
	}
	conn, err := pgx.NewConnPool(*poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	return conn, nil
}
