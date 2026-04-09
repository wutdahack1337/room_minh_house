package server

import (
	"context"
	"strings"
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

func TestCreateMeterReading_InvalidPeriodFormat(t *testing.T) {
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

	invalidPeriods := []string{
		"2025-08",   // wrong order
		"8-2025",    // missing leading zero
		"13-2025",   // invalid month
		"00-2025",   // month 00
		"08-25",     // year too short
		"08-202500", // year too long
		"abc-defg",  // non-numeric
		"08/2025",   // wrong separator
	}

	for _, invalidPeriod := range invalidPeriods {
		response, err := service.CreateMeterReading(ctx, &types.CreateMeterReadingRequest{
			MeterReading: &types.MeterReading{
				RoomNumber:         ROOM_NUMBER,
				Period:             invalidPeriod,
				ElectricityReading: ELECTRICITY_READING,
				WaterReading:       WATER_READING,
			},
		})
		if err == nil {
			t.Fatalf("expected error for invalid period format '%s', got nil", invalidPeriod)
		}
		if !containsSubstring(err.Error(), "period must be in MM-YYYY format") {
			t.Fatalf("expected error to contain 'period must be in MM-YYYY format', got '%s'", err.Error())
		}

		if response != nil {
			t.Fatal("expected nil response for invalid period format")
		}
	}
}

func TestCreateMeterReading_NegativeReadings(t *testing.T) {
	ctx := context.Background()
	service := getServiceForMeterReading()

	_, err := service.CreateRoom(ctx, &types.CreateRoomRequest{
		Room: &types.Room{RoomNumber: ROOM_NUMBER, Price: PRICE, IsAvailable: false},
	})
	if err != nil {
		t.Fatal("Failed to create room: " + err.Error())
	}

	response, err := service.CreateMeterReading(ctx, &types.CreateMeterReadingRequest{
		MeterReading: &types.MeterReading{
			RoomNumber:         ROOM_NUMBER,
			Period:             PERIOD,
			ElectricityReading: -1,
			WaterReading:       10,
		},
	})
	if err == nil {
		t.Fatal("expected error for negative electricity reading, got nil")
	}
	if !containsSubstring(err.Error(), "electricity reading must be greater than or equal to 0") {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if response != nil {
		t.Fatal("expected nil response for negative readings")
	}
}

func TestCreateMeterReading_PeriodMustBeGreaterThanLatest(t *testing.T) {
	ctx := context.Background()
	service := getServiceForMeterReading()

	_, err := service.CreateRoom(ctx, &types.CreateRoomRequest{
		Room: &types.Room{RoomNumber: ROOM_NUMBER, Price: PRICE, IsAvailable: false},
	})
	if err != nil {
		t.Fatal("Failed to create room: " + err.Error())
	}

	_, err = service.CreateMeterReading(ctx, &types.CreateMeterReadingRequest{
		MeterReading: &types.MeterReading{
			RoomNumber:         ROOM_NUMBER,
			Period:             "02-2026",
			ElectricityReading: 100,
			WaterReading:       50,
		},
	})
	if err != nil {
		t.Fatal("Failed to create initial meter reading: " + err.Error())
	}

	response, err := service.CreateMeterReading(ctx, &types.CreateMeterReadingRequest{
		MeterReading: &types.MeterReading{
			RoomNumber:         ROOM_NUMBER,
			Period:             "01-2026",
			ElectricityReading: 110,
			WaterReading:       60,
		},
	})
	if err == nil {
		t.Fatal("expected error for non-increasing period, got nil")
	}
	if !containsSubstring(err.Error(), "period must be greater than latest existing period") {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if response != nil {
		t.Fatal("expected nil response for non-increasing period")
	}
}

func TestCreateMeterReading_ReadingCanBeEqualToLatestPeriod(t *testing.T) {
	ctx := context.Background()
	service := getServiceForMeterReading()

	_, err := service.CreateRoom(ctx, &types.CreateRoomRequest{
		Room: &types.Room{RoomNumber: ROOM_NUMBER, Price: PRICE, IsAvailable: false},
	})
	if err != nil {
		t.Fatal("Failed to create room: " + err.Error())
	}

	_, err = service.CreateMeterReading(ctx, &types.CreateMeterReadingRequest{
		MeterReading: &types.MeterReading{
			RoomNumber:         ROOM_NUMBER,
			Period:             "06-2026",
			ElectricityReading: 250,
			WaterReading:       50,
		},
	})
	if err != nil {
		t.Fatal("Failed to create initial meter reading: " + err.Error())
	}

	response, err := service.CreateMeterReading(ctx, &types.CreateMeterReadingRequest{
		MeterReading: &types.MeterReading{
			RoomNumber:         ROOM_NUMBER,
			Period:             "07-2026",
			ElectricityReading: 250,
			WaterReading:       50,
		},
	})
	if err != nil {
		t.Fatalf("expected success when readings are equal to previous period, got error: %s", err.Error())
	}
	if response == nil {
		t.Fatal("expected non-nil response when readings are equal to previous period")
	}
}

func TestCreateMeterReading_ReadingCannotBeLowerThanLatestPeriod(t *testing.T) {
	ctx := context.Background()
	service := getServiceForMeterReading()

	_, err := service.CreateRoom(ctx, &types.CreateRoomRequest{
		Room: &types.Room{RoomNumber: ROOM_NUMBER, Price: PRICE, IsAvailable: false},
	})
	if err != nil {
		t.Fatal("Failed to create room: " + err.Error())
	}

	_, err = service.CreateMeterReading(ctx, &types.CreateMeterReadingRequest{
		MeterReading: &types.MeterReading{
			RoomNumber:         ROOM_NUMBER,
			Period:             "06-2026",
			ElectricityReading: 250,
			WaterReading:       50,
		},
	})
	if err != nil {
		t.Fatal("Failed to create initial meter reading: " + err.Error())
	}

	response, err := service.CreateMeterReading(ctx, &types.CreateMeterReadingRequest{
		MeterReading: &types.MeterReading{
			RoomNumber:         ROOM_NUMBER,
			Period:             "07-2026",
			ElectricityReading: 10,
			WaterReading:       10,
		},
	})
	if err == nil {
		t.Fatal("expected error for lower reading values, got nil")
	}
	if !containsSubstring(err.Error(), "electricity reading of new period must be greater than or equal to previous period") {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if response != nil {
		t.Fatal("expected nil response for lower reading values")
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

func TestUpdateMeterReading_InvalidPeriodFormat(t *testing.T) {
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

	invalidPeriods := []string{
		"2025-08",
		"8-2025",
		"13-2025",
		"00-2025",
		"08-25",
		"08-202500",
		"abc-defg",
	}

	for _, invalidPeriod := range invalidPeriods {
		response, err := service.UpdateMeterReading(ctx, &types.UpdateMeterReadingRequest{
			MeterReading: &types.MeterReading{
				RoomNumber:         ROOM_NUMBER,
				Period:             invalidPeriod,
				ElectricityReading: ELECTRICITY_READING,
				WaterReading:       WATER_READING,
			},
		})
		if err == nil {
			t.Fatalf("expected error for invalid period format '%s', got nil", invalidPeriod)
		}
		if !containsSubstring(err.Error(), "period must be in MM-YYYY format") {
			t.Fatalf("expected error to contain 'period must be in MM-YYYY format', got '%s'", err.Error())
		}

		if response != nil {
			t.Fatal("expected nil response for invalid period format")
		}
	}
}

func TestUpdateMeterReading_NegativeReadings(t *testing.T) {
	ctx := context.Background()
	service := getServiceForMeterReading()

	_, err := service.CreateRoom(ctx, &types.CreateRoomRequest{
		Room: &types.Room{RoomNumber: ROOM_NUMBER, Price: PRICE, IsAvailable: false},
	})
	if err != nil {
		t.Fatal("Failed to create room: " + err.Error())
	}

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

	response, err := service.UpdateMeterReading(ctx, &types.UpdateMeterReadingRequest{
		MeterReading: &types.MeterReading{
			RoomNumber:         ROOM_NUMBER,
			Period:             PERIOD,
			ElectricityReading: -10,
			WaterReading:       20,
		},
	})
	if err == nil {
		t.Fatal("expected error for negative electricity reading, got nil")
	}
	if !containsSubstring(err.Error(), "electricity reading must be greater than or equal to 0") {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if response != nil {
		t.Fatal("expected nil response for negative readings")
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

func TestGetMeterReading_InvalidPeriodFormat(t *testing.T) {
	service := getServiceForMeterReading()

	invalidPeriods := []string{
		"2025-08",
		"8-2025",
		"13-2025",
		"00-2025",
		"08-25",
		"08-202500",
		"abc-defg",
		"08/2025",
	}

	for _, invalidPeriod := range invalidPeriods {
		response, err := service.GetMeterReading(context.Background(), &types.GetMeterReadingRequest{
			RoomNumber: ROOM_NUMBER,
			Period:     invalidPeriod,
		})
		if err == nil {
			t.Fatalf("expected error for invalid period format '%s', got nil", invalidPeriod)
		}
		if !containsSubstring(err.Error(), "period must be in MM-YYYY format") {
			t.Fatalf("expected error to contain 'period must be in MM-YYYY format', got '%s'", err.Error())
		}

		if response != nil {
			t.Fatal("expected nil response for invalid period format")
		}
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

// helpers
func getServiceForMeterReading() *Service {
	ctx := context.Background()

	db, err := connectDatabase("mongodb://localhost:27017/?directConnection=true", "room_minh_house")
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// removes all documents and indexes
	db.Collection("rooms").Drop(ctx)
	db.Collection("meter_readings").Drop(ctx)

	dbService, err := database.NewMongoService(db)
	if err != nil {
		panic("Failed to create database service: " + err.Error())
	}

	return NewService(dbService)
}

func containsSubstring(str, substr string) bool {
	return strings.Contains(str, substr)
}

// ─── Cascade Tests ───────────────────────────────────────────────────────────

func TestUpdateMeterReading_CascadesInvoiceItems(t *testing.T) {
	ctx := context.Background()
	service := getServiceForInvoice() // drops all collections

	// seed user, room, contract, prev & current meter readings
	contractID := seedInvoice(t, service)

	// create invoice for INVOICE_PERIOD
	_, err := service.CreateInvoice(ctx, &types.CreateInvoiceRequest{
		ContractId: contractID, Period: INVOICE_PERIOD,
		ElectricityUnitPrice: ELEC_UNIT_PRICE, WaterUnitPrice: WATER_UNIT_PRICE,
	})
	if err != nil {
		t.Fatal("setup: CreateInvoice failed: " + err.Error())
	}

	// original electricity qty = ELEC_READING_CURR - ELEC_READING_PREV = 50
	// update current meter reading with higher values
	newElecReading := ELEC_READING_CURR + 20  // 170
	newWaterReading := WATER_READING_CURR + 5 // 75

	_, err = service.UpdateMeterReading(ctx, &types.UpdateMeterReadingRequest{
		MeterReading: &types.MeterReading{
			RoomNumber: ROOM_NUMBER, Period: INVOICE_PERIOD,
			ElectricityReading: newElecReading,
			WaterReading:       newWaterReading,
		},
	})
	if err != nil {
		t.Fatal("UpdateMeterReading failed: " + err.Error())
	}

	// verify invoice items were updated
	itemsResp, err := service.ListInvoiceItems(ctx, &types.ListInvoiceItemsRequest{
		ContractId: contractID, Period: INVOICE_PERIOD,
	})
	if err != nil {
		t.Fatal("ListInvoiceItems failed: " + err.Error())
	}

	newElecQty := newElecReading - ELEC_READING_PREV    // 70
	newWaterQty := newWaterReading - WATER_READING_PREV // 25

	for _, item := range itemsResp.InvoiceItems {
		switch item.ItemType {
		case ITEM_TYPE_ELECTRICITY:
			if item.Quantity != newElecQty {
				t.Fatalf("electricity qty: expected %.2f, got %.2f", newElecQty, item.Quantity)
			}
			expectedAmount := newElecQty * ELEC_UNIT_PRICE
			if item.Amount != expectedAmount {
				t.Fatalf("electricity amount: expected %.2f, got %.2f", expectedAmount, item.Amount)
			}
		case ITEM_TYPE_WATER:
			if item.Quantity != newWaterQty {
				t.Fatalf("water qty: expected %.2f, got %.2f", newWaterQty, item.Quantity)
			}
			expectedAmount := newWaterQty * WATER_UNIT_PRICE
			if item.Amount != expectedAmount {
				t.Fatalf("water amount: expected %.2f, got %.2f", expectedAmount, item.Amount)
			}
		}
	}

	// verify invoice total_amount was recalculated
	invoiceResp, err := service.GetInvoice(ctx, &types.GetInvoiceRequest{
		ContractId: contractID, Period: INVOICE_PERIOD,
	})
	if err != nil {
		t.Fatal("GetInvoice failed: " + err.Error())
	}

	expectedTotal := PRICE + newElecQty*ELEC_UNIT_PRICE + newWaterQty*WATER_UNIT_PRICE
	if invoiceResp.Invoice.TotalAmount != expectedTotal {
		t.Fatalf("invoice total: expected %.2f, got %.2f", expectedTotal, invoiceResp.Invoice.TotalAmount)
	}
}
