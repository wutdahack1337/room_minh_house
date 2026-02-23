package server

import (
	"context"
	"testing"
	"time"

	"github.com/wutdahack1337/room_minh_house/backend/internal/database"
	types "github.com/wutdahack1337/room_minh_house/backend/pkg/pb"
)

const (
	CONTRACT_ID         = "test-contract-id"
	INVOICE_PERIOD      = "02-2026"
	PREV_INVOICE_PERIOD = "01-2026"
	ELEC_UNIT_PRICE     = 3500.0
	WATER_UNIT_PRICE    = 15000.0
	ELEC_READING_PREV   = 100.0
	ELEC_READING_CURR   = 150.0
	WATER_READING_PREV  = 50.0
	WATER_READING_CURR  = 70.0
)

// seedInvoice creates a user, room, contract, and two meter readings (previous and current month)
func seedInvoice(t *testing.T, service *Service) string {
	t.Helper()
	ctx := context.Background()

	// create user
	_, err := service.CreateUser(ctx, &types.CreateUserRequest{
		User: &types.User{
			Pid:         P_ID,
			PhoneNumber: PHONE_NUMBER,
			Role:        ROLE,
			IdCard:      newTestIdCard(),
			Account:     newTestAccount(),
		},
	})
	if err != nil {
		t.Fatal("seed: failed to create user: " + err.Error())
	}

	// create room
	_, err = service.CreateRoom(ctx, &types.CreateRoomRequest{
		Room: &types.Room{
			RoomNumber:  ROOM_NUMBER,
			Price:       PRICE,
			IsAvailable: false,
		},
	})
	if err != nil {
		t.Fatal("seed: failed to create room: " + err.Error())
	}

	// create contract
	contractResp, err := service.CreateContract(ctx, &types.CreateContractRequest{
		Contract: &types.Contract{
			UserPid:        P_ID,
			RoomNumber:     ROOM_NUMBER,
			StartDate:      time.Now().UTC().Format(time.RFC3339),
			DurationMonths: 3,
		},
	})
	if err != nil {
		t.Fatal("seed: failed to create contract: " + err.Error())
	}
	contractID := contractResp.Contract.ContractId

	// create previous month meter reading
	_, err = service.CreateMeterReading(ctx, &types.CreateMeterReadingRequest{
		MeterReading: &types.MeterReading{
			RoomNumber:         ROOM_NUMBER,
			Period:             PREV_INVOICE_PERIOD,
			ElectricityReading: ELEC_READING_PREV,
			WaterReading:       WATER_READING_PREV,
		},
	})
	if err != nil {
		t.Fatal("seed: failed to create prev meter reading: " + err.Error())
	}

	// create current month meter reading
	_, err = service.CreateMeterReading(ctx, &types.CreateMeterReadingRequest{
		MeterReading: &types.MeterReading{
			RoomNumber:         ROOM_NUMBER,
			Period:             INVOICE_PERIOD,
			ElectricityReading: ELEC_READING_CURR,
			WaterReading:       WATER_READING_CURR,
		},
	})
	if err != nil {
		t.Fatal("seed: failed to create curr meter reading: " + err.Error())
	}

	return contractID
}

func TestCreateInvoice_HappyPath(t *testing.T) {
	ctx := context.Background()
	service := getServiceForInvoice()

	contractID := seedInvoice(t, service)

	response, err := service.CreateInvoice(ctx, &types.CreateInvoiceRequest{
		ContractId:           contractID,
		Period:               INVOICE_PERIOD,
		ElectricityUnitPrice: ELEC_UNIT_PRICE,
		WaterUnitPrice:       WATER_UNIT_PRICE,
	})
	if err != nil {
		t.Fatal("CreateInvoice failed: " + err.Error())
	}

	if response == nil {
		t.Fatal("CreateInvoice returned nil response")
	}
	if response.Invoice.Status != "unpaid" {
		t.Fatalf("expected status 'unpaid', got '%s'", response.Invoice.Status)
	}
	if len(response.InvoiceItems) != 3 {
		t.Fatalf("expected 3 invoice items, got %d", len(response.InvoiceItems))
	}

	expectedElecQty := ELEC_READING_CURR - ELEC_READING_PREV
	expectedWaterQty := WATER_READING_CURR - WATER_READING_PREV
	expectedTotal := PRICE + expectedElecQty*ELEC_UNIT_PRICE + expectedWaterQty*WATER_UNIT_PRICE
	if expectedTotal != response.Invoice.TotalAmount {
		t.Fatalf("expected total_amount %.2f, got %.2f", expectedTotal, response.Invoice.TotalAmount)
	}
}

