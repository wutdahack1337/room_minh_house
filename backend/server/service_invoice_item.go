package server

import (
	"context"
	"fmt"

	types "github.com/wutdahack1337/room_minh_house/backend/pkg/pb"
)

func (s *Service) UpdateInvoiceItem(ctx context.Context, request *types.UpdateInvoiceItemRequest) (*types.UpdateInvoiceItemResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("UpdateInvoiceItem error: empty request")
	}
	if request.InvoiceItem == nil {
		return nil, fmt.Errorf("UpdateInvoiceItem error: empty invoice_item")
	}
	if request.InvoiceItem.ContractId == "" {
		return nil, fmt.Errorf("UpdateInvoiceItem error: contract_id is required")
	}
	if request.InvoiceItem.Period == "" {
		return nil, fmt.Errorf("UpdateInvoiceItem error: period is required")
	}
	if request.InvoiceItem.ItemType == "" {
		return nil, fmt.Errorf("UpdateInvoiceItem error: item_type is required")
	}
	if request.InvoiceItem.UnitPrice <= 0 {
		return nil, fmt.Errorf("UpdateInvoiceItem error: unit_price must be greater than 0")
	}

	// validate period format
	if err := validatePeriod(request.InvoiceItem.Period); err != nil {
		return nil, fmt.Errorf("UpdateInvoiceItem error: %v", err)
	}

	// recalculate amount from quantity * unit_price
	request.InvoiceItem.Amount = request.InvoiceItem.Quantity * request.InvoiceItem.UnitPrice

	if err := s.dbService.UpdateInvoiceItem(ctx, request.InvoiceItem); err != nil {
		return nil, err
	}

	// recalculate and update invoice total_amount
	if err := s.dbService.UpdateInvoiceTotalAmount(ctx, request.InvoiceItem.ContractId, request.InvoiceItem.Period); err != nil {
		return nil, err
	}

	return &types.UpdateInvoiceItemResponse{InvoiceItem: request.InvoiceItem}, nil
}

func (s *Service) ListInvoiceItems(ctx context.Context, request *types.ListInvoiceItemsRequest) (*types.ListInvoiceItemsResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("ListInvoiceItems error: empty request")
	}
	if request.ContractId == "" {
		return nil, fmt.Errorf("ListInvoiceItems error: contract_id is required")
	}
	if request.Period == "" {
		return nil, fmt.Errorf("ListInvoiceItems error: period is required")
	}

	// validate period format
	if err := validatePeriod(request.Period); err != nil {
		return nil, fmt.Errorf("ListInvoiceItems error: %v", err)
	}

	items, err := s.dbService.ListInvoiceItems(ctx, request.ContractId, request.Period)
	if err != nil {
		return nil, err
	}

	return &types.ListInvoiceItemsResponse{InvoiceItems: items}, nil
}
