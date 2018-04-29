package main

import (
	"fmt"
	"log"

	"os"
	"os/signal"
	"syscall"

	"github.com/gabolaev/tpark_db/config"
	"github.com/gabolaev/tpark_db/database"
	"github.com/gabolaev/tpark_db/logger"
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
		logger.Instance.Debug("SIGINT CATCHED")
		database.Instance.Disconnect()
		os.Exit(0)
	}()

	if err := database.Instance.LoadSchema(config.Instance.Database.SchemaFile); err != nil {
		fmt.Println(err)
		return
	}

	log.Fatal(fasthttp.ListenAndServe(":5000", router.LogHandledRequests(router.Instance.Handler)))
}
