package server

import (
	"context"
	"testing"
	"time"

	"github.com/wutdahack1337/room_minh_house/backend/internal/database"
	types "github.com/wutdahack1337/room_minh_house/backend/pkg/pb"
)

func TestCreateContract_HappyPath(t *testing.T) {
	ctx := context.Background()
	service := getServiceForContract()

	// Create a valid user and room first
	user := &types.User{
		Pid:         P_ID,
		PhoneNumber: PHONE_NUMBER,
		Role:        ROLE,
		IdCard:      newTestIdCard(),
		Account:     newTestAccount(),
	}
	room := &types.Room{
		RoomNumber:  ROOM_NUMBER,
		Price:       PRICE,
		IsAvailable: true,
	}

	_, err := service.CreateUser(ctx, &types.CreateUserRequest{
		User: user,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = service.CreateRoom(ctx, &types.CreateRoomRequest{
		Room: room,
	})
	if err != nil {
		t.Fatal(err)
	}

	contract := &types.Contract{
		UserPid:        P_ID,
		RoomNumber:     ROOM_NUMBER,
		StartDate:      time.Now().UTC().Format(time.RFC3339),
		DurationMonths: 3,
	}

	response, err := service.CreateContract(ctx, &types.CreateContractRequest{
		Contract: contract,
	})
	if err != nil {
		t.Fatal(err)
	}

	if response == nil {
		t.Fatal("CreateContract returned nil response")
	}
	if response.Contract.ContractId == "" {
		t.Fatal("CreateContract did not generate contract_id")
	}
	if response.Contract.Status != "active" {
		t.Fatalf("expected status 'active', got '%s'", response.Contract.Status)
	}
	if response.Contract.EndDate == "" {
		t.Fatal("CreateContract did not calculate end_date")
	}
}

func TestCreateContract_MissingUserPid(t *testing.T) {
	contract := &types.Contract{
		RoomNumber:     ROOM_NUMBER,
		StartDate:      time.Now().UTC().Format(time.RFC3339),
		DurationMonths: 3,
	}

	service := getServiceForContract()
	response, err := service.CreateContract(context.Background(), &types.CreateContractRequest{
		Contract: contract,
	})

	if err == nil {
		t.Fatal("expected error for missing user_pid, got nil")
	}

	if response != nil {
		t.Fatal("expected nil response for missing user_pid, got non-nil")
	}
}

func TestCreateContract_MissingRoomNumber(t *testing.T) {
	contract := &types.Contract{
		UserPid:        P_ID,
		StartDate:      time.Now().UTC().Format(time.RFC3339),
		DurationMonths: 3,
	}

	service := getServiceForContract()
	response, err := service.CreateContract(context.Background(), &types.CreateContractRequest{
		Contract: contract,
	})

	if err == nil {
		t.Fatal("expected error for missing room_number, got nil")
	}

	if response != nil {
		t.Fatal("expected nil response for missing room_number, got non-nil")
	}
}

func TestCreateContract_UserNotFound(t *testing.T) {
	contract := &types.Contract{
		UserPid:        "i'm not exist",
		RoomNumber:     ROOM_NUMBER,
		StartDate:      time.Now().UTC().Format(time.RFC3339),
		DurationMonths: 3,
	}

	service := getServiceForContract()
	response, err := service.CreateContract(context.Background(), &types.CreateContractRequest{
		Contract: contract,
	})

	if err == nil {
		t.Fatal("expected error for non-existent user, got nil")
	}
	if err.Error() != "database error: user with pid i'm not exist not found" {
		t.Fatalf("expected 'database error: user with pid \"i'm not exist\" not found', got '%s'", err.Error())
	}

	if response != nil {
		t.Fatal("expected nil response for non-existent user")
	}
}

func TestCreateContract_RoomNotFound(t *testing.T) {
	ctx := context.Background()
	service := getServiceForContract()

	// create a valid user
	user := &types.User{
		Pid:         P_ID,
		PhoneNumber: PHONE_NUMBER,
		Role:        ROLE,
		IdCard:      newTestIdCard(),
		Account:     newTestAccount(),
	}

	_, err := service.CreateUser(ctx, &types.CreateUserRequest{
		User: user,
	})
	if err != nil {
		t.Fatal("Failed to create user: " + err.Error())
	}

	contract := &types.Contract{
		UserPid:        P_ID,
		RoomNumber:     999_999,
		StartDate:      time.Now().UTC().Format(time.RFC3339),
		DurationMonths: 3,
	}

	response, err := service.CreateContract(ctx, &types.CreateContractRequest{
		Contract: contract,
	})

	if err == nil {
		t.Fatal("expected error for non-existent room, got nil")
	}
	if err.Error() != "database error: room with number 999999 not found" {
		t.Fatalf("expected 'database error: room with number 999999 not found', got '%s'", err.Error())
	}

	if response != nil {
		t.Fatal("expected nil response for non-existent room")
	}
}

func TestUpdateContract_HappyPath(t *testing.T) {
	ctx := context.Background()

	service := getServiceForContract()

	// Create a valid user and room first
	user := &types.User{
		Pid:         P_ID,
		PhoneNumber: PHONE_NUMBER,
		Role:        ROLE,
		IdCard:      newTestIdCard(),
		Account:     newTestAccount(),
	}
	room := &types.Room{
		RoomNumber:  ROOM_NUMBER,
		Price:       PRICE,
		IsAvailable: true,
	}

	_, err := service.CreateUser(ctx, &types.CreateUserRequest{
		User: user,
	})
	if err != nil {
		t.Fatal("Failed to create user: " + err.Error())
	}
	_, err = service.CreateRoom(ctx, &types.CreateRoomRequest{
		Room: room,
	})
	if err != nil {
		t.Fatal("Failed to create room: " + err.Error())
	}

	contract := &types.Contract{
		UserPid:        P_ID,
		RoomNumber:     ROOM_NUMBER,
		StartDate:      time.Now().UTC().Format(time.RFC3339),
		DurationMonths: 3,
	}
	createResponse, _ := service.CreateContract(ctx, &types.CreateContractRequest{
		Contract: contract,
	})

	contractId := createResponse.Contract.ContractId

	// update contract status to "expired"
	contract.ContractId = contractId
	contract.Status = "expired"

	_, err = service.UpdateContract(ctx, &types.UpdateContractRequest{
		Contract: contract,
	})
	if err != nil {
		t.Fatal("UpdateContract failed: " + err.Error())
	}

	response, err := service.GetContract(ctx, &types.GetContractRequest{
		ContractId: contractId,
	})

	if response == nil {
		t.Fatal("GetContract returned nil response")
	}
	if response.Contract.ContractId != contractId {
		t.Fatalf("expected contract_id %s, got %s", contractId, response.Contract.ContractId)
	}
	if response.Contract.Status != "expired" {
		t.Fatalf("expected status 'expired', got '%s'", response.Contract.Status)
	}
}

func TestGetContract_HappyPath(t *testing.T) {
	ctx := context.Background()
	service := getServiceForContract()

	// Create a valid user and room first
	user := &types.User{
		Pid:         P_ID,
		PhoneNumber: PHONE_NUMBER,
		Role:        ROLE,
		IdCard:      newTestIdCard(),
		Account:     newTestAccount(),
	}
	room := &types.Room{
		RoomNumber:  ROOM_NUMBER,
		Price:       666_666,
		IsAvailable: true,
	}
	_, err := service.CreateUser(ctx, &types.CreateUserRequest{
		User: user,
	})
	if err != nil {
		t.Fatal("Failed to create user: " + err.Error())
	}

	_, err = service.CreateRoom(ctx, &types.CreateRoomRequest{
		Room: room,
	})
	if err != nil {
		t.Fatal("Failed to create room: " + err.Error())
	}

	contract := &types.Contract{
		UserPid:        P_ID,
		RoomNumber:     ROOM_NUMBER,
		StartDate:      time.Now().UTC().Format(time.RFC3339),
		DurationMonths: 3,
	}

	createResponse, _ := service.CreateContract(ctx, &types.CreateContractRequest{
		Contract: contract,
	})

	contractId := createResponse.Contract.ContractId

	getResponse, err := service.GetContract(ctx, &types.GetContractRequest{
		ContractId: contractId,
	})
	if err != nil {
		t.Fatal("GetContract failed: " + err.Error())
	}

	if getResponse == nil {
		t.Fatal("GetContract returned nil response")
	}
	if getResponse.Contract.ContractId != contractId {
		t.Fatalf("expected contract_id %s, got %s", contractId, getResponse.Contract.ContractId)
	}
}

func TestGetContract_NotFound(t *testing.T) {
	service := getServiceForContract()

	response, err := service.GetContract(context.Background(), &types.GetContractRequest{
		ContractId: "non-existent-contract-id",
	})
	if err == nil {
		t.Fatal("expected error for non-existent contract, got nil")
	}

	if response != nil {
		t.Fatal("expected nil response for non-existent contract")
	}
}

func TestDeleteContract_HappyPath(t *testing.T) {
	ctx := context.Background()
	service := getServiceForContract()

	// Create a valid user and room first
	user := &types.User{
		Pid:         P_ID,
		PhoneNumber: PHONE_NUMBER,
		Role:        ROLE,
		IdCard:      newTestIdCard(),
		Account:     newTestAccount(),
	}
	room := &types.Room{
		RoomNumber:  ROOM_NUMBER,
		Price:       666_666,
		IsAvailable: true,
	}

	_, err := service.CreateUser(ctx, &types.CreateUserRequest{
		User: user,
	})
	if err != nil {
		t.Fatal("Failed to create user: " + err.Error())
	}

	_, err = service.CreateRoom(ctx, &types.CreateRoomRequest{
		Room: room,
	})
	if err != nil {
		t.Fatal("Failed to create room: " + err.Error())
	}

	contract := &types.Contract{
		UserPid:        P_ID,
		RoomNumber:     ROOM_NUMBER,
		StartDate:      time.Now().UTC().Format(time.RFC3339),
		DurationMonths: 3,
	}

	createResponse, _ := service.CreateContract(ctx, &types.CreateContractRequest{
		Contract: contract,
	})

	contractId := createResponse.Contract.ContractId

	_, err = service.DeleteContract(ctx, &types.DeleteContractRequest{
		ContractId: contractId,
	})
	if err != nil {
		t.Fatal("DeleteContract failed: " + err.Error())
	}

	// Try to get deleted contract
	getResp, err := service.GetContract(ctx, &types.GetContractRequest{
		ContractId: contractId,
	})
	if err == nil {
		t.Fatal("expected error for deleted contract, got nil")
	}

	if getResp != nil {
		t.Fatal("expected nil response for deleted contract")
	}
}

func TestListContracts_HappyPath(t *testing.T) {
	ctx := context.Background()
	service := getServiceForContract()

	users := []*types.User{
		{Pid: "user1", PhoneNumber: "0111111111", Role: ROLE, IdCard: newTestIdCard(), Account: newTestAccount()},
		{Pid: "user2", PhoneNumber: "0222222222", Role: ROLE, IdCard: newTestIdCard(), Account: newTestAccount()},
		{Pid: "user3", PhoneNumber: "0333333333", Role: ROLE, IdCard: newTestIdCard(), Account: newTestAccount()},
	}
	for _, user := range users {
		_, err := service.CreateUser(ctx, &types.CreateUserRequest{User: user})
		if err != nil {
			t.Fatal("Failed to create user: " + err.Error())
		}
	}

	rooms := []*types.Room{
		{RoomNumber: 666, Price: 666_666, IsAvailable: true},
		{RoomNumber: 693, Price: 693_693, IsAvailable: true},
		{RoomNumber: 999, Price: 999_999, IsAvailable: true},
	}
	for _, room := range rooms {
		_, _ = service.CreateRoom(ctx, &types.CreateRoomRequest{Room: room})
	}

	startDate := time.Now().UTC().Format(time.RFC3339)
	contracts := []*types.Contract{
		{
			UserPid:        "user1",
			RoomNumber:     666,
			StartDate:      startDate,
			DurationMonths: 3,
		},
		{
			UserPid:        "user2",
			RoomNumber:     693,
			StartDate:      startDate,
			DurationMonths: 6,
		},
		{
			UserPid:        "user3",
			RoomNumber:     999,
			StartDate:      startDate,
			DurationMonths: 3,
		},
	}

	for _, contract := range contracts {
		_, _ = service.CreateContract(ctx, &types.CreateContractRequest{
			Contract: contract,
		})
	}

	response, err := service.ListContracts(ctx, &types.ListContractsRequest{})
	if err != nil {
		t.Fatal("ListContracts failed: " + err.Error())
	}

	if response == nil {
		t.Fatal("ListContracts returned nil response")
	}
	if len(response.Contracts) != 3 {
		t.Fatalf("expected 3 contracts, got %d", len(response.Contracts))
	}
}

func TestListContracts_Empty(t *testing.T) {
	service := getServiceForContract()

	response, err := service.ListContracts(context.Background(), &types.ListContractsRequest{})
	if err != nil {
		t.Fatal("ListContracts failed: " + err.Error())
	}

	if response == nil {
		t.Fatal("ListContracts returned nil response")
	}
	if len(response.Contracts) != 0 {
		t.Fatalf("expected 0 contracts, got %d", len(response.Contracts))
	}
}

// helper
func getServiceForContract() *Service {
	db, err := connectDatabase("mongodb://localhost:27017/?directConnection=true", "room_minh_house")
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// removes all documents and indexes
	db.Collection("users").Drop(context.Background())
	db.Collection("rooms").Drop(context.Background())
	db.Collection("contracts").Drop(context.Background())

	dbService, err := database.NewMongoService(db)
	if err != nil {
		panic("Failed to create database service: " + err.Error())
	}

	return NewService(dbService)
}
