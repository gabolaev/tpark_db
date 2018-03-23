package database

import (
	"fmt"
	"os"

	"github.com/jackc/pgx"
)

// Connect is
func Connect() {
	config, err := pgx.ParseEnvLibpq()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to parse environment:", err)
		os.Exit(1)
	} else {
		fmt.Println("SUCCESS")
		fmt.Println(config)
	}
}
