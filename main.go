package main

import (
	"fmt"

	"github.com/gabolaev/tpark_db/database"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = database.LoadSchema(db, "/Users/gabolaev/go/src/github.com/gabolaev/tpark_db/sql-schema/create.sql")
	if err != nil {
		fmt.Println(err)
		return
	}
}
