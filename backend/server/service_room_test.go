package server

import (
	"context"
	"testing"

	"github.com/wutdahack1337/room_minh_house/backend/internal/database"
	types "github.com/wutdahack1337/room_minh_house/backend/pkg/pb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCreateRoom_HappyPath(t *testing.T) {
	room := &types.Room{
		RoomNumber:  666,
		Price:       9_999_999,
		IsAvailable: true,
	}

	service := getService()
	response, err := service.CreateRoom(context.Background(), &types.CreateRoomRequest{
		Room: room,
	})

	if err != nil {
		t.Fatal(err)
	}

	t.Log(response)
	if response == nil {
		t.Fatal("CreateRoom returned nil response")
	}
}

func TestCreateRoom_DuplicateRoomNumber(t *testing.T) {
	ctx := context.Background()

	room := &types.Room{
		RoomNumber:  666,
		Price:       9_999_999,
		IsAvailable: true,
	}

	service := getService()
	_, _ = service.CreateRoom(ctx, &types.CreateRoomRequest{
		Room: room,
	})

	response, err := service.CreateRoom(ctx, &types.CreateRoomRequest{
		Room: room,
	})

	t.Log(err)
	if err == nil {
		t.Fatal("expected error for duplicate room number, got nil")
	}

	if response != nil {
		t.Fatal("expected nil response on duplicate room number error")
	}
}

func TestUpdateRoom_HappyPath(t *testing.T) {
	ctx := context.Background()

	room := &types.Room{
		RoomNumber:  666,
		Price:       6_666_666,
		IsAvailable: true,
	}

	service := getService()
	_, _ = service.CreateRoom(ctx, &types.CreateRoomRequest{
		Room: room,
	})

	room = &types.Room{
		RoomNumber:  666,
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
		RoomNumber: 666,
	})

	t.Log(response)
	if response == nil {
		t.Fatal("GetRoom returned nil response")
	}

	if response.Room.Price != room.Price || response.Room.IsAvailable != room.IsAvailable {
		t.Fatalf("expected room price %f and isAvailable %t, got %f and %t", room.Price, room.IsAvailable, response.Room.Price, response.Room.IsAvailable)
	}
}

func TestUpdateRoom_MissingRoomNumber(t *testing.T) {
	room := &types.Room{
		Price:       6_666_666,
		IsAvailable: true,
	}

	service := getService()
	_, err := service.UpdateRoom(context.Background(), &types.UpdateRoomRequest{
		Room: room,
	})

	t.Log(err)
	if err == nil {
		t.Fatal("expected error for missing room number, got nil")
	}
}

func TestUpdateRoom_NotFound(t *testing.T) {
	room := &types.Room{
		RoomNumber:  999,
		Price:       9_999_999,
		IsAvailable: true,
	}

	service := getService()
	response, err := service.UpdateRoom(context.Background(), &types.UpdateRoomRequest{
		Room: room,
	})

	t.Log(err)
	if err == nil {
		t.Fatal("expected error for room not found, got nil")
	}
	if response != nil {
		t.Fatal("expected nil response when room not found")
	}
}

func TestDeleteRoom_HappyPath(t *testing.T) {
	ctx := context.Background()

	room := &types.Room{
		RoomNumber:  666,
		Price:       6_666_666,
		IsAvailable: true,
	}

	service := getService()
	createResponse, err := service.CreateRoom(ctx, &types.CreateRoomRequest{
		Room: room,
	})

	if err != nil {
		t.Fatal("CreateRoom failed: " + err.Error())
	}
	if createResponse == nil {
		t.Fatal("CreateRoom returned nil response")
	}

	deleteResponse, err := service.DeleteRoom(ctx, &types.DeleteRoomRequest{
		RoomNumber: 666,
	})

	if err != nil {
		t.Fatal("DeleteRoom failed: " + err.Error())
	}
	if deleteResponse == nil {
		t.Fatal("DeleteRoom returned nil response")
	}

	// Verify room is deleted
	getResponse, err := service.GetRoom(ctx, &types.GetRoomRequest{
		RoomNumber: 666,
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
		RoomNumber: 999,
	})

	t.Log(err)
	if err == nil {
		t.Fatal("expected error for room not found, got nil")
	}
	if deleteResponse != nil {
		t.Fatal("expected nil response when deleting non-existent room")
	}
}

func TestGetRoom_HappyPath(t *testing.T) {
	ctx := context.Background()

	room := &types.Room{
		RoomNumber:  666,
		Price:       666_666,
		IsAvailable: true,
	}

	service := getService()
	_, _ = service.CreateRoom(ctx, &types.CreateRoomRequest{
		Room: room,
	})

	response, err := service.GetRoom(ctx, &types.GetRoomRequest{
		RoomNumber: 666,
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
	if response.Room.RoomNumber != 666 {
		t.Fatalf("expected room number 666, got %d", response.Room.RoomNumber)
	}
	if response.Room.Price != 666_666 {
		t.Fatalf("expected price 666_666, got %f", response.Room.Price)
	}
}

func TestGetRoom_NotFound(t *testing.T) {
	service := getService()
	response, err := service.GetRoom(context.Background(), &types.GetRoomRequest{
		RoomNumber: 999,
	})

	t.Log(err)
	if err == nil {
		t.Fatal("expected error for room not found, got nil")
	}
	if response != nil {
		t.Fatal("expected nil response when room not found")
	}
}

func TestListRooms_HappyPath(t *testing.T) {
	ctx := context.Background()

	// Create multiple rooms
	rooms := []*types.Room{
		{RoomNumber: 666, Price: 666_666, IsAvailable: true},
		{RoomNumber: 963, Price: 963_963, IsAvailable: false},
		{RoomNumber: 999, Price: 999_999, IsAvailable: true},
	}

	service := getService()
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

	t.Log("ListRooms returned empty list as expected")
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
	ctx := context.Background()

	db, err := connectDatabase("mongodb://localhost:27017/?directConnection=true", "room_minh_house")
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// removes all documents and indexes
	db.Collection("rooms").Drop(ctx)

	dbService, err := database.NewMongoService(db)
	if err != nil {
		panic("Failed to create database service: " + err.Error())
	}

	service := NewService(dbService)
	return service
}
