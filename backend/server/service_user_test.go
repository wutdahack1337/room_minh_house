package server

import (
	"context"
	"testing"

	"github.com/wutdahack1337/room_minh_house/backend/internal/database"
	types "github.com/wutdahack1337/room_minh_house/backend/pkg/pb"
)

func TestCreateUser_HappyPath(t *testing.T) {
	user := &types.User{
		Pid:         "0123456789",
		PhoneNumber: "0949447156",
	}

	service := getUserService()
	response, err := service.CreateUser(context.Background(), &types.CreateUserRequest{
		User: user,
	})

	if err != nil {
		t.Fatal(err)
	}

	t.Log(response)
	if response == nil {
		t.Fatal("CreateUser returned nil response")
	}
}

func TestCreateUser_DuplicatePid(t *testing.T) {
	ctx := context.Background()

	user := &types.User{
		Pid:         "0123456789",
		PhoneNumber: "0949447156",
	}

	service := getUserService()
	_, _ = service.CreateUser(ctx, &types.CreateUserRequest{
		User: user,
	})

	response, err := service.CreateUser(ctx, &types.CreateUserRequest{
		User: user,
	})

	t.Log(err)
	if err == nil {
		t.Fatal("expected error for duplicate pid, got nil")
	}

	if response != nil {
		t.Fatal("expected nil response for duplicate pid, got non-nil")
	}
}

func TestCreateUser_MissingPid(t *testing.T) {
	service := getUserService()
	response, err := service.CreateUser(context.Background(), &types.CreateUserRequest{
		User: &types.User{
			PhoneNumber: "0949447156",
		},
	})

	t.Log(err)
	if err == nil {
		t.Fatal("expected error for empty pid, got nil")
	}

	if response != nil {
		t.Fatal("expected nil response for empty pid, got non-nil")
	}
}

func TestUpdateUser_HappyPath(t *testing.T) {
	ctx := context.Background()
	service := getUserService()

	user := &types.User{
		Pid:         "0123456789",
		PhoneNumber: "0949447156",
	}

	_, _ = service.CreateUser(ctx, &types.CreateUserRequest{
		User: user,
	})

	user.PhoneNumber = "0918441887"
	_, err := service.UpdateUser(ctx, &types.UpdateUserRequest{
		User: user,
	})

	if err != nil {
		t.Fatal("UpdateUser failed: " + err.Error())
	}

	response, err := service.GetUser(ctx, &types.GetUserRequest{
		Pid: user.Pid,
	})

	if err != nil {
		t.Fatal("GetUser failed: " + err.Error())
	}

	if response.User.PhoneNumber != user.PhoneNumber {
		t.Fatalf("expected phone number %s, got %s", user.PhoneNumber, response.User.PhoneNumber)
	}
}

