package database

import (
	"io/ioutil"

	"github.com/jackc/pgx"
)

// Connect is
func Connect() (*pgx.ConnPool, error) {
	connConfig, err := pgx.ParseEnvLibpq()
	if err != nil {
		return nil, err
	}

	db, err := pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig:     connConfig,
			MaxConnections: 8,
		})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// LoadSchema is
func LoadSchema(db *pgx.ConnPool, path string) error {
	schema, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	schemaStr := string(schema)

	_, err = db.Exec(schemaStr)
	if err != nil {
		return err
	}
	return nil
}
