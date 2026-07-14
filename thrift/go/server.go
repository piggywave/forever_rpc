package main

import (
	"fmt"
	"net"

	"github.com/apache/thrift/lib/go/thrift"
	"./user"
)

var users = make(map[int64]*user.User)
var nextID int64 = 1

type UserServiceHandler struct{}

func (h *UserServiceHandler) GetUser(request *user.GetUserRequest) (*user.GetUserResponse, error) {
	userData := users[request.UserId]
	if userData != nil {
		return &user.GetUserResponse{
			User:    userData,
			Code:    200,
			Message: "success",
		}, nil
	}
	return &user.GetUserResponse{
		User:    &user.User{},
		Code:    404,
		Message: "user not found",
	}, nil
}

func (h *UserServiceHandler) CreateUser(request *user.CreateUserRequest) (*user.CreateUserResponse, error) {
	userID := nextID
	nextID++
	userData := &user.User{
		Id:      userID,
		Name:    request.Name,
		Email:   request.Email,
		Age:     request.Age,
		Address: request.Address,
		Tags:    request.Tags,
	}
	users[userID] = userData
	return &user.CreateUserResponse{
		UserId:  userID,
		Code:    201,
		Message: "user created",
	}, nil
}

func (h *UserServiceHandler) UpdateUser(request *user.UpdateUserRequest) (*user.UpdateUserResponse, error) {
	userData := users[request.Id]
	if userData == nil {
		return &user.UpdateUserResponse{
			Success: false,
			Code:    404,
			Message: "user not found",
		}, nil
	}
	if request.Name != nil {
		userData.Name = *request.Name
	}
	if request.Email != nil {
		userData.Email = *request.Email
	}
	if request.Age != nil {
		userData.Age = *request.Age
	}
	if request.Address != nil {
		userData.Address = *request.Address
	}
	return &user.UpdateUserResponse{
		Success: true,
		Code:    200,
		Message: "user updated",
	}, nil
}

func (h *UserServiceHandler) DeleteUser(request *user.DeleteUserRequest) (*user.DeleteUserResponse, error) {
	if _, exists := users[request.UserId]; exists {
		delete(users, request.UserId)
		return &user.DeleteUserResponse{
			Success: true,
			Code:    200,
			Message: "user deleted",
		}, nil
	}
	return &user.DeleteUserResponse{
		Success: false,
		Code:    404,
		Message: "user not found",
	}, nil
}

func (h *UserServiceHandler) ListUsers(request *user.ListUsersRequest) (*user.ListUsersResponse, error) {
	var allUsers []*user.User
	for _, u := range users {
		allUsers = append(allUsers, u)
	}

	if request.Keyword != nil {
		keyword := *request.Keyword
		var filtered []*user.User
		for _, u := range allUsers {
			if contains(u.Name, keyword) || contains(u.Email, keyword) {
				filtered = append(filtered, u)
			}
		}
		allUsers = filtered
	}

	start := (request.Page - 1) * request.PageSize
	end := start + request.PageSize
	if start > int32(len(allUsers)) {
		start = int32(len(allUsers))
	}
	if end > int32(len(allUsers)) {
		end = int32(len(allUsers))
	}

	return &user.ListUsersResponse{
		Users:   allUsers[start:end],
		Total:   int32(len(allUsers)),
		Code:    200,
		Message: "success",
	}, nil
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func main() {
	transport, err := thrift.NewTServerSocket(":9090")
	if err != nil {
		fmt.Printf("Error creating server socket: %v\n", err)
		return
	}

	handler := &UserServiceHandler{}
	processor := user.NewUserServiceProcessor(handler)
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transportFactory := thrift.NewTBufferedTransportFactory(8192)

	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)
	fmt.Println("Thrift server started on port 9090")
	if err := server.Serve(); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}