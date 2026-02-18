package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const roomsCollection = "rooms"

type MongoService struct {
	database *mongo.Database
}

func NewMongoService(database *mongo.Database) *MongoService {
	return &MongoService{
		database: database,
	}
}

func (s *MongoService) Health(ctx context.Context) error {
	return s.database.Client().Ping(ctx, readpref.Primary())
}
