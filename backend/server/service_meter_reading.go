package server

import (
	"context"
	"fmt"
	"time"

	types "github.com/wutdahack1337/room_minh_house/backend/pkg/pb"
	"go.mongodb.org/mongo-driver/mongo"
)

func parsePeriod(period string) (time.Time, error) {
	t, err := time.Parse("01-2006", period)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid period format, expected MM-YYYY: %s", period)
	}
	return t, nil
}

func (s *Service) validateCreateMeterReadingOrderAndValue(ctx context.Context, meterReading *types.MeterReading) error {
	newPeriod, err := parsePeriod(meterReading.Period)
	if err != nil {
		return err
	}

	meterReadings, err := s.dbService.ListMeterReadings(ctx, meterReading.RoomNumber)
	if err != nil {
		return err
	}
	if len(meterReadings) == 0 {
		return nil
	}

	latest := time.Time{}
	var latestReading *types.MeterReading
	for _, reading := range meterReadings {
		t, err := parsePeriod(reading.Period)
		if err != nil {
			return err
		}
		if latest.IsZero() || t.After(latest) {
			latest = t
			latestReading = reading
		}
	}

	if !newPeriod.After(latest) {
		return fmt.Errorf("period must be greater than latest existing period (%s)", latest.Format("01-2006"))
	}

	if latestReading != nil {
		if meterReading.ElectricityReading < latestReading.ElectricityReading {
			return fmt.Errorf("electricity reading of new period must be greater than or equal to previous period (%s)", latestReading.Period)
		}
		if meterReading.WaterReading < latestReading.WaterReading {
			return fmt.Errorf("water reading of new period must be greater than or equal to previous period (%s)", latestReading.Period)
		}
	}

	return nil
}

func (s *Service) CreateMeterReading(ctx context.Context, request *types.CreateMeterReadingRequest) (*types.CreateMeterReadingResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("CreateMeterReading error: empty request")
	}
	if request.MeterReading == nil {
		return nil, fmt.Errorf("CreateMeterReading error: empty meter reading")
	}
	if request.MeterReading.RoomNumber == 0 {
		return nil, fmt.Errorf("CreateMeterReading error: room number is required")
	}
	if request.MeterReading.Period == "" {
		return nil, fmt.Errorf("CreateMeterReading error: period is required")
	}
	if err := validatePeriod(request.MeterReading.Period); err != nil {
		return nil, fmt.Errorf("CreateMeterReading error: %v", err)
	}
	if request.MeterReading.ElectricityReading < 0 {
		return nil, fmt.Errorf("CreateMeterReading error: electricity reading must be greater than or equal to 0")
	}
	if request.MeterReading.WaterReading < 0 {
		return nil, fmt.Errorf("CreateMeterReading error: water reading must be greater than or equal to 0")
	}
	if err := s.validateCreateMeterReadingOrderAndValue(ctx, request.MeterReading); err != nil {
		return nil, fmt.Errorf("CreateMeterReading error: %v", err)
	}

	if err := s.dbService.CreateMeterReading(ctx, request.MeterReading); err != nil {
		return nil, err
	}

	return &types.CreateMeterReadingResponse{MeterReading: request.MeterReading}, nil
}

