package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	roomCollection         = "rooms"
	userCollection         = "users"
	contractCollection     = "contracts"
	meterReadingCollection = "meter_readings"
	invoiceCollection      = "invoices"
	invoiceItemCollection  = "invoice_items"
)

type MongoService struct {
	database *mongo.Database
}

func NewMongoService(database *mongo.Database) (*MongoService, error) {
	// set room_number as unique index
	collection := database.Collection(roomCollection)
	_, err := collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "room_number", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}

	// set pid as unique index
	collection = database.Collection(userCollection)
	_, err = collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "pid", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}

	// set contract id as unique index
	collection = database.Collection(contractCollection)
	_, err = collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "contract_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}

	// set meter reading compound index on (room_number, period)
	collection = database.Collection(meterReadingCollection)
	_, err = collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "room_number", Value: 1}, {Key: "period", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}

	// set invoice compound index on (contract_id, period)
	collection = database.Collection(invoiceCollection)
	_, err = collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "contract_id", Value: 1}, {Key: "period", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}

	// set invoice_item compound index on (contract_id, period, item_type)
	collection = database.Collection(invoiceItemCollection)
	_, err = collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "contract_id", Value: 1}, {Key: "period", Value: 1}, {Key: "item_type", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}

	return &MongoService{database: database}, nil
}

func (s *MongoService) CheckHealth(ctx context.Context) error {
	return s.database.Client().Ping(ctx, readpref.Primary())
}
