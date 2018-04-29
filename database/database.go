package database

import (
	"fmt"
	"io/ioutil"

	"github.com/gabolaev/tpark_db/config"
	"github.com/gabolaev/tpark_db/models"
	"github.com/jackc/pgx"
)

// Database structure
type Database struct {
	Pool   *pgx.ConnPool
	Status models.Status
}

// Instance of database
var Instance Database

func (i *Database) Clear() error {
	tx := StartTransaction()
	defer tx.Rollback()

	schema, err := ioutil.ReadFile(config.Instance.Database.EraseFile)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(string(schema)); err != nil {
		return err
	}
	CommitTransaction(tx)
	i.Status.Forum = 0
	i.Status.User = 0
	i.Status.Thread = 0
	i.Status.Post = 0
	return nil
}

// Connect method for Instance
func (i *Database) Connect() error {
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

func (i *Database) Disconnect() {
	fmt.Println("Disconnecting")
	i.Pool.Close()
	fmt.Println("Disconnected")
}

// LoadSchema is
func (i *Database) LoadSchema(path string) error {
	tx := StartTransaction()
	defer tx.Rollback()

	schema, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(string(schema)); err != nil {
		return err
	}
	CommitTransaction(tx)
	Instance.Status = models.Status{}
	tx = StartTransaction()
	err = tx.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&i.Status.User)
	err = tx.QueryRow(`SELECT COUNT(*) FROM forums`).Scan(&i.Status.Forum)
	err = tx.QueryRow(`SELECT COUNT(*) FROM threads`).Scan(&i.Status.Thread)
	err = tx.QueryRow(`SELECT COUNT(*) FROM posts`).Scan(&i.Status.Post)
	if err != nil {
		return err
	}
	CommitTransaction(tx)
	return nil
}

func StartTransaction() *pgx.Tx {
	tx, err := Instance.Pool.Begin()
	// error probability is so small...
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return tx
}

func CommitTransaction(tx *pgx.Tx) {
	// error probability is so small...
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		fmt.Println(err)
	}
}