func TestUpdateUser_MissingPid(t *testing.T) {
	user := &types.User{
		PhoneNumber: "0949447156",
	}

	service := getUserService()
	response, err := service.UpdateUser(context.Background(), &types.UpdateUserRequest{
		User: user,
	})

	t.Log(err)
	if err == nil {
		t.Fatal("expected error for missing pid, got nil")
	}

	if response != nil {
		t.Fatal("expected nil response for missing pid, got non-nil")
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	user := &types.User{
		Pid:         "9999999999",
		PhoneNumber: "0949447156",
	}

	service := getUserService()
	response, err := service.UpdateUser(context.Background(), &types.UpdateUserRequest{
		User: user,
	})

	t.Log(err)
	if err == nil {
		t.Fatal("expected error for user not found, got nil")
	}

	if response != nil {
		t.Fatal("expected nil response for user not found, got non-nil")
	}
}

func TestDeleteUser_HappyPath(t *testing.T) {
	ctx := context.Background()

	user := &types.User{
		Pid:         "0123456789",
		PhoneNumber: "0949447156",
	}

	service := getUserService()
	_, _ = service.CreateUser(ctx, &types.CreateUserRequest{
		User: user,
	})

	_, err := service.DeleteUser(ctx, &types.DeleteUserRequest{
		Pid: user.Pid,
	})

	if err != nil {
		t.Fatal("DeleteUser failed: " + err.Error())
	}

	response, err := service.GetUser(ctx, &types.GetUserRequest{
		Pid: user.Pid,
	})

	t.Log(err)
	if err == nil {
		t.Fatal("expected error for user not found after deletion, got nil")
	}

	if response != nil {
		t.Fatal("expected nil response for user not found after deletion, got non-nil")
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	ctx := context.Background()
	service := getUserService()

	_, err := service.DeleteUser(ctx, &types.DeleteUserRequest{
		Pid: "9999999999",
	})

	t.Log(err)
	if err == nil {
		t.Fatal("expected error for user not found, got nil")
	}
}

func TestGetUser_HappyPath(t *testing.T) {
	ctx := context.Background()

	user := &types.User{
		Pid:         "0123456789",
		PhoneNumber: "0949447156",
	}

	service := getUserService()
	_, _ = service.CreateUser(ctx, &types.CreateUserRequest{
		User: user,
	})

	response, err := service.GetUser(ctx, &types.GetUserRequest{
		Pid: user.Pid,
	})

	if err != nil {
		t.Fatal("GetUser failed: " + err.Error())
	}
	if response == nil {
		t.Fatal("GetUser returned nil response")
	}
	if response.User.Pid != user.Pid || response.User.PhoneNumber != user.PhoneNumber {
		t.Fatalf("expected user with pid %s and phone number %s, got pid %s and phone number %s", user.Pid, user.PhoneNumber, response.User.Pid, response.User.PhoneNumber)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	ctx := context.Background()
	service := getUserService()

	response, err := service.GetUser(ctx, &types.GetUserRequest{
		Pid: "9999999999",
	})

	t.Log(err)
	if err == nil {
		t.Fatal("expected error for user not found, got nil")
	}
	if response != nil {
		t.Fatal("expected nil response for user not found")
	}
}

func TestListUsers_HappyPath(t *testing.T) {
	ctx := context.Background()

	users := []*types.User{
		{
			Pid:         "0123456789",
			PhoneNumber: "0949447156",
		},
		{
			Pid:         "9876543210",
			PhoneNumber: "0918441887",
		},
	}

	service := getUserService()
	for _, user := range users {
		_, _ = service.CreateUser(ctx, &types.CreateUserRequest{
			User: user,
		})
	}

	response, err := service.ListUsers(ctx, &types.ListUsersRequest{})
	if err != nil {
		t.Fatal("ListUsers failed: " + err.Error())
	}
	if response == nil {
		t.Fatal("ListUsers returned nil response")
	}
	if len(response.Users) != len(users) {
		t.Fatalf("expected %d users, got %d", len(users), len(response.Users))
	}

	for i, user := range users {
		if response.Users[i].Pid != user.Pid || response.Users[i].PhoneNumber != user.PhoneNumber {
			t.Fatalf("expected user with pid %s and phone number %s, got pid %s and phone number %s", user.Pid, user.PhoneNumber, response.Users[i].Pid, response.Users[i].PhoneNumber)
		}
	}
}

func TestListUsers_Empty(t *testing.T) {
	ctx := context.Background()
	service := getUserService()

	response, err := service.ListUsers(ctx, &types.ListUsersRequest{})
	if err != nil {
		t.Fatal("ListUsers failed: " + err.Error())
	}
	if response == nil {
		t.Fatal("ListUsers returned nil response")
	}
	if len(response.Users) != 0 {
		t.Fatalf("expected 0 users, got %d", len(response.Users))
	}
}

// helper
func getUserService() *Service {
	ctx := context.Background()

	db, err := connectDatabase("mongodb://localhost:27017/?directConnection=true", "room_minh_house")
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// removes all documents and indexes
	db.Collection("users").Drop(ctx)

	dbService, err := database.NewMongoService(db)
	if err != nil {
		panic("Failed to create database service: " + err.Error())
	}

	service := NewService(dbService)
	return service
}
