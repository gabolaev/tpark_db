package database

import (
	"fmt"
	"os"

	"github.com/jackc/pgx"
)

// Connect is
func Connect() (*pgx.ConnPool, error) {
	connConfig, err := pgx.ParseEnvLibpq()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to parse environment:", err)
		os.Exit(1)
		return nil, err
	}

	pool, err := pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig:     connConfig,
			MaxConnections: 8,
		})

	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to create connection pool", err)
		os.Exit(1)
		return nil, err
	}
	return pool, nil
}
