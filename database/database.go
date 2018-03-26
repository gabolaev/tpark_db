package database

import (
	"io/ioutil"

	"github.com/jackc/pgx"
)

// Database structure
type Database struct {
	Pool *pgx.ConnPool
}

// Instance of database
var Instance = Database{}

// Connect method for Instance
func (i Database) Connect() error {
	if connConfig, err := pgx.ParseEnvLibpq(); err != nil {
		return nil
	} else {
		if Instance.Pool, err = pgx.NewConnPool(
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
func (i Database) LoadSchema(path string) error {
	schema, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	schemaStr := string(schema)

	if _, err := i.Pool.Exec(schemaStr); err != nil {
		return err
	}
	return nil
}
