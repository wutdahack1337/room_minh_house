package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/InjectiveLabs/suplog"
	cli "github.com/jawher/mow.cli"
	"github.com/wutdahack1337/introductory_capstone/backend/cmd/config"
)

type Config struct {
	*config.ServerOptions
}

func serverCmd(c *cli.Cmd) {
	config := &Config{
		ServerOptions: config.New(config.WithAllServerOptions(c)),
	}

	c.Action = func() {
		fmt.Println("Hello World!")

		// init server
		httpAddress := config.HttpAddress
		fmt.Println(httpAddress)

		httpMultiplexer := http.NewServeMux()
		go func() {
			if err := http.ListenAndServe(httpAddress, httpMultiplexer); err != nil {
				log.WithError(err).Errorln("HTTP server failed")
			}
		}()

		// graceful shutdown
		sigC := make(chan os.Signal, 1)
		signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)
		<-sigC
	}
}
