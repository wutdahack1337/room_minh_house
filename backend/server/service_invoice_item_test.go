package server

import (
	"context"
	"testing"

	types "github.com/wutdahack1337/room_minh_house/backend/pkg/pb"
)

func TestUpdateInvoiceItem_HappyPath(t *testing.T) {
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

	newUnitPrice := 4000.0
	elecQty := ELEC_READING_CURR - ELEC_READING_PREV

	response, err := service.UpdateInvoiceItem(ctx, &types.UpdateInvoiceItemRequest{
		InvoiceItem: &types.InvoiceItem{
			ContractId: contractID,
			Period:     INVOICE_PERIOD,
			ItemType:   ITEM_TYPE_ELECTRICITY,
			Quantity:   elecQty,
			UnitPrice:  newUnitPrice,
		},
	})
	if err != nil {
		t.Fatal("UpdateInvoiceItem failed: " + err.Error())
	}

	expectedAmount := elecQty * newUnitPrice
	if response.InvoiceItem.Amount != expectedAmount {
		t.Fatalf("expected amount %.2f, got %.2f", expectedAmount, response.InvoiceItem.Amount)
	}

	// verify invoice total_amount was recalculated
	invoiceResp, err := service.GetInvoice(ctx, &types.GetInvoiceRequest{
		ContractId: contractID,
		Period:     INVOICE_PERIOD,
	})
	if err != nil {
		t.Fatal("GetInvoice failed: " + err.Error())
	}

	waterQty := WATER_READING_CURR - WATER_READING_PREV
	expectedTotal := PRICE + elecQty*newUnitPrice + waterQty*WATER_UNIT_PRICE
	if invoiceResp.Invoice.TotalAmount != expectedTotal {
		t.Fatalf("expected recalculated total %.2f, got %.2f", expectedTotal, invoiceResp.Invoice.TotalAmount)
	}
}

func TestUpdateInvoiceItem_MissingContractId(t *testing.T) {
	service := getServiceForInvoice()

	response, err := service.UpdateInvoiceItem(context.Background(), &types.UpdateInvoiceItemRequest{
		InvoiceItem: &types.InvoiceItem{
			Period:    INVOICE_PERIOD,
			ItemType:  ITEM_TYPE_ELECTRICITY,
			Quantity:  50,
			UnitPrice: ELEC_UNIT_PRICE,
		},
	})

	if err == nil {
		t.Fatal("expected error for missing contract_id, got nil")
	}
	if err.Error() != "UpdateInvoiceItem error: contract_id is required" {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if response != nil {
		t.Fatal("expected nil response")
	}
}

func TestUpdateInvoiceItem_MissingItemType(t *testing.T) {
	service := getServiceForInvoice()

	response, err := service.UpdateInvoiceItem(context.Background(), &types.UpdateInvoiceItemRequest{
		InvoiceItem: &types.InvoiceItem{
			ContractId: "some-id",
			Period:     INVOICE_PERIOD,
			Quantity:   50,
			UnitPrice:  ELEC_UNIT_PRICE,
		},
	})

	if err == nil {
		t.Fatal("expected error for missing item_type, got nil")
	}
	if err.Error() != "UpdateInvoiceItem error: item_type is required" {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	if response != nil {
		t.Fatal("expected nil response")
	}
}

func TestListInvoiceItems_HappyPath(t *testing.T) {
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

	response, err := service.ListInvoiceItems(ctx, &types.ListInvoiceItemsRequest{
		ContractId: contractID,
		Period:     INVOICE_PERIOD,
	})
	if err != nil {
		t.Fatal("ListInvoiceItems failed: " + err.Error())
	}

	if len(response.InvoiceItems) != 3 {
		t.Fatalf("expected 3 invoice items, got %d", len(response.InvoiceItems))
	}
}

func TestListInvoiceItems_Empty(t *testing.T) {
	service := getServiceForInvoice()

	response, err := service.ListInvoiceItems(context.Background(), &types.ListInvoiceItemsRequest{
		ContractId: "no-invoice",
		Period:     INVOICE_PERIOD,
	})
	if err != nil {
		t.Fatal("ListInvoiceItems failed: " + err.Error())
	}
	if len(response.InvoiceItems) != 0 {
		t.Fatalf("expected 0 items, got %d", len(response.InvoiceItems))
	}
}
