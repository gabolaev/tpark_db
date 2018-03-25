package main

import (
	"fmt"

	"github.com/gabolaev/tpark_db/database"
	"github.com/gabolaev/tpark_db/router"
	"github.com/valyala/fasthttp"
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

	router := router.Instance.Handler
	fasthttp.ListenAndServe(":8080", router)

}
