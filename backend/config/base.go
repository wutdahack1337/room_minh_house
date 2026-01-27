package config

import cli "github.com/jawher/mow.cli"

func New[T any](options ...func(*T)) *T {
	config := new(T)
	for _, o := range options {
		o(config)
	}
	return config
}

type DbOptions struct {
	DbURI  string
	DbName string
}

type DbConfig func(*DbOptions)

func WithDbOptions(c *cli.Cmd) DbConfig {
	return func(opts *DbOptions) {
		c.StringPtr(&opts.DbURI, cli.StringOpt{
			Name:   "db-uri",
			Desc:   "MongoDB uri",
			EnvVar: "DB_URI",
			Value:  "mongodb://localhost:27017/?directConnection=true",
		})
		c.StringPtr(&opts.DbName, cli.StringOpt{
			Name:   "db-name",
			Desc:   "Service db name",
			EnvVar: "DB_NAME",
			Value:  "capstone",
		})
	}

}
