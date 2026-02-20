package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const roomCollection = "rooms"
const userCollection = "users"

type MongoService struct {
	database *mongo.Database
}

func NewMongoService(database *mongo.Database) (*MongoService, error) {
	// set room_number as unique index
	collection := database.Collection(roomCollection)
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "room_number", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return nil, err
	}

	return &MongoService{
		database: database,
	}, nil
}

func (s *MongoService) Health(ctx context.Context) error {
	return s.database.Client().Ping(ctx, readpref.Primary())
}
