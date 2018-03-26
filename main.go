package main

import (
	"fmt"
	"io/ioutil"

	"github.com/mailru/easyjson"

	"github.com/gabolaev/tpark_db/config"
	"github.com/gabolaev/tpark_db/database"
	"github.com/gabolaev/tpark_db/router"
	"github.com/valyala/fasthttp"
)

func main() {
	var config config.Config
	configBytes, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := easyjson.Unmarshal(configBytes, &config); err != nil {
		fmt.Println(err)
		return
	}

	if err := database.Instance.Connect(); err != nil {
		fmt.Println(err)
		return
	}

	if err := database.Instance.LoadSchema(config.Database.SchemaFile); err != nil {
		fmt.Println(err)
		return
	}

	fasthttp.ListenAndServe(":8080", router.Instance.Handler)
}
