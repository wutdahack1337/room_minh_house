package server

import (
	"context"
	"fmt"

	types "github.com/wutdahack1337/room_minh_house/backend/pkg/pb"
)

func (s *Service) CreateUser(ctx context.Context, request *types.CreateUserRequest) (*types.CreateUserResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("CreateUser error: empty request")
	}
	if request.User == nil {
		return nil, fmt.Errorf("CreateUser error: empty user")
	}
	if request.User.Pid == "" {
		return nil, fmt.Errorf("CreateUser error: pid is required")
	}
	if request.User.Role == "" {
		return nil, fmt.Errorf("CreateUser error: role is required")
	}

	card := request.User.IdCard
	if card == nil {
		return nil, fmt.Errorf("CreateUser error: id card is required")
	}
	if card.Name == "" {
		return nil, fmt.Errorf("CreateUser error: id card name is required")
	}
	if card.Sex == "" {
		return nil, fmt.Errorf("CreateUser error: id card sex is required")
	}
	if card.DateOfBirth == "" {
		return nil, fmt.Errorf("CreateUser error: id card date of birth is required")
	}
	if card.PlaceOfBirth == "" {
		return nil, fmt.Errorf("CreateUser error: id card place of birth is required")
	}

	if err := s.dbService.CreateUser(ctx, request.User); err != nil {
		return nil, err
	}

	return &types.CreateUserResponse{User: request.User}, nil
}

func (s *Service) UpdateUser(ctx context.Context, request *types.UpdateUserRequest) (*types.UpdateUserResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("UpdateUser error: empty request")
	}
	if request.User == nil {
		return nil, fmt.Errorf("UpdateUser error: empty user")
	}
	if request.User.Pid == "" {
		return nil, fmt.Errorf("UpdateUser error: pid is required")
	}
	if request.User.Role == "" {
		return nil, fmt.Errorf("UpdateUser error: role is required")
	}

	card := request.User.IdCard
	if card == nil {
		return nil, fmt.Errorf("UpdateUser error: id card is required")
	}
	if card.Name == "" {
		return nil, fmt.Errorf("UpdateUser error: id card name is required")
	}
	if card.Sex == "" {
		return nil, fmt.Errorf("UpdateUser error: id card sex is required")
	}
	if card.DateOfBirth == "" {
		return nil, fmt.Errorf("UpdateUser error: id card date of birth is required")
	}
	if card.PlaceOfBirth == "" {
		return nil, fmt.Errorf("UpdateUser error: id card place of birth is required")
	}

	if err := s.dbService.UpdateUser(ctx, request.User); err != nil {
		return nil, err
	}

	return &types.UpdateUserResponse{User: request.User}, nil
}

func (s *Service) DeleteUser(ctx context.Context, request *types.DeleteUserRequest) (*types.DeleteUserResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("DeleteUser error: empty request")
	}
	if request.Pid == "" {
		return nil, fmt.Errorf("DeleteUser error: pid is required")
	}

	// cascade: delete all contracts (each call also deletes invoices + invoice_items)
	contracts, err := s.dbService.ListContractsByUserPid(ctx, request.Pid)
	if err != nil {
		return nil, fmt.Errorf("DeleteUser error: failed to list contracts: %w", err)
	}
	for _, c := range contracts {
		if _, err := s.DeleteContract(ctx, &types.DeleteContractRequest{ContractId: c.ContractId}); err != nil {
			return nil, fmt.Errorf("DeleteUser error: failed to delete contract %s: %w", c.ContractId, err)
		}
	}

	if err := s.dbService.DeleteUser(ctx, request.Pid); err != nil {
		return nil, err
	}

	return &types.DeleteUserResponse{}, nil
}

func (s *Service) GetUser(ctx context.Context, request *types.GetUserRequest) (*types.GetUserResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("GetUser error: empty request")
	}
	if request.Pid == "" {
		return nil, fmt.Errorf("GetUser error: pid is required")
	}

	user, err := s.dbService.GetUser(ctx, request.Pid)
	if err != nil {
		return nil, err
	}

	return &types.GetUserResponse{User: user}, nil
}

func (s *Service) ListUsers(ctx context.Context, request *types.ListUsersRequest) (*types.ListUsersResponse, error) {
	users, err := s.dbService.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	return &types.ListUsersResponse{Users: users}, nil
}
