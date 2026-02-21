package database

import (
	"context"

	types "github.com/wutdahack1337/room_minh_house/backend/pkg/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *MongoService) CreateUser(ctx context.Context, user *types.User) error {
	collection := s.database.Collection(userCollection)
	_, err := collection.InsertOne(ctx, user)
	return err
}

func (s *MongoService) UpdateUser(ctx context.Context, user *types.User) error {
	collection := s.database.Collection(userCollection)
	res, err := collection.UpdateOne(ctx, bson.M{"pid": user.Pid}, bson.M{"$set": user})
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (s *MongoService) DeleteUser(ctx context.Context, pid string) error {
	collection := s.database.Collection(userCollection)
	result, err := collection.DeleteOne(ctx, bson.M{"pid": pid})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (s *MongoService) GetUser(ctx context.Context, pid string) (*types.User, error) {
	var user types.User
	collection := s.database.Collection(userCollection)
	if err := collection.FindOne(ctx, bson.M{"pid": pid}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *MongoService) ListUsers(ctx context.Context) ([]*types.User, error) {
	collection := s.database.Collection(userCollection)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*types.User
	for cursor.Next(ctx) {
		var user types.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, cursor.Err()
}
