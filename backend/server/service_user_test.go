package server

import (
	"context"
	"testing"

	"github.com/wutdahack1337/room_minh_house/backend/internal/database"
	types "github.com/wutdahack1337/room_minh_house/backend/pkg/pb"
	"golang.org/x/crypto/bcrypt"
)

const (
	P_ID         = "0123456789"
	PHONE_NUMBER = "0949447156"
	ROLE         = "tenant"
)

func newTestIdCard() *types.IdCard {
	return &types.IdCard{
		Name:         "Trần Quang Minh",
		Sex:          "Nam",
		DateOfBirth:  "02-08-2005",
		PlaceOfBirth: "Cần Thơ",
	}
}

func newTestAccount() *types.Account {
	return &types.Account{
		Username: "test_username",
		Password: "test_password",
	}
}

func TestCreateUser_HappyPath(t *testing.T) {
	service := getUserService()

	response, err := service.CreateUser(context.Background(), &types.CreateUserRequest{
		User: &types.User{
			Pid:         P_ID,
			PhoneNumber: PHONE_NUMBER,
			Role:        ROLE,
			IdCard:      newTestIdCard(),
			Account:     newTestAccount(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if response == nil {
		t.Fatal("CreateUser returned nil response")
	}
}

func TestCreateUser_DuplicatePid(t *testing.T) {
	ctx := context.Background()
	service := getUserService()

	user := &types.User{
		Pid:         P_ID,
		PhoneNumber: PHONE_NUMBER,
		Role:        ROLE,
		IdCard:      newTestIdCard(),
		Account:     newTestAccount(),
	}

	_, _ = service.CreateUser(ctx, &types.CreateUserRequest{
		User: user,
	})

	response, err := service.CreateUser(ctx, &types.CreateUserRequest{
		User: user,
	})

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
			PhoneNumber: PHONE_NUMBER,
			Role:        ROLE,
			IdCard:      newTestIdCard(),
			Account:     newTestAccount(),
		},
	})

	if err == nil {
		t.Fatal("expected error for empty pid, got nil")
	}

	if response != nil {
		t.Fatal("expected nil response for empty pid, got non-nil")
	}
}

func TestCreateUser_MissingIdCard(t *testing.T) {
	service := getUserService()
	response, err := service.CreateUser(context.Background(), &types.CreateUserRequest{
		User: &types.User{
			Pid:         P_ID,
			PhoneNumber: PHONE_NUMBER,
			Role:        ROLE,
			Account:     newTestAccount(),
		},
	})

	if err == nil {
		t.Fatal("expected error for missing IdCard, got nil")
	}

	if response != nil {
		t.Fatal("expected nil response for missing IdCard, got non-nil")
	}
}

func TestCreateUser_IdCardMissingName(t *testing.T) {
	service := getUserService()

	idCard := newTestIdCard()
	idCard.Name = ""
	response, err := service.CreateUser(context.Background(), &types.CreateUserRequest{
		User: &types.User{
			Pid:         P_ID,
			PhoneNumber: PHONE_NUMBER,
			Role:        ROLE,
			IdCard:      idCard,
			Account:     newTestAccount(),
		},
	})

	if err == nil {
		t.Fatal("expected error for missing IdCard.Name, got nil")
	}

	if response != nil {
		t.Fatal("expected nil response for missing IdCard.Name, got non-nil")
	}
}

func TestCreateUser_IdCardMissingSex(t *testing.T) {
	service := getUserService()

	idCard := newTestIdCard()
	idCard.Sex = ""
	response, err := service.CreateUser(context.Background(), &types.CreateUserRequest{
		User: &types.User{
			Pid:         P_ID,
			PhoneNumber: PHONE_NUMBER,
			Role:        ROLE,
			IdCard:      idCard,
			Account:     newTestAccount(),
		},
	})

	if err == nil {
		t.Fatal("expected error for missing IdCard.Sex, got nil")
	}

	if response != nil {
		t.Fatal("expected nil response for missing IdCard.Sex, got non-nil")
	}
}

func TestCreateUser_IdCardMissingDateOfBirth(t *testing.T) {
	service := getUserService()

	idCard := newTestIdCard()
	idCard.DateOfBirth = ""
	response, err := service.CreateUser(context.Background(), &types.CreateUserRequest{
		User: &types.User{
			Pid:         P_ID,
			PhoneNumber: PHONE_NUMBER,
			Role:        ROLE,
			IdCard:      idCard,
			Account:     newTestAccount(),
		},
	})

	if err == nil {
		t.Fatal("expected error for missing IdCard.DateOfBirth, got nil")
	}

	if response != nil {
		t.Fatal("expected nil response for missing IdCard.DateOfBirth, got non-nil")
	}
}

func TestCreateUser_IdCardMissingPlaceOfBirth(t *testing.T) {
	service := getUserService()

	idCard := newTestIdCard()
	idCard.PlaceOfBirth = ""
	response, err := service.CreateUser(context.Background(), &types.CreateUserRequest{
		User: &types.User{
			Pid:         P_ID,
			PhoneNumber: PHONE_NUMBER,
			Role:        ROLE,
			IdCard:      idCard,
			Account:     newTestAccount(),
		},
	})

	if err == nil {
		t.Fatal("expected error for missing IdCard.PlaceOfBirth, got nil")
	}

	if response != nil {
		t.Fatal("expected nil response for missing IdCard.PlaceOfBirth, got non-nil")
	}
}

func TestUpdateUser_HappyPath(t *testing.T) {
	ctx := context.Background()
	service := getUserService()

	_, _ = service.CreateUser(ctx, &types.CreateUserRequest{
		User: &types.User{
			Pid:         P_ID,
			PhoneNumber: PHONE_NUMBER,
			Role:        ROLE,
			IdCard:      newTestIdCard(),
			Account:     newTestAccount(),
		},
	})

	user := &types.User{
		Pid:         P_ID,
		PhoneNumber: "0918441887",
		Role:        ROLE,
		IdCard:      newTestIdCard(),
		Account:     newTestAccount(),
	}

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
	service := getUserService()
	response, err := service.UpdateUser(context.Background(), &types.UpdateUserRequest{
		User: &types.User{
			PhoneNumber: PHONE_NUMBER,
			Role:        ROLE,
			IdCard:      newTestIdCard(),
			Account:     newTestAccount(),
		},
	})

	if err == nil {
		t.Fatal("expected error for missing pid, got nil")
	}

	if response != nil {
		t.Fatal("expected nil response for missing pid, got non-nil")
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	service := getUserService()
	response, err := service.UpdateUser(context.Background(), &types.UpdateUserRequest{
		User: &types.User{
			Pid:         "9999999999",
			PhoneNumber: PHONE_NUMBER,
			Role:        ROLE,
			IdCard:      newTestIdCard(),
			Account:     newTestAccount(),
		},
	})

	if err == nil {
		t.Fatal("expected error for user not found, got nil")
	}

	if response != nil {
		t.Fatal("expected nil response for user not found, got non-nil")
	}
}

func TestUpdateUser_IdCardUpdated(t *testing.T) {
	ctx := context.Background()
	service := getUserService()

	user := &types.User{
		Pid:         P_ID,
		PhoneNumber: PHONE_NUMBER,
		Role:        ROLE,
		IdCard:      newTestIdCard(),
		Account:     newTestAccount(),
	}

	_, _ = service.CreateUser(ctx, &types.CreateUserRequest{
		User: user,
	})

	user.IdCard = &types.IdCard{
		Name:         "Trần Tấn Khang",
		Sex:          "Nam",
		DateOfBirth:  "09-02-2012",
		PlaceOfBirth: "Cần Thơ",
	}

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

	if response.User.IdCard.Name != user.IdCard.Name {
		t.Fatalf("expected IdCard.Name %s, got %s", user.IdCard.Name, response.User.IdCard.Name)
	}
}

func TestDeleteUser_HappyPath(t *testing.T) {
	ctx := context.Background()
	service := getUserService()

	user := &types.User{
		Pid:         P_ID,
		PhoneNumber: PHONE_NUMBER,
		Role:        ROLE,
		IdCard:      newTestIdCard(),
		Account:     newTestAccount(),
	}

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

	if err == nil {
		t.Fatal("expected error for user not found after deletion, got nil")
	}

	if response != nil {
		t.Fatal("expected nil response for user not found after deletion, got non-nil")
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	service := getUserService()
	_, err := service.DeleteUser(context.Background(), &types.DeleteUserRequest{
		Pid: "9999999999",
	})

	if err == nil {
		t.Fatal("expected error for user not found, got nil")
	}
}

func TestGetUser_HappyPath(t *testing.T) {
	ctx := context.Background()
	service := getUserService()

	user := &types.User{
		Pid:         P_ID,
		PhoneNumber: PHONE_NUMBER,
		Role:        ROLE,
		IdCard:      newTestIdCard(),
		Account:     newTestAccount(),
	}

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
	service := getUserService()

	response, err := service.GetUser(context.Background(), &types.GetUserRequest{
		Pid: "9999999999",
	})
	if err == nil {
		t.Fatal("expected error for user not found, got nil")
	}

	if response != nil {
		t.Fatal("expected nil response for user not found")
	}
}

func TestListUsers_HappyPath(t *testing.T) {
	ctx := context.Background()
	service := getUserService()

	users := []*types.User{
		{
			Pid:         "0123456789",
			PhoneNumber: "0949447156",
			Role:        ROLE,
			IdCard:      newTestIdCard(),
			Account:     newTestAccount(),
		},
		{
			Pid:         "9876543210",
			PhoneNumber: "0918441887",
			Role:        ROLE,
			IdCard:      newTestIdCard(),
			Account:     newTestAccount(),
		},
	}
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
	service := getUserService()

	response, err := service.ListUsers(context.Background(), &types.ListUsersRequest{})
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

func TestCreateUser_MissingRole(t *testing.T) {
	service := getUserService()
	response, err := service.CreateUser(context.Background(), &types.CreateUserRequest{
		User: &types.User{
			Pid:         P_ID,
			PhoneNumber: PHONE_NUMBER,
			IdCard:      newTestIdCard(),
			Account:     newTestAccount(),
		},
	})

	if err == nil {
		t.Fatal("expected error for missing role, got nil")
	}
	if response != nil {
		t.Fatal("expected nil response for missing role, got non-nil")
	}
}

func TestCreateUser_MissingAccount(t *testing.T) {
	service := getUserService()
	response, err := service.CreateUser(context.Background(), &types.CreateUserRequest{
		User: &types.User{
			Pid:         P_ID,
			PhoneNumber: PHONE_NUMBER,
			Role:        ROLE,
			IdCard:      newTestIdCard(),
		},
	})

	if err == nil {
		t.Fatal("expected error for missing account, got nil")
	}
	if response != nil {
		t.Fatal("expected nil response for missing account, got non-nil")
	}
}

func TestCreateUser_AccountMissingUsername(t *testing.T) {
	service := getUserService()

	account := newTestAccount()
	account.Username = ""
	response, err := service.CreateUser(context.Background(), &types.CreateUserRequest{
		User: &types.User{
			Pid:         P_ID,
			PhoneNumber: PHONE_NUMBER,
			Role:        ROLE,
			IdCard:      newTestIdCard(),
			Account:     account,
		},
	})

	if err == nil {
		t.Fatal("expected error for missing account username, got nil")
	}
	if response != nil {
		t.Fatal("expected nil response for missing account username, got non-nil")
	}
}

func TestCreateUser_AccountMissingPassword(t *testing.T) {
	service := getUserService()

	account := newTestAccount()
	account.Password = ""
	response, err := service.CreateUser(context.Background(), &types.CreateUserRequest{
		User: &types.User{
			Pid:         P_ID,
			PhoneNumber: PHONE_NUMBER,
			Role:        ROLE,
			IdCard:      newTestIdCard(),
			Account:     account,
		},
	})

	if err == nil {
		t.Fatal("expected error for missing account password, got nil")
	}
	if response != nil {
		t.Fatal("expected nil response for missing account password, got non-nil")
	}
}

func TestUpdateUser_MissingAccount(t *testing.T) {
	ctx := context.Background()
	service := getUserService()

	_, _ = service.CreateUser(ctx, &types.CreateUserRequest{
		User: &types.User{
			Pid:         P_ID,
			PhoneNumber: PHONE_NUMBER,
			IdCard:      newTestIdCard(),
			Account:     newTestAccount(),
		},
	})

	response, err := service.UpdateUser(context.Background(), &types.UpdateUserRequest{
		User: &types.User{
			Pid:         P_ID,
			PhoneNumber: PHONE_NUMBER,
			IdCard:      newTestIdCard(),
		},
	})
	if err == nil {
		t.Fatal("expected error for missing account, got nil")
	}

	if response != nil {
		t.Fatal("expected nil response for missing account, got non-nil")
	}
}

func TestUpdateUser_AccountMissingUsername(t *testing.T) {
	ctx := context.Background()
	service := getUserService()

	_, err := service.CreateUser(ctx, &types.CreateUserRequest{
		User: &types.User{
			Pid:         P_ID,
			PhoneNumber: PHONE_NUMBER,
			Role:        ROLE,
			IdCard:      newTestIdCard(),
			Account:     newTestAccount(),
		},
	})
	if err != nil {
		t.Fatal("CreateUser failed: " + err.Error())
	}

	account := newTestAccount()
	account.Username = ""
	response, err := service.UpdateUser(context.Background(), &types.UpdateUserRequest{
		User: &types.User{
			Pid:         P_ID,
			PhoneNumber: PHONE_NUMBER,
			Role:        ROLE,
			IdCard:      newTestIdCard(),
			Account:     account,
		},
	})

	if err == nil {
		t.Fatal("expected error for missing account username, got nil")
	}
	if response != nil {
		t.Fatal("expected nil response for missing account username, got non-nil")
	}
}

func TestUpdateUser_AccountMissingPassword(t *testing.T) {
	ctx := context.Background()
	service := getUserService()

	_, _ = service.CreateUser(ctx, &types.CreateUserRequest{
		User: &types.User{
			Pid:         P_ID,
			PhoneNumber: PHONE_NUMBER,
			Role:        ROLE,
			IdCard:      newTestIdCard(),
			Account:     newTestAccount(),
		},
	})

	account := newTestAccount()
	account.Password = ""
	response, err := service.UpdateUser(context.Background(), &types.UpdateUserRequest{
		User: &types.User{
			Pid:         P_ID,
			PhoneNumber: PHONE_NUMBER,
			Role:        ROLE,
			IdCard:      newTestIdCard(),
			Account:     account,
		},
	})

	if err == nil {
		t.Fatal("expected error for missing account password, got nil")
	}
	if response != nil {
		t.Fatal("expected nil response for missing account password, got non-nil")
	}
}

func TestUpdateUser_MissingRole(t *testing.T) {
	ctx := context.Background()
	service := getUserService()

	_, _ = service.CreateUser(ctx, &types.CreateUserRequest{
		User: &types.User{
			Pid:         P_ID,
			PhoneNumber: PHONE_NUMBER,
			Role:        ROLE,
			IdCard:      newTestIdCard(),
			Account:     newTestAccount(),
		},
	})

	response, err := service.UpdateUser(context.Background(), &types.UpdateUserRequest{
		User: &types.User{
			Pid:         P_ID,
			PhoneNumber: PHONE_NUMBER,
			IdCard:      newTestIdCard(),
			Account:     newTestAccount(),
		},
	})
	if err == nil {
		t.Fatal("expected error for missing role, got nil")
	}

	if response != nil {
		t.Fatal("expected nil response for missing role, got non-nil")
	}
}

func TestGetUser_PasswordIsHashed(t *testing.T) {
	ctx := context.Background()
	service := getUserService()

	plainPassword := "test_password_secure_123"
	user := &types.User{
		Pid:         P_ID,
		PhoneNumber: PHONE_NUMBER,
		Role:        ROLE,
		IdCard:      newTestIdCard(),
		Account: &types.Account{
			Username: "test_user",
			Password: plainPassword,
		},
	}

	_, err := service.CreateUser(ctx, &types.CreateUserRequest{
		User: user,
	})
	if err != nil {
		t.Fatal("CreateUser failed: " + err.Error())
	}

	response, err := service.GetUser(ctx, &types.GetUserRequest{
		Pid: user.Pid,
	})
	if err != nil {
		t.Fatal("GetUser failed: " + err.Error())
	}

	if response.User.Account.Password == plainPassword {
		t.Fatal("expected password to be hashed, but got plaintext")
	}

	// Verify the hash is valid bcrypt by comparing with plaintext
	if err := bcrypt.CompareHashAndPassword([]byte(response.User.Account.Password), []byte(plainPassword)); err != nil {
		t.Fatal("password hash does not match plaintext password: " + err.Error())
	}
}

func TestListUsers_PasswordsAreHashed(t *testing.T) {
	ctx := context.Background()
	service := getUserService()

	plainPassword := "test_password_list_456"
	users := []*types.User{
		{
			Pid:         "0111111111",
			PhoneNumber: "0911111111",
			Role:        ROLE,
			IdCard:      newTestIdCard(),
			Account: &types.Account{
				Username: "user1_list",
				Password: plainPassword,
			},
		},
		{
			Pid:         "0222222222",
			PhoneNumber: "0922222222",
			Role:        ROLE,
			IdCard:      newTestIdCard(),
			Account: &types.Account{
				Username: "user2_list",
				Password: plainPassword,
			},
		},
	}

	for _, user := range users {
		_, err := service.CreateUser(ctx, &types.CreateUserRequest{
			User: user,
		})
		if err != nil {
			t.Fatal("CreateUser failed for user " + user.Pid + ": " + err.Error())
		}
	}

	response, err := service.ListUsers(ctx, &types.ListUsersRequest{})
	if err != nil {
		t.Fatal("ListUsers failed: " + err.Error())
	}

	if len(response.Users) == 0 {
		t.Fatal("expected users in response, got empty list")
	}

	for _, user := range response.Users {
		if user.Account.Password == plainPassword {
			t.Fatal("expected password to be hashed for user " + user.Pid + ", but got plaintext")
		}

		// Verify the hash is valid bcrypt
		if err := bcrypt.CompareHashAndPassword([]byte(user.Account.Password), []byte(plainPassword)); err != nil {
			t.Fatal("password hash does not match plaintext password for user " + user.Pid + ": " + err.Error())
		}
	}
}

// helper
func getUserService() *Service {
	db, err := connectDatabase("mongodb://localhost:27017/?directConnection=true", "room_minh_house")
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// removes all documents and indexes
	db.Collection("users").Drop(context.Background())

	dbService, err := database.NewMongoService(db)
	if err != nil {
		panic("Failed to create database service: " + err.Error())
	}

	return NewService(dbService)
}
