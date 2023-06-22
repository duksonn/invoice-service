package main

import (
	"fmt"
	"invoice-service/cmd/config"
	"invoice-service/cmd/dependencies"
	"invoice-service/cmd/http"
	"os"
	"os/signal"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("main: can't read config: %s", err.Error()))
	}

	var dependencies = dependencies.Init(cfg)
	fmt.Println("dependencies init successfully")

	err = http.StartServer(cfg, dependencies)

	// Wait for terminate signal to shut down server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
