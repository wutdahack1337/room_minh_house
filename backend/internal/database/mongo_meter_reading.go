package database

import (
	"context"
	"fmt"

	types "github.com/wutdahack1337/room_minh_house/backend/pkg/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *MongoService) CreateMeterReading(ctx context.Context, meterReading *types.MeterReading) error {
	// check if room exists
	roomColl := s.database.Collection(roomCollection)
	var room types.Room
	if err := roomColl.FindOne(ctx, bson.M{"room_number": meterReading.RoomNumber}).Decode(&room); err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("database error: room with number %d not found", meterReading.RoomNumber)
		}
		return err
	}

	collection := s.database.Collection(meterReadingCollection)
	_, err := collection.InsertOne(ctx, meterReading)
	return err
}

func (s *MongoService) UpdateMeterReading(ctx context.Context, meterReading *types.MeterReading) error {
	collection := s.database.Collection(meterReadingCollection)

	res, err := collection.UpdateOne(ctx, bson.M{"room_number": meterReading.RoomNumber, "period": meterReading.Period}, bson.M{"$set": meterReading})
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (s *MongoService) DeleteMeterReadingsByRoomNumber(ctx context.Context, roomNumber int32) error {
	collection := s.database.Collection(meterReadingCollection)
	_, err := collection.DeleteMany(ctx, bson.M{"room_number": roomNumber})
	return err
}

func (s *MongoService) GetMeterReading(ctx context.Context, roomNumber int32, period string) (*types.MeterReading, error) {
	var meterReading types.MeterReading
	collection := s.database.Collection(meterReadingCollection)
	if err := collection.FindOne(ctx, bson.M{"room_number": roomNumber, "period": period}).Decode(&meterReading); err != nil {
		return nil, err
	}

	return &meterReading, nil
}

func (s *MongoService) ListMeterReadings(ctx context.Context, roomNumber int32) ([]*types.MeterReading, error) {
	collection := s.database.Collection(meterReadingCollection)

	cursor, err := collection.Find(ctx, bson.M{"room_number": roomNumber})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var res []*types.MeterReading
	for cursor.Next(ctx) {
		var meterReading types.MeterReading
		if err := cursor.Decode(&meterReading); err != nil {
			return nil, err
		}
		res = append(res, &meterReading)
	}

	return res, cursor.Err()
}
