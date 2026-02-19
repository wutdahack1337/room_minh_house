package server

import (
	"context"
	"fmt"

	"github.com/wutdahack1337/room_minh_house/backend/internal/database"
	types "github.com/wutdahack1337/room_minh_house/backend/pkg/pb"
)

type Service struct {
	types.UnimplementedQueryServer
	dbService database.DatabaseService
}

func NewService(dbService database.DatabaseService) *Service {
	return &Service{
		dbService: dbService,
	}
}

func (s *Service) CreateRoom(ctx context.Context, request *types.CreateRoomRequest) (*types.CreateRoomResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("CreateRoom error: empty request")
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
	if request.Room.RoomNumber == 0 {
		return nil, fmt.Errorf("UpdateRoom error: room number is required")
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
	if err := s.dbService.DeleteRoom(ctx, request.RoomNumber); err != nil {
		return nil, err
	}
	return &types.DeleteRoomResponse{}, nil
}

func (s *Service) GetRoom(ctx context.Context, request *types.GetRoomRequest) (*types.GetRoomResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("GetRoom error: empty request")
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

func (s *Service) CheckHealth(ctx context.Context, request *types.CheckHealthRequest) (*types.CheckHealthResponse, error) {
	if err := s.dbService.Health(ctx); err != nil {
		return &types.CheckHealthResponse{
			Status: fmt.Sprintf("database error: %v", err),
		}, err
	}

	return &types.CheckHealthResponse{
		Status: "i'm fine",
	}, nil
}
