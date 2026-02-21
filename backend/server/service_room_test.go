package server

import (
	"context"
	"testing"

	"github.com/wutdahack1337/room_minh_house/backend/internal/database"
	types "github.com/wutdahack1337/room_minh_house/backend/pkg/pb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ROOM_NUMBER = 666
	PRICE       = 666_666
)

func TestCreateRoom_HappyPath(t *testing.T) {
	service := getService()

	response, err := service.CreateRoom(context.Background(), &types.CreateRoomRequest{
		Room: &types.Room{
			RoomNumber:  ROOM_NUMBER,
			Price:       PRICE,
			IsAvailable: true,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if response == nil {
		t.Fatal("CreateRoom returned nil response")
	}
}

func TestCreateRoom_DuplicateRoomNumber(t *testing.T) {
	ctx := context.Background()
	service := getService()

	room := &types.Room{
		RoomNumber:  ROOM_NUMBER,
		Price:       PRICE,
		IsAvailable: true,
	}

	_, _ = service.CreateRoom(ctx, &types.CreateRoomRequest{
		Room: room,
	})

	response, err := service.CreateRoom(ctx, &types.CreateRoomRequest{
		Room: room,
	})

	if err == nil {
		t.Fatal("expected error for duplicate room number, got nil")
	}

	if response != nil {
		t.Fatal("expected nil response on duplicate room number error")
	}
}

func TestCreateRoom_MissingRoomNumber(t *testing.T) {
	service := getService()
	response, err := service.CreateRoom(context.Background(), &types.CreateRoomRequest{
		Room: &types.Room{
			Price:       PRICE,
			IsAvailable: true,
		},
	})

	if err == nil {
		t.Fatal("expected error for missing room number, got nil")
	}

	if response != nil {
		t.Fatal("expected nil response for missing room number, got non-nil")
	}
}

func TestUpdateRoom_HappyPath(t *testing.T) {
	ctx := context.Background()
	service := getService()

	_, _ = service.CreateRoom(ctx, &types.CreateRoomRequest{
		Room: &types.Room{
			RoomNumber:  ROOM_NUMBER,
			Price:       PRICE,
			IsAvailable: true,
		},
	})

	room := &types.Room{
		RoomNumber:  ROOM_NUMBER,
		Price:       9_999_999,
		IsAvailable: false,
	}

	_, err := service.UpdateRoom(ctx, &types.UpdateRoomRequest{
		Room: room,
	})

	if err != nil {
		t.Fatal("UpdateRoom failed: " + err.Error())
	}

	response, err := service.GetRoom(ctx, &types.GetRoomRequest{
		RoomNumber: ROOM_NUMBER,
	})

	if response == nil {
		t.Fatal("GetRoom returned nil response")
	}
	if response.Room.Price != room.Price || response.Room.IsAvailable != room.IsAvailable {
		t.Fatalf("expected room price %f and isAvailable %t, got %f and %t", room.Price, room.IsAvailable, response.Room.Price, response.Room.IsAvailable)
	}
}

func TestUpdateRoom_MissingRoomNumber(t *testing.T) {
	service := getService()
	_, err := service.UpdateRoom(context.Background(), &types.UpdateRoomRequest{
		Room: &types.Room{
			Price:       6_666_666,
			IsAvailable: true,
		},
	})

	if err == nil {
		t.Fatal("expected error for missing room number, got nil")
	}
}

func TestUpdateRoom_NotFound(t *testing.T) {
	service := getService()

	response, err := service.UpdateRoom(context.Background(), &types.UpdateRoomRequest{
		Room: &types.Room{
			RoomNumber:  999_999,
			Price:       9_999_999,
			IsAvailable: true,
		},
	})

	if err == nil {
		t.Fatal("expected error for room not found, got nil")
	}

	if response != nil {
		t.Fatal("expected nil response when room not found")
	}
}

func TestDeleteRoom_HappyPath(t *testing.T) {
	ctx := context.Background()
	service := getService()

	createResponse, err := service.CreateRoom(ctx, &types.CreateRoomRequest{
		Room: &types.Room{
			RoomNumber:  ROOM_NUMBER,
			Price:       PRICE,
			IsAvailable: true,
		},
	})
	if err != nil {
		t.Fatal("CreateRoom failed: " + err.Error())
	}
	if createResponse == nil {
		t.Fatal("CreateRoom returned nil response")
	}

	deleteResponse, err := service.DeleteRoom(ctx, &types.DeleteRoomRequest{
		RoomNumber: ROOM_NUMBER,
	})
	if err != nil {
		t.Fatal("DeleteRoom failed: " + err.Error())
	}
	if deleteResponse == nil {
		t.Fatal("DeleteRoom returned nil response")
	}

	// Verify room is deleted
	getResponse, err := service.GetRoom(ctx, &types.GetRoomRequest{
		RoomNumber: ROOM_NUMBER,
	})
	if err == nil {
		t.Fatal("expected error when getting deleted room, got nil")
	}
	if getResponse != nil {
		t.Fatal("expected nil response when getting deleted room")
	}
}

func TestDeleteRoom_NotFound(t *testing.T) {
	service := getService()
	deleteResponse, err := service.DeleteRoom(context.Background(), &types.DeleteRoomRequest{
		RoomNumber: 999_999,
	})

	if err == nil {
		t.Fatal("expected error for room not found, got nil")
	}
	if deleteResponse != nil {
		t.Fatal("expected nil response when deleting non-existent room")
	}
}

func TestGetRoom_HappyPath(t *testing.T) {
	ctx := context.Background()
	service := getService()

	room := &types.Room{
		RoomNumber:  ROOM_NUMBER,
		Price:       PRICE,
		IsAvailable: true,
	}
	_, _ = service.CreateRoom(ctx, &types.CreateRoomRequest{
		Room: room,
	})

	response, err := service.GetRoom(ctx, &types.GetRoomRequest{
		RoomNumber: ROOM_NUMBER,
	})
	if err != nil {
		t.Fatal("GetRoom failed: " + err.Error())
	}

	if response == nil {
		t.Fatal("GetRoom returned nil response")
	}
	if response.Room == nil {
		t.Fatal("GetRoom response contains nil room")
	}
	if response.Room.RoomNumber != room.RoomNumber {
		t.Fatalf("expected room number %d, got %d", room.RoomNumber, response.Room.RoomNumber)
	}
	if response.Room.Price != room.Price {
		t.Fatalf("expected price %f, got %f", room.Price, response.Room.Price)
	}
}

func TestGetRoom_NotFound(t *testing.T) {
	service := getService()
	response, err := service.GetRoom(context.Background(), &types.GetRoomRequest{
		RoomNumber: 999_999,
	})

	if err == nil {
		t.Fatal("expected error for room not found, got nil")
	}
	if response != nil {
		t.Fatal("expected nil response when room not found")
	}
}

func TestListRooms_HappyPath(t *testing.T) {
	ctx := context.Background()
	service := getService()

	rooms := []*types.Room{
		{RoomNumber: 666, Price: 666_666, IsAvailable: true},
		{RoomNumber: 963, Price: 963_963, IsAvailable: false},
		{RoomNumber: 999, Price: 999_999, IsAvailable: true},
	}
	for _, room := range rooms {
		_, err := service.CreateRoom(ctx, &types.CreateRoomRequest{
			Room: room,
		})
		if err != nil {
			t.Fatalf("CreateRoom failed for room %d: %v", room.RoomNumber, err)
		}
	}

	response, err := service.ListRooms(ctx, &types.ListRoomsRequest{})
	if err != nil {
		t.Fatal("ListRooms failed: " + err.Error())
	}

	if response == nil {
		t.Fatal("ListRooms returned nil response")
	}
	if len(response.Rooms) != 3 {
		t.Fatalf("expected 3 rooms, got %d", len(response.Rooms))
	}

	for i, room := range rooms {
		if response.Rooms[i].RoomNumber != room.RoomNumber || response.Rooms[i].Price != room.Price || response.Rooms[i].IsAvailable != room.IsAvailable {
			t.Fatalf("expected room number %d, price %f, isAvailable %t, got room number %d, price %f, isAvailable %t", room.RoomNumber, room.Price, room.IsAvailable, response.Rooms[i].RoomNumber, response.Rooms[i].Price, response.Rooms[i].IsAvailable)
		}
	}
}

func TestListRooms_Empty(t *testing.T) {
	service := getService()
	response, err := service.ListRooms(context.Background(), &types.ListRoomsRequest{})

	if err != nil {
		t.Fatal("ListRooms failed: " + err.Error())
	}
	if response == nil {
		t.Fatal("ListRooms returned nil response")
	}
	if len(response.Rooms) != 0 {
		t.Fatalf("expected 0 rooms, got %d", len(response.Rooms))
	}
}

// helper
func connectDatabase(uri, dbName string) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	return client.Database(dbName), nil
}

func getService() *Service {
	db, err := connectDatabase("mongodb://localhost:27017/?directConnection=true", "room_minh_house")
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// removes all documents and indexes
	db.Collection("rooms").Drop(context.Background())

	dbService, err := database.NewMongoService(db)
	if err != nil {
		panic("Failed to create database service: " + err.Error())
	}

	return NewService(dbService)
}
