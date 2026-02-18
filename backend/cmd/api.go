package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	cli "github.com/jawher/mow.cli"
	"github.com/wutdahack1337/room_minh_house/backend/config"
	"github.com/wutdahack1337/room_minh_house/backend/internal/database"
	"github.com/wutdahack1337/room_minh_house/backend/server"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ApiConfig struct {
	*config.ApiOptions
	*config.DbOptions
}

func apiCmd(c *cli.Cmd) {
	config := &ApiConfig{
		ApiOptions: config.New(config.WithApiOptions(c)),
		DbOptions:  config.New(config.WithDbOptions(c)),
	}

	c.Action = func() {
		db, err := connectDatabase(config.DbURI, config.DbName)
		if err != nil {
			panic(err)
		}
		defer db.Client().Disconnect(context.Background())

		dbService := database.NewMongoService(db)

		server.StartServer(
			"0.0.0.0:8000",
			config.ApiOptions.HttpUrl,
			dbService,
		)

		sigC := make(chan os.Signal, 1)
		signal.Notify(sigC, os.Interrupt, syscall.SIGTERM)
		<-sigC
	}
}

func connectDatabase(uri, dbName string) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	return client.Database(dbName), nil
}
