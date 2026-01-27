package main

import (
	"bufio"
	"bytes"
	"os"
	"strings"

	log "github.com/InjectiveLabs/suplog"
	cli "github.com/jawher/mow.cli"
)

func main() {
	readEnv()
	app := cli.App("capstone", "Minh's rooming house website")

	app.Command("api", "start Minh's rooming house api", apiCmd)

	app.Run(os.Args)
}

func readEnv() {
	if envData, _ := os.ReadFile(".env"); len(envData) > 0 {
		s := bufio.NewScanner(bytes.NewReader(envData))
		for s.Scan() {
			txt := s.Text()
			valIdx := strings.IndexByte(txt, '=')
			if valIdx < 0 {
				continue
			}

			strValue := strings.Trim(txt[valIdx+1:], `"`)
			if err := os.Setenv(txt[:valIdx], strValue); err != nil {
				log.WithField("name", txt[:valIdx]).WithError(err).Warningln("failed to override ENV variable")
			}
		}
	}
}
