package config

import cli "github.com/jawher/mow.cli"

type ApiOptions struct {
	HttpUrl string
}

type ServerConfig func(*ApiOptions)

func WithApiOptions(c *cli.Cmd) ServerConfig {
	return func(opts *ApiOptions) {
		c.StringPtr(&opts.HttpUrl, cli.StringOpt{
			Name:   "http-url",
			Desc:   "Service http url",
			EnvVar: "HTTP_URL",
			Value:  "0.0.0.0:8001",
		})
	}
}
