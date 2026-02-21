package database

import (
	"context"

	types "github.com/wutdahack1337/room_minh_house/backend/pkg/pb"
)

type Database interface {
	CreateRoom(ctx context.Context, room *types.Room) error
	UpdateRoom(ctx context.Context, room *types.Room) error
	DeleteRoom(ctx context.Context, roomNumber int32) error
	GetRoom(ctx context.Context, roomNumber int32) (*types.Room, error)
	ListRooms(ctx context.Context) ([]*types.Room, error)

	CreateUser(ctx context.Context, user *types.User) error
	UpdateUser(ctx context.Context, user *types.User) error
	DeleteUser(ctx context.Context, pid string) error
	GetUser(ctx context.Context, pid string) (*types.User, error)
	ListUsers(ctx context.Context) ([]*types.User, error)

	CreateContract(ctx context.Context, contract *types.Contract) error
	UpdateContract(ctx context.Context, contract *types.Contract) error
	DeleteContract(ctx context.Context, contractId string) error
	GetContract(ctx context.Context, contractId string) (*types.Contract, error)
	ListContracts(ctx context.Context) ([]*types.Contract, error)

	Health(ctx context.Context) error
}
