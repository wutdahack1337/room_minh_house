package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	cli "github.com/jawher/mow.cli"
	"github.com/wutdahack1337/introductory_capstone/backend/config"
	"github.com/wutdahack1337/introductory_capstone/backend/server"
)

type ApiConfig struct {
	*config.ApiOptions
}

func apiCmd(c *cli.Cmd) {
	config := &ApiConfig{
		ApiOptions: config.New(config.WithApiOptions(c)),
	}

	c.Action = func() {
		fmt.Println("Hello World!")

		server.StartServer(
			"0.0.0.0:8000",
			config.ApiOptions.HttpUrl,
		)

		sigC := make(chan os.Signal, 1)
		signal.Notify(sigC, os.Interrupt, syscall.SIGTERM)
		<-sigC
	}
}