func (s *Service) UpdateMeterReading(ctx context.Context, request *types.UpdateMeterReadingRequest) (*types.UpdateMeterReadingResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("UpdateMeterReading error: empty request")
	}
	if request.MeterReading == nil {
		return nil, fmt.Errorf("UpdateMeterReading error: empty meter reading")
	}
	if request.MeterReading.RoomNumber == 0 {
		return nil, fmt.Errorf("UpdateMeterReading error: room number is required")
	}
	if request.MeterReading.Period == "" {
		return nil, fmt.Errorf("UpdateMeterReading error: period is required")
	}
	if err := validatePeriod(request.MeterReading.Period); err != nil {
		return nil, fmt.Errorf("UpdateMeterReading error: %v", err)
	}
	if request.MeterReading.ElectricityReading < 0 {
		return nil, fmt.Errorf("UpdateMeterReading error: electricity reading must be greater than or equal to 0")
	}
	if request.MeterReading.WaterReading < 0 {
		return nil, fmt.Errorf("UpdateMeterReading error: water reading must be greater than or equal to 0")
	}

	if err := s.dbService.UpdateMeterReading(ctx, request.MeterReading); err != nil {
		return nil, err
	}

	// cascade: recalculate invoice_items for any invoice in this period
	contracts, err := s.dbService.ListContractsByRoomNumber(ctx, request.MeterReading.RoomNumber)
	if err != nil {
		return nil, fmt.Errorf("UpdateMeterReading error: failed to list contracts: %w", err)
	}

	period := request.MeterReading.Period
	prev, err := prevPeriod(period)
	if err != nil {
		return nil, fmt.Errorf("UpdateMeterReading error: %v", err)
	}

	for _, c := range contracts {
		// only cascade if an invoice exists for this period
		_, err := s.dbService.GetInvoice(ctx, c.ContractId, period)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				continue // no invoice for this period, nothing to cascade
			}
			return nil, fmt.Errorf("UpdateMeterReading error: failed to get invoice: %w", err)
		}

		// get previous meter reading for qty calculation
		prevReading, err := s.dbService.GetMeterReading(ctx, request.MeterReading.RoomNumber, prev)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				continue // no previous reading yet, can't compute delta
			}
			return nil, fmt.Errorf("UpdateMeterReading error: failed to get previous meter reading: %w", err)
		}

		elecQty := request.MeterReading.ElectricityReading - prevReading.ElectricityReading
		waterQty := request.MeterReading.WaterReading - prevReading.WaterReading

		// get existing items to preserve unit prices
		items, err := s.dbService.ListInvoiceItems(ctx, c.ContractId, period)
		if err != nil {
			return nil, fmt.Errorf("UpdateMeterReading error: failed to list invoice items: %w", err)
		}

		unitPrices := make(map[string]float64)
		for _, item := range items {
			unitPrices[item.ItemType] = item.UnitPrice
		}

		elecUnitPrice := unitPrices[ITEM_TYPE_ELECTRICITY]
		waterUnitPrice := unitPrices[ITEM_TYPE_WATER]

		if elecUnitPrice > 0 {
			if err := s.dbService.UpdateInvoiceItem(ctx, &types.InvoiceItem{
				ContractId: c.ContractId,
				Period:     period,
				ItemType:   ITEM_TYPE_ELECTRICITY,
				Quantity:   elecQty,
				UnitPrice:  elecUnitPrice,
				Amount:     elecQty * elecUnitPrice,
			}); err != nil {
				return nil, fmt.Errorf("UpdateMeterReading error: failed to update electricity item: %w", err)
			}
		}

		if waterUnitPrice > 0 {
			if err := s.dbService.UpdateInvoiceItem(ctx, &types.InvoiceItem{
				ContractId: c.ContractId,
				Period:     period,
				ItemType:   ITEM_TYPE_WATER,
				Quantity:   waterQty,
				UnitPrice:  waterUnitPrice,
				Amount:     waterQty * waterUnitPrice,
			}); err != nil {
				return nil, fmt.Errorf("UpdateMeterReading error: failed to update water item: %w", err)
			}
		}

		if err := s.dbService.UpdateInvoiceTotalAmount(ctx, c.ContractId, period); err != nil {
			return nil, fmt.Errorf("UpdateMeterReading error: failed to update invoice total: %w", err)
		}
	}

	return &types.UpdateMeterReadingResponse{MeterReading: request.MeterReading}, nil
}

func (s *Service) GetMeterReading(ctx context.Context, request *types.GetMeterReadingRequest) (*types.GetMeterReadingResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("GetMeterReading error: empty request")
	}
	if request.RoomNumber == 0 {
		return nil, fmt.Errorf("GetMeterReading error: room number is required")
	}
	if request.Period == "" {
		return nil, fmt.Errorf("GetMeterReading error: period is required")
	}
	if err := validatePeriod(request.Period); err != nil {
		return nil, fmt.Errorf("GetMeterReading error: %v", err)
	}

	meterReading, err := s.dbService.GetMeterReading(ctx, request.RoomNumber, request.Period)
	if err != nil {
		return nil, err
	}

	return &types.GetMeterReadingResponse{MeterReading: meterReading}, nil
}

func (s *Service) ListMeterReadings(ctx context.Context, request *types.ListMeterReadingsRequest) (*types.ListMeterReadingsResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("ListMeterReadings error: empty request")
	}
	if request.RoomNumber == 0 {
		return nil, fmt.Errorf("ListMeterReadings error: room number is required")
	}

	meterReadings, err := s.dbService.ListMeterReadings(ctx, request.RoomNumber)
	if err != nil {
		return nil, err
	}

	return &types.ListMeterReadingsResponse{MeterReadings: meterReadings}, nil
}
