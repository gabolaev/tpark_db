package main

import (
	"fmt"

	"os"
	"os/signal"
	"syscall"

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
	syscallChan := make(chan os.Signal, 1)
	signal.Notify(syscallChan, syscall.SIGINT)

	go func() {
		<-syscallChan
		fmt.Println("Signal emited")
		database.Instance.Disconnect()
		os.Exit(0)
	}()

	if err := database.Instance.LoadSchema(config.Instance.Database.SchemaFile); err != nil {
		fmt.Println(err)
		return
	}

	fasthttp.ListenAndServe(":8080", router.Instance.Handler)
}
