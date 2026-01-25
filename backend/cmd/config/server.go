package config

import cli "github.com/jawher/mow.cli"

type ServerOptions struct {
	HttpAddress string
}

type ServerConfig func(*ServerOptions)

func WithAllServerOptions(c *cli.Cmd) ServerConfig {
	return func(opts *ServerOptions) {
		c.StringPtr(&opts.HttpAddress, cli.StringOpt{
			Name:   "http-address",
			Desc:   "HTTP server listening address",
			EnvVar: "HTTP_ADDRESS",
			Value:  "0.0.0.0:8080",
		})
	}
}
