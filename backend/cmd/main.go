package main

import (
	"os"

	log "github.com/InjectiveLabs/suplog"
	cli "github.com/jawher/mow.cli"
	"google.golang.org/grpc"
)

func main() {
	app := cli.App("capstone", "Minh's rooming house website")

	app.Command("server", "start Minh's rooming house server", serverCmd)

	app.Run(os.Args)
}

func orShutdown(err error) {
	if err != nil && err != grpc.ErrServerStopped {
		log.WithError(err).Fatalln("unable to start")
	}
}
