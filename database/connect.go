package database

import (
	"io/ioutil"

	"github.com/jackc/pgx"
)

// Instance of database
var Instance *pgx.ConnPool

// Connect method for Instance
func Connect() error {
	if connConfig, err := pgx.ParseEnvLibpq(); err != nil {
		return nil
	} else {
		if Instance, err = pgx.NewConnPool(
			pgx.ConnPoolConfig{
				ConnConfig:     connConfig,
				MaxConnections: 8,
			}); err != nil {
			return err
		}
	}
	return nil
}

// LoadSchema is
func LoadSchema(db *pgx.ConnPool, path string) error {
	schema, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	schemaStr := string(schema)

	if _, err := db.Exec(schemaStr); err != nil {
		return err
	}
	return nil
}
