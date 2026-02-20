package database

import (
	"context"
	"fmt"

	types "github.com/wutdahack1337/room_minh_house/backend/pkg/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *MongoService) CreateContract(ctx context.Context, contract *types.Contract) error {
	// check if user exists
	userCollection := s.database.Collection(userCollection)
	var user types.User
	if err := userCollection.FindOne(ctx, bson.M{"pid": contract.UserPid}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("database error: user with pid %s not found", contract.UserPid)
		}
		return err
	}

	// check if room exists
	roomCollection := s.database.Collection(roomCollection)
	var room types.Room
	if err := roomCollection.FindOne(ctx, bson.M{"room_number": contract.RoomNumber}).Decode(&room); err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("database error: room with number %d not found", contract.RoomNumber)
		}
		return err
	}

	// insert contract
	contractCollection := s.database.Collection(contractCollection)
	_, err := contractCollection.InsertOne(ctx, contract)
	return err
}

func (s *MongoService) UpdateContract(ctx context.Context, contract *types.Contract) error {
	collection := s.database.Collection(contractCollection)

	res, err := collection.UpdateOne(ctx, bson.M{"contract_id": contract.ContractId}, bson.M{"$set": contract})
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (s *MongoService) DeleteContract(ctx context.Context, contractID string) error {
	collection := s.database.Collection(contractCollection)

	result, err := collection.DeleteOne(ctx, bson.M{"contract_id": contractID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (s *MongoService) GetContract(ctx context.Context, contractID string) (*types.Contract, error) {
	collection := s.database.Collection(contractCollection)

	var contract types.Contract
	if err := collection.FindOne(ctx, bson.M{"contract_id": contractID}).Decode(&contract); err != nil {
		return nil, err
	}

	return &contract, nil
}

func (s *MongoService) ListContracts(ctx context.Context) ([]*types.Contract, error) {
	collection := s.database.Collection(contractCollection)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var res []*types.Contract
	for cursor.Next(ctx) {
		var contract types.Contract
		if err := cursor.Decode(&contract); err != nil {
			return nil, err
		}
		res = append(res, &contract)
	}

	return res, cursor.Err()
}
