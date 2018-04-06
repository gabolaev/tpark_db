package main

import (
	"fmt"

	"github.com/gabolaev/tpark_db/config"
	"github.com/gabolaev/tpark_db/database"
	"github.com/gabolaev/tpark_db/router"
	"github.com/valyala/fasthttp"
)

func main() {
	if err := database.Instance.Connect(); err != nil {
		fmt.Println(err)
		return
	}
	defer database.Instance.Disconnect()

	if err := database.Instance.LoadSchema(config.Instance.Database.SchemaFile); err != nil {
		fmt.Println(err)
		return
	}

	fasthttp.ListenAndServe(":8080", router.Instance.Handler)
}