func TestCreateInvoice_MissingContractId(t *testing.T) {
	service := getServiceForInvoice()

	response, err := service.CreateInvoice(context.Background(), &types.CreateInvoiceRequest{
		Period:               INVOICE_PERIOD,
		ElectricityUnitPrice: ELEC_UNIT_PRICE,
		WaterUnitPrice:       WATER_UNIT_PRICE,
	})
	if err == nil {
		t.Fatal("expected error for missing contract_id, got nil")
	}
	if err.Error() != "CreateInvoice error: contract_id is required" {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	if response != nil {
		t.Fatal("expected nil response")
	}
}

func TestCreateInvoice_MissingPeriod(t *testing.T) {
	service := getServiceForInvoice()

	response, err := service.CreateInvoice(context.Background(), &types.CreateInvoiceRequest{
		ContractId:           CONTRACT_ID,
		ElectricityUnitPrice: ELEC_UNIT_PRICE,
		WaterUnitPrice:       WATER_UNIT_PRICE,
	})
	if err == nil {
		t.Fatal("expected error for missing period, got nil")
	}
	if err.Error() != "CreateInvoice error: period is required" {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	if response != nil {
		t.Fatal("expected nil response")
	}
}

func TestCreateInvoice_InvalidPeriodFormat(t *testing.T) {
	service := getServiceForInvoice()

	response, err := service.CreateInvoice(context.Background(), &types.CreateInvoiceRequest{
		ContractId:           CONTRACT_ID,
		Period:               "2026-02",
		ElectricityUnitPrice: ELEC_UNIT_PRICE,
		WaterUnitPrice:       WATER_UNIT_PRICE,
	})
	if err == nil {
		t.Fatal("expected error for invalid period format, got nil")
	}

	if response != nil {
		t.Fatal("expected nil response")
	}
}

func TestCreateInvoice_ContractNotFound(t *testing.T) {
	service := getServiceForInvoice()

	response, err := service.CreateInvoice(context.Background(), &types.CreateInvoiceRequest{
		ContractId:           "non-existent-contract",
		Period:               INVOICE_PERIOD,
		ElectricityUnitPrice: ELEC_UNIT_PRICE,
		WaterUnitPrice:       WATER_UNIT_PRICE,
	})
	if err == nil {
		t.Fatal("expected error for non-existent contract, got nil")
	}

	if response != nil {
		t.Fatal("expected nil response")
	}
}

func TestCreateInvoice_MeterReadingNotFound(t *testing.T) {
	ctx := context.Background()
	service := getServiceForInvoice()

	contractID := seedInvoice(t, service)

	response, err := service.CreateInvoice(ctx, &types.CreateInvoiceRequest{
		ContractId:           contractID,
		Period:               "01-0001", // no meter reading exists for this period
		ElectricityUnitPrice: ELEC_UNIT_PRICE,
		WaterUnitPrice:       WATER_UNIT_PRICE,
	})
	if err == nil {
		t.Fatal("expected error for missing meter reading, got nil")
	}

	if response != nil {
		t.Fatal("expected nil response")
	}
}

func TestGetInvoice_HappyPath(t *testing.T) {
	ctx := context.Background()
	service := getServiceForInvoice()

	contractID := seedInvoice(t, service)

	_, err := service.CreateInvoice(ctx, &types.CreateInvoiceRequest{
		ContractId:           contractID,
		Period:               INVOICE_PERIOD,
		ElectricityUnitPrice: ELEC_UNIT_PRICE,
		WaterUnitPrice:       WATER_UNIT_PRICE,
	})
	if err != nil {
		t.Fatal("setup: CreateInvoice failed: " + err.Error())
	}

	response, err := service.GetInvoice(ctx, &types.GetInvoiceRequest{
		ContractId: contractID,
		Period:     INVOICE_PERIOD,
	})
	if err != nil {
		t.Fatal("GetInvoice failed: " + err.Error())
	}

	if response.Invoice.ContractId != contractID {
		t.Fatalf("expected contract_id %s, got %s", contractID, response.Invoice.ContractId)
	}
	if response.Invoice.Period != INVOICE_PERIOD {
		t.Fatalf("expected period %s, got %s", INVOICE_PERIOD, response.Invoice.Period)
	}
}

func TestListInvoices_HappyPath(t *testing.T) {
	ctx := context.Background()
	service := getServiceForInvoice()

	contractID := seedInvoice(t, service)

	_, err := service.CreateInvoice(ctx, &types.CreateInvoiceRequest{
		ContractId:           contractID,
		Period:               INVOICE_PERIOD,
		ElectricityUnitPrice: ELEC_UNIT_PRICE,
		WaterUnitPrice:       WATER_UNIT_PRICE,
	})
	if err != nil {
		t.Fatal("setup: CreateInvoice failed: " + err.Error())
	}

	response, err := service.ListInvoices(ctx, &types.ListInvoicesRequest{
		ContractId: contractID,
	})
	if err != nil {
		t.Fatal("ListInvoices failed: " + err.Error())
	}

	if len(response.Invoices) != 1 {
		t.Fatalf("expected 1 invoice, got %d", len(response.Invoices))
	}
}

func TestListInvoices_Empty(t *testing.T) {
	service := getServiceForInvoice()

	response, err := service.ListInvoices(context.Background(), &types.ListInvoicesRequest{
		ContractId: "no-invoices-here",
	})
	if err != nil {
		t.Fatal("ListInvoices failed: " + err.Error())
	}
	if len(response.Invoices) != 0 {
		t.Fatalf("expected 0 invoices, got %d", len(response.Invoices))
	}
}

func TestDeleteInvoice_HappyPath(t *testing.T) {
	ctx := context.Background()
	service := getServiceForInvoice()

	contractID := seedInvoice(t, service)

	_, err := service.CreateInvoice(ctx, &types.CreateInvoiceRequest{
		ContractId:           contractID,
		Period:               INVOICE_PERIOD,
		ElectricityUnitPrice: ELEC_UNIT_PRICE,
		WaterUnitPrice:       WATER_UNIT_PRICE,
	})
	if err != nil {
		t.Fatal("setup: CreateInvoice failed: " + err.Error())
	}

	_, err = service.DeleteInvoice(ctx, &types.DeleteInvoiceRequest{
		ContractId: contractID,
		Period:     INVOICE_PERIOD,
	})
	if err != nil {
		t.Fatal("DeleteInvoice failed: " + err.Error())
	}
}

// helper
func getServiceForInvoice() *Service {
	ctx := context.Background()

	db, err := connectDatabase("mongodb://localhost:27017/?directConnection=true", "room_minh_house")
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	db.Collection("meter_readings").Drop(ctx)
	db.Collection("users").Drop(ctx)
	db.Collection("rooms").Drop(ctx)
	db.Collection("contracts").Drop(ctx)
	db.Collection("invoices").Drop(ctx)
	db.Collection("invoice_items").Drop(ctx)

	dbService, err := database.NewMongoService(db)
	if err != nil {
		panic("Failed to create database service: " + err.Error())
	}

	return NewService(dbService)
}
