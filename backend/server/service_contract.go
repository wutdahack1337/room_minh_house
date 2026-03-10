package server

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	types "github.com/wutdahack1337/room_minh_house/backend/pkg/pb"
)

func (s *Service) CreateContract(ctx context.Context, request *types.CreateContractRequest) (*types.CreateContractResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("CreateContract error: empty request")
	}
	if request.Contract == nil {
		return nil, fmt.Errorf("CreateContract error: empty contract")
	}
	if request.Contract.UserPid == "" {
		return nil, fmt.Errorf("CreateContract error: user pid is required")
	}
	if request.Contract.RoomNumber == 0 {
		return nil, fmt.Errorf("CreateContract error: room number is required")
	}
	if request.Contract.StartDate == "" {
		return nil, fmt.Errorf("CreateContract error: start date is required")
	}
	if request.Contract.DurationMonths <= 0 {
		return nil, fmt.Errorf("CreateContract error: duration months must be greater than 0")
	}

	// generate UUID
	request.Contract.ContractId = uuid.NewString()
	// set initial status
	request.Contract.Status = "active"

	// calculate end_date
	startDate, err := time.Parse(time.RFC3339, request.Contract.StartDate)
	if err != nil {
		return nil, fmt.Errorf("CreateContract error: invalid start_date format, must be ISO 8601")
	}
	endDate := startDate.AddDate(0, int(request.Contract.DurationMonths), 0)
	request.Contract.EndDate = endDate.Format(time.RFC3339)

	if err := s.dbService.CreateContract(ctx, request.Contract); err != nil {
		return nil, err
	}

	return &types.CreateContractResponse{Contract: request.Contract}, nil
}

func (s *Service) UpdateContract(ctx context.Context, request *types.UpdateContractRequest) (*types.UpdateContractResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("UpdateContract error: empty request")
	}
	if request.Contract == nil {
		return nil, fmt.Errorf("UpdateContract error: empty contract")
	}
	if request.Contract.ContractId == "" {
		return nil, fmt.Errorf("UpdateContract error: contract id is required")
	}
	if request.Contract.UserPid == "" {
		return nil, fmt.Errorf("UpdateContract error: user pid is required")
	}
	if request.Contract.RoomNumber == 0 {
		return nil, fmt.Errorf("UpdateContract error: room number is required")
	}
	if request.Contract.StartDate == "" {
		return nil, fmt.Errorf("UpdateContract error: start date is required")
	}
	if request.Contract.DurationMonths <= 0 {
		return nil, fmt.Errorf("UpdateContract error: duration months must be greater than 0")
	}

	// recalculate end_date from start_date + duration_months
	startDate, err := time.Parse(time.RFC3339, request.Contract.StartDate)
	if err != nil {
		return nil, fmt.Errorf("UpdateContract error: invalid start_date format, must be ISO 8601")
	}
	request.Contract.EndDate = startDate.AddDate(0, int(request.Contract.DurationMonths), 0).Format(time.RFC3339)

	if err := s.dbService.UpdateContract(ctx, request.Contract); err != nil {
		return nil, err
	}

	return &types.UpdateContractResponse{Contract: request.Contract}, nil
}

func (s *Service) DeleteContract(ctx context.Context, request *types.DeleteContractRequest) (*types.DeleteContractResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("DeleteContract error: empty request")
	}
	if request.ContractId == "" {
		return nil, fmt.Errorf("DeleteContract error: contract id is required")
	}

	// cascade: delete all invoices (each call also deletes invoice_items)
	if err := s.dbService.DeleteInvoicesByContractId(ctx, request.ContractId); err != nil {
		return nil, fmt.Errorf("DeleteContract error: failed to delete invoices: %w", err)
	}

	if err := s.dbService.DeleteContract(ctx, request.ContractId); err != nil {
		return nil, err
	}

	return &types.DeleteContractResponse{}, nil
}

func (s *Service) GetContract(ctx context.Context, request *types.GetContractRequest) (*types.GetContractResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("GetContract error: empty request")
	}
	if request.ContractId == "" {
		return nil, fmt.Errorf("GetContract error: contract id is required")
	}

	contract, err := s.dbService.GetContract(ctx, request.ContractId)
	if err != nil {
		return nil, err
	}

	return &types.GetContractResponse{Contract: contract}, nil
}

func (s *Service) ListContracts(ctx context.Context, request *types.ListContractsRequest) (*types.ListContractsResponse, error) {
	contracts, err := s.dbService.ListContracts(ctx)
	if err != nil {
		return nil, err
	}

	return &types.ListContractsResponse{Contracts: contracts}, nil
}
