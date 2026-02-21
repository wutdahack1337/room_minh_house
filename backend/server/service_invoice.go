package server

import (
	"context"
	"fmt"
	"time"

	types "github.com/wutdahack1337/room_minh_house/backend/pkg/pb"
)

const (
	ITEM_TYPE_ROOM_RENT   = "room_rent"
	ITEM_TYPE_ELECTRICITY = "electricity"
	ITEM_TYPE_WATER       = "water"
)

// prevPeriod returns the previous month period string.
// format: MM-YYYY
func prevPeriod(period string) (string, error) {
	t, err := time.Parse("01-2006", period)
	if err != nil {
		return "", fmt.Errorf("invalid period format, expected MM-YYYY: %s", period)
	}
	prev := t.AddDate(0, -1, 0)
	return prev.Format("01-2006"), nil
}

func (s *Service) CreateInvoice(ctx context.Context, request *types.CreateInvoiceRequest) (*types.CreateInvoiceResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("CreateInvoice error: empty request")
	}
	if request.ContractId == "" {
		return nil, fmt.Errorf("CreateInvoice error: contract_id is required")
	}
	if request.Period == "" {
		return nil, fmt.Errorf("CreateInvoice error: period is required")
	}
	if request.ElectricityUnitPrice <= 0 {
		return nil, fmt.Errorf("CreateInvoice error: electricity_unit_price must be greater than 0")
	}
	if request.WaterUnitPrice <= 0 {
		return nil, fmt.Errorf("CreateInvoice error: water_unit_price must be greater than 0")
	}

	// validate period format
	if err := validatePeriod(request.Period); err != nil {
		return nil, fmt.Errorf("CreateInvoice error: %v", err)
	}

	// get contract to find room_number
	contract, err := s.dbService.GetContract(ctx, request.ContractId)
	if err != nil {
		return nil, err
	}

	// get room to find price
	room, err := s.dbService.GetRoom(ctx, contract.RoomNumber)
	if err != nil {
		return nil, err
	}

	// get current period meter reading
	currentReading, err := s.dbService.GetMeterReading(ctx, contract.RoomNumber, request.Period)
	if err != nil {
		return nil, fmt.Errorf("CreateInvoice error: meter reading for period %s not found", request.Period)
	}

	// get previous period meter reading
	prev, err := prevPeriod(request.Period)
	if err != nil {
		return nil, fmt.Errorf("CreateInvoice error: %v", err)
	}
	prevReading, err := s.dbService.GetMeterReading(ctx, contract.RoomNumber, prev)
	if err != nil {
		return nil, fmt.Errorf("CreateInvoice error: meter reading for previous period %s not found", prev)
	}

	// calculate quantities
	elecQty := currentReading.ElectricityReading - prevReading.ElectricityReading
	waterQty := currentReading.WaterReading - prevReading.WaterReading

	// calculate amounts
	roomAmount := room.Price
	elecAmount := elecQty * request.ElectricityUnitPrice
	waterAmount := waterQty * request.WaterUnitPrice
	totalAmount := roomAmount + elecAmount + waterAmount

	// build invoice
	invoice := &types.Invoice{
		ContractId:  request.ContractId,
		Period:      request.Period,
		TotalAmount: totalAmount,
		Status:      "unpaid",
	}

	if err := s.dbService.CreateInvoice(ctx, invoice); err != nil {
		return nil, err
	}

	// build invoice items
	items := []*types.InvoiceItem{
		{
			ContractId: request.ContractId,
			Period:     request.Period,
			ItemType:   ITEM_TYPE_ROOM_RENT,
			Quantity:   1,
			UnitPrice:  room.Price,
			Amount:     roomAmount,
		},
		{
			ContractId: request.ContractId,
			Period:     request.Period,
			ItemType:   ITEM_TYPE_ELECTRICITY,
			Quantity:   elecQty,
			UnitPrice:  request.ElectricityUnitPrice,
			Amount:     elecAmount,
		},
		{
			ContractId: request.ContractId,
			Period:     request.Period,
			ItemType:   ITEM_TYPE_WATER,
			Quantity:   waterQty,
			UnitPrice:  request.WaterUnitPrice,
			Amount:     waterAmount,
		},
	}

	for _, item := range items {
		if err := s.dbService.CreateInvoiceItem(ctx, item); err != nil {
			return nil, err
		}
	}

	return &types.CreateInvoiceResponse{
		Invoice:      invoice,
		InvoiceItems: items,
	}, nil
}

func (s *Service) UpdateInvoice(ctx context.Context, request *types.UpdateInvoiceRequest) (*types.UpdateInvoiceResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("UpdateInvoice error: empty request")
	}
	if request.Invoice == nil {
		return nil, fmt.Errorf("UpdateInvoice error: empty invoice")
	}
	if request.Invoice.ContractId == "" {
		return nil, fmt.Errorf("UpdateInvoice error: contract_id is required")
	}
	if request.Invoice.Period == "" {
		return nil, fmt.Errorf("UpdateInvoice error: period is required")
	}

	if err := s.dbService.UpdateInvoice(ctx, request.Invoice); err != nil {
		return nil, err
	}

	return &types.UpdateInvoiceResponse{Invoice: request.Invoice}, nil
}

func (s *Service) DeleteInvoice(ctx context.Context, request *types.DeleteInvoiceRequest) (*types.DeleteInvoiceResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("DeleteInvoice error: empty request")
	}
	if request.ContractId == "" {
		return nil, fmt.Errorf("DeleteInvoice error: contract_id is required")
	}
	if request.Period == "" {
		return nil, fmt.Errorf("DeleteInvoice error: period is required")
	}

	if err := s.dbService.DeleteInvoice(ctx, request.ContractId, request.Period); err != nil {
		return nil, err
	}

	return &types.DeleteInvoiceResponse{}, nil
}

func (s *Service) GetInvoice(ctx context.Context, request *types.GetInvoiceRequest) (*types.GetInvoiceResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("GetInvoice error: empty request")
	}
	if request.ContractId == "" {
		return nil, fmt.Errorf("GetInvoice error: contract_id is required")
	}
	if request.Period == "" {
		return nil, fmt.Errorf("GetInvoice error: period is required")
	}

	invoice, err := s.dbService.GetInvoice(ctx, request.ContractId, request.Period)
	if err != nil {
		return nil, err
	}

	return &types.GetInvoiceResponse{Invoice: invoice}, nil
}

func (s *Service) ListInvoices(ctx context.Context, request *types.ListInvoicesRequest) (*types.ListInvoicesResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("ListInvoices error: empty request")
	}
	if request.ContractId == "" {
		return nil, fmt.Errorf("ListInvoices error: contract_id is required")
	}

	invoices, err := s.dbService.ListInvoices(ctx, request.ContractId)
	if err != nil {
		return nil, err
	}

	return &types.ListInvoicesResponse{Invoices: invoices}, nil
}
