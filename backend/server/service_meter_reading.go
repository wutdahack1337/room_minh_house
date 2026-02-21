package server

import (
	"context"
	"fmt"

	types "github.com/wutdahack1337/room_minh_house/backend/pkg/pb"
)

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

	if err := s.dbService.UpdateMeterReading(ctx, request.MeterReading); err != nil {
		return nil, err
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
