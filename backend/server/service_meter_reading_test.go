package server

import (
	"context"
	"testing"

	"github.com/wutdahack1337/room_minh_house/backend/internal/database"
	types "github.com/wutdahack1337/room_minh_house/backend/pkg/pb"
)

const (
	PERIOD              = "08-2005"
	ELECTRICITY_READING = 69.0
	WATER_READING       = 67.0
)

func TestCreateMeterReading_HappyPath(t *testing.T) {
	ctx := context.Background()
	service := getServiceForMeterReading()

	// Create a room first
	room := &types.Room{
		RoomNumber:  ROOM_NUMBER,
		Price:       PRICE,
		IsAvailable: false,
	}
	_, err := service.CreateRoom(ctx, &types.CreateRoomRequest{
		Room: room,
	})
	if err != nil {
		t.Fatal("Failed to create room: " + err.Error())
	}

	response, err := service.CreateMeterReading(ctx, &types.CreateMeterReadingRequest{
		MeterReading: &types.MeterReading{
			RoomNumber:         ROOM_NUMBER,
			Period:             PERIOD,
			ElectricityReading: ELECTRICITY_READING,
			WaterReading:       WATER_READING,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if response == nil {
		t.Fatal("CreateMeterReading returned nil response")
	}
	if response.MeterReading.Period != PERIOD {
		t.Fatalf("expected period %s, got %s", PERIOD, response.MeterReading.Period)
	}
}

func TestCreateMeterReading_MissingRoomNumber(t *testing.T) {
	service := getServiceForMeterReading()

	meterReading := &types.MeterReading{
		Period:             PERIOD,
		ElectricityReading: ELECTRICITY_READING,
		WaterReading:       WATER_READING,
	}

	response, err := service.CreateMeterReading(context.Background(), &types.CreateMeterReadingRequest{
		MeterReading: meterReading,
	})
	if err == nil {
		t.Fatal("expected error for missing room number, got nil")
	}
	if err.Error() != "CreateMeterReading error: room number is required" {
		t.Fatalf("expected 'CreateMeterReading error: room number is required', got '%s'", err.Error())
	}

	if response != nil {
		t.Fatal("expected nil response for missing room number, got non-nil")
	}
}

func TestCreateMeterReading_MissingPeriod(t *testing.T) {
	ctx := context.Background()
	service := getServiceForMeterReading()

	// Create a room first
	room := &types.Room{
		RoomNumber:  ROOM_NUMBER,
		Price:       PRICE,
		IsAvailable: false,
	}
	_, err := service.CreateRoom(ctx, &types.CreateRoomRequest{Room: room})
	if err != nil {
		t.Fatal("Failed to create room: " + err.Error())
	}

	meterReading := &types.MeterReading{
		RoomNumber:         ROOM_NUMBER,
		ElectricityReading: ELECTRICITY_READING,
		WaterReading:       WATER_READING,
	}

	response, err := service.CreateMeterReading(ctx, &types.CreateMeterReadingRequest{
		MeterReading: meterReading,
	})
	if err == nil {
		t.Fatal("expected error for missing period, got nil")
	}
	if err.Error() != "CreateMeterReading error: period is required" {
		t.Fatalf("expected 'CreateMeterReading error: period is required', got '%s'", err.Error())
	}

	if response != nil {
		t.Fatal("expected nil response for missing period, got non-nil")
	}
}

func TestCreateMeterReading_RoomNotFound(t *testing.T) {
	service := getServiceForMeterReading()

	meterReading := &types.MeterReading{
		RoomNumber:         999_999_999,
		Period:             PERIOD,
		ElectricityReading: ELECTRICITY_READING,
		WaterReading:       WATER_READING,
	}

	response, err := service.CreateMeterReading(context.Background(), &types.CreateMeterReadingRequest{
		MeterReading: meterReading,
	})

	if err == nil {
		t.Fatal("expected error for non-existent room, got nil")
	}
	if err.Error() != "database error: room with number 999999999 not found" {
		t.Fatalf("expected 'database error: room with number 999999999 not found', got '%s'", err.Error())
	}

	if response != nil {
		t.Fatal("expected nil response for non-existent room")
	}
}

func TestUpdateMeterReading_HappyPath(t *testing.T) {
	ctx := context.Background()
	service := getServiceForMeterReading()

	// Create a room first
	_, err := service.CreateRoom(ctx, &types.CreateRoomRequest{
		Room: &types.Room{
			RoomNumber:  ROOM_NUMBER,
			Price:       PRICE,
			IsAvailable: false,
		},
	})
	if err != nil {
		t.Fatal("Failed to create room: " + err.Error())
	}

	// Create meter reading
	meterReading := &types.MeterReading{
		RoomNumber:         ROOM_NUMBER,
		Period:             PERIOD,
		ElectricityReading: ELECTRICITY_READING,
		WaterReading:       WATER_READING,
	}

	_, err = service.CreateMeterReading(ctx, &types.CreateMeterReadingRequest{
		MeterReading: meterReading,
	})
	if err != nil {
		t.Fatal("Failed to create meter reading: " + err.Error())
	}

	// Update meter reading
	meterReading.ElectricityReading = 200.5
	meterReading.WaterReading = 30.5

	_, err = service.UpdateMeterReading(ctx, &types.UpdateMeterReadingRequest{
		MeterReading: meterReading,
	})
	if err != nil {
		t.Fatal("UpdateMeterReading failed: " + err.Error())
	}

	response, err := service.GetMeterReading(ctx, &types.GetMeterReadingRequest{
		RoomNumber: ROOM_NUMBER,
		Period:     PERIOD,
	})
	if err != nil {
		t.Fatal("GetMeterReading failed: " + err.Error())
	}

	if response == nil {
		t.Fatal("GetMeterReading returned nil response")
	}
	if response.MeterReading.Period != PERIOD {
		t.Fatalf("expected period %s, got %s", PERIOD, response.MeterReading.Period)
	}
	if response.MeterReading.ElectricityReading != 200.5 {
		t.Fatalf("expected electricity reading 200.5, got %f", response.MeterReading.ElectricityReading)
	}
	if response.MeterReading.WaterReading != 30.5 {
		t.Fatalf("expected water reading 30.5, got %f", response.MeterReading.WaterReading)
	}
}

func TestUpdateMeterReading_NotFound(t *testing.T) {
	service := getServiceForMeterReading()

	response, err := service.UpdateMeterReading(context.Background(), &types.UpdateMeterReadingRequest{
		MeterReading: &types.MeterReading{
			RoomNumber:         ROOM_NUMBER,
			Period:             "08-2105",
			ElectricityReading: ELECTRICITY_READING,
			WaterReading:       WATER_READING,
		},
	})
	if err == nil {
		t.Fatal("expected error for non-existent meter reading, got nil")
	}
	if err.Error() != "mongo: no documents in result" {
		t.Fatalf("expected 'mongo: no documents in result', got '%s'", err.Error())
	}

	if response != nil {
		t.Fatal("expected nil response for non-existent meter reading")
	}
}

func TestGetMeterReading_HappyPath(t *testing.T) {
	ctx := context.Background()
	service := getServiceForMeterReading()

	// Create a room first
	room := &types.Room{
		RoomNumber:  ROOM_NUMBER,
		Price:       PRICE,
		IsAvailable: false,
	}
	_, err := service.CreateRoom(ctx, &types.CreateRoomRequest{
		Room: room,
	})
	if err != nil {
		t.Fatal("Failed to create room: " + err.Error())
	}

	// Create meter reading
	_, err = service.CreateMeterReading(ctx, &types.CreateMeterReadingRequest{
		MeterReading: &types.MeterReading{
			RoomNumber:         ROOM_NUMBER,
			Period:             PERIOD,
			ElectricityReading: ELECTRICITY_READING,
			WaterReading:       WATER_READING,
		},
	})
	if err != nil {
		t.Fatal("Failed to create meter reading: " + err.Error())
	}

	response, err := service.GetMeterReading(ctx, &types.GetMeterReadingRequest{
		RoomNumber: ROOM_NUMBER,
		Period:     PERIOD,
	})
	if err != nil {
		t.Fatal("GetMeterReading failed: " + err.Error())
	}

	if response == nil {
		t.Fatal("GetMeterReading returned nil response")
	}
	if response.MeterReading.RoomNumber != ROOM_NUMBER {
		t.Fatalf("expected room number %d, got %d", ROOM_NUMBER, response.MeterReading.RoomNumber)
	}
	if response.MeterReading.Period != PERIOD {
		t.Fatalf("expected period %s, got %s", PERIOD, response.MeterReading.Period)
	}
}

func TestGetMeterReading_NotFound(t *testing.T) {
	service := getServiceForMeterReading()

	response, err := service.GetMeterReading(context.Background(), &types.GetMeterReadingRequest{
		RoomNumber: ROOM_NUMBER,
		Period:     "08-2105",
	})
	if err == nil {
		t.Fatal("expected error for non-existent meter reading, got nil")
	}
	if err.Error() != "mongo: no documents in result" {
		t.Fatalf("expected 'mongo: no documents in result', got '%s'", err.Error())
	}

	if response != nil {
		t.Fatal("expected nil response for non-existent meter reading")
	}
}

func TestListMeterReadings_HappyPath(t *testing.T) {
	ctx := context.Background()
	service := getServiceForMeterReading()

	// Create rooms
	rooms := []*types.Room{
		{RoomNumber: 666_666, Price: 666_666, IsAvailable: false},
		{RoomNumber: 999_999, Price: 999_999, IsAvailable: false},
	}
	for _, room := range rooms {
		_, err := service.CreateRoom(ctx, &types.CreateRoomRequest{Room: room})
		if err != nil {
			t.Fatal("Failed to create room: " + err.Error())
		}
	}

	// Create meter readings
	meterReadings := []*types.MeterReading{
		{RoomNumber: 666_666, Period: "08-2005", ElectricityReading: 150.5, WaterReading: 25.3},
		{RoomNumber: 666_666, Period: "02-2012", ElectricityReading: 170.5, WaterReading: 28.3},
		{RoomNumber: 999_999, Period: "09-1979", ElectricityReading: 120.0, WaterReading: 20.0},
	}

	for _, meterReading := range meterReadings {
		_, err := service.CreateMeterReading(ctx, &types.CreateMeterReadingRequest{
			MeterReading: meterReading,
		})
		if err != nil {
			t.Fatal("Failed to create meter reading: " + err.Error())
		}
	}

	response, err := service.ListMeterReadings(ctx, &types.ListMeterReadingsRequest{
		RoomNumber: 666_666,
	})
	if err != nil {
		t.Fatal("ListMeterReadings failed: " + err.Error())
	}

	if response == nil {
		t.Fatal("ListMeterReadings returned nil response")
	}
	if len(response.MeterReadings) != 2 {
		t.Fatalf("expected 2 meter readings for room 666_666, got %d", len(response.MeterReadings))
	}

	for _, meterReading := range response.MeterReadings {
		if meterReading.RoomNumber != 666_666 {
			t.Fatalf("expected room number 666_666, got %d", meterReading.RoomNumber)
		}
	}
}

func TestListMeterReadings_Empty(t *testing.T) {
	service := getServiceForMeterReading()

	response, err := service.ListMeterReadings(context.Background(), &types.ListMeterReadingsRequest{
		RoomNumber: ROOM_NUMBER,
	})
	if err != nil {
		t.Fatal("ListMeterReadings failed: " + err.Error())
	}

	if response == nil {
		t.Fatal("ListMeterReadings returned nil response")
	}
	if len(response.MeterReadings) != 0 {
		t.Fatalf("expected 0 meter readings, got %d", len(response.MeterReadings))
	}
}

// helper
func getServiceForMeterReading() *Service {
	db, err := connectDatabase("mongodb://localhost:27017/?directConnection=true", "room_minh_house")
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// removes all documents and indexes
	db.Collection("rooms").Drop(context.Background())
	db.Collection("meter_readings").Drop(context.Background())

	dbService, err := database.NewMongoService(db)
	if err != nil {
		panic("Failed to create database service: " + err.Error())
	}

	return NewService(dbService)
}
