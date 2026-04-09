package server

import (
	"context"
	"fmt"

	types "github.com/wutdahack1337/room_minh_house/backend/pkg/pb"
)

func (s *Service) CreateRoom(ctx context.Context, request *types.CreateRoomRequest) (*types.CreateRoomResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("CreateRoom error: empty request")
	}
	if request.Room == nil {
		return nil, fmt.Errorf("CreateRoom error: empty room")
	}
	if request.Room.RoomNumber == 0 {
		return nil, fmt.Errorf("CreateRoom error: room number is required")
	}
	if request.Room.Price <= 0 {
		return nil, fmt.Errorf("CreateRoom error: room price must be greater than zero")
	}

	if err := s.dbService.CreateRoom(ctx, request.Room); err != nil {
		return nil, err
	}

	return &types.CreateRoomResponse{
		Room: request.Room,
	}, nil
}

func (s *Service) UpdateRoom(ctx context.Context, request *types.UpdateRoomRequest) (*types.UpdateRoomResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("UpdateRoom error: empty request")
	}
	if request.Room == nil {
		return nil, fmt.Errorf("UpdateRoom error: empty room")
	}
	if request.Room.RoomNumber == 0 {
		return nil, fmt.Errorf("UpdateRoom error: room number is required")
	}
	if request.Room.Price <= 0 {
		return nil, fmt.Errorf("UpdateRoom error: room price must be greater than zero")
	}

	if err := s.dbService.UpdateRoom(ctx, request.Room); err != nil {
		return nil, err
	}

	return &types.UpdateRoomResponse{
		Room: request.Room,
	}, nil
}

func (s *Service) DeleteRoom(ctx context.Context, request *types.DeleteRoomRequest) (*types.DeleteRoomResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("DeleteRoom error: empty request")
	}
	if request.RoomNumber == 0 {
		return nil, fmt.Errorf("DeleteRoom error: room number is required")
	}

	// cascade: delete all contracts (each call also deletes invoices + invoice_items)
	contracts, err := s.dbService.ListContractsByRoomNumber(ctx, request.RoomNumber)
	if err != nil {
		return nil, fmt.Errorf("DeleteRoom error: failed to list contracts: %w", err)
	}
	for _, c := range contracts {
		if _, err := s.DeleteContract(ctx, &types.DeleteContractRequest{ContractId: c.ContractId}); err != nil {
			return nil, fmt.Errorf("DeleteRoom error: failed to delete contract %s: %w", c.ContractId, err)
		}
	}

	// cascade: delete all meter readings
	if err := s.dbService.DeleteMeterReadingsByRoomNumber(ctx, request.RoomNumber); err != nil {
		return nil, fmt.Errorf("DeleteRoom error: failed to delete meter readings: %w", err)
	}

	if err := s.dbService.DeleteRoom(ctx, request.RoomNumber); err != nil {
		return nil, err
	}

	return &types.DeleteRoomResponse{}, nil
}

func (s *Service) GetRoom(ctx context.Context, request *types.GetRoomRequest) (*types.GetRoomResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("GetRoom error: empty request")
	}
	if request.RoomNumber == 0 {
		return nil, fmt.Errorf("GetRoom error: room number is required")
	}

	room, err := s.dbService.GetRoom(ctx, request.RoomNumber)
	if err != nil {
		return nil, err
	}

	return &types.GetRoomResponse{
		Room: room,
	}, nil
}

func (s *Service) ListRooms(ctx context.Context, request *types.ListRoomsRequest) (*types.ListRoomsResponse, error) {
	rooms, err := s.dbService.ListRooms(ctx)
	if err != nil {
		return nil, err
	}

	return &types.ListRoomsResponse{
		Rooms: rooms,
	}, nil
}
