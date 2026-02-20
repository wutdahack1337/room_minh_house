package main

import (
	"os"

	cli "github.com/jawher/mow.cli"
)

func main() {
	app := cli.App("room_minh_house", "Minh's rooming house website")

	app.Command("api", "start room_minh_house apis", apiCmd)

	app.Run(os.Args)
}
