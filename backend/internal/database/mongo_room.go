package database

import (
	"context"

	types "github.com/wutdahack1337/room_minh_house/backend/pkg/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *MongoService) CreateRoom(ctx context.Context, room *types.Room) error {
	collection := s.database.Collection(roomsCollection)
	_, err := collection.InsertOne(ctx, room)
	return err
}

func (s *MongoService) UpdateRoom(ctx context.Context, room *types.Room) error {
	collection := s.database.Collection(roomsCollection)
	res, err := collection.UpdateOne(ctx, bson.M{"room_number": room.RoomNumber}, bson.M{"$set": room})

	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (s *MongoService) DeleteRoom(ctx context.Context, roomNumber int32) error {
	collection := s.database.Collection(roomsCollection)
	result, err := collection.DeleteOne(ctx, bson.M{"room_number": roomNumber})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (s *MongoService) GetRoom(ctx context.Context, roomNumber int32) (*types.Room, error) {
	collection := s.database.Collection(roomsCollection)

	var room types.Room
	if err := collection.FindOne(ctx, bson.M{"room_number": roomNumber}).Decode(&room); err != nil {
		return nil, err
	}

	return &room, nil
}

func (s *MongoService) ListRooms(ctx context.Context) ([]*types.Room, error) {
	collection := s.database.Collection(roomsCollection)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var res []*types.Room
	for cursor.Next(ctx) {
		var room types.Room
		if err := cursor.Decode(&room); err != nil {
			return nil, err
		}
		res = append(res, &room)
	}
	return res, cursor.Err()
}
