package main

import (
	"context"
	"log"
	"net"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
	kitex "github.com/piggywave/forever_rpc/kitex/kitex_gen/user/userservice"
	"github.com/piggywave/forever_rpc/kitex/kitex_gen/user"
)

var users = make(map[int64]*user.User)
var nextID = int64(1)

type UserServiceImpl struct{}

func (s *UserServiceImpl) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	u := users[req.UserId]
	if u != nil {
		return &user.GetUserResponse{
			User:    u,
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

func (s *UserServiceImpl) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {
	userID := nextID
	nextID++
	u := &user.User{
		Id:      userID,
		Name:    req.Name,
		Email:   req.Email,
		Age:     req.Age,
		Address: req.Address,
		Tags:    req.Tags,
	}
	users[userID] = u
	return &user.CreateUserResponse{
		UserId:  userID,
		Code:    201,
		Message: "user created",
	}, nil
}

func (s *UserServiceImpl) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.UpdateUserResponse, error) {
	u := users[req.Id]
	if u == nil {
		return &user.UpdateUserResponse{
			Success: false,
			Code:    404,
			Message: "user not found",
		}, nil
	}
	if req.Name != nil {
		u.Name = *req.Name
	}
	if req.Email != nil {
		u.Email = *req.Email
	}
	if req.Age != nil {
		u.Age = *req.Age
	}
	if req.Address != nil {
		u.Address = *req.Address
	}
	return &user.UpdateUserResponse{
		Success: true,
		Code:    200,
		Message: "user updated",
	}, nil
}

func (s *UserServiceImpl) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*user.DeleteUserResponse, error) {
	if _, ok := users[req.UserId]; ok {
		delete(users, req.UserId)
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

func (s *UserServiceImpl) ListUsers(ctx context.Context, req *user.ListUsersRequest) (*user.ListUsersResponse, error) {
	var allUsers []*user.User
	for _, u := range users {
		allUsers = append(allUsers, u)
	}
	if req.Keyword != nil {
		keyword := *req.Keyword
		var filtered []*user.User
		for _, u := range allUsers {
			if contains(u.Name, keyword) || contains(u.Email, keyword) {
				filtered = append(filtered, u)
			}
		}
		allUsers = filtered
	}
	start := (req.Page - 1) * req.PageSize
	end := start + req.PageSize
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

func contains(str, substr string) bool {
	if str == "" || substr == "" {
		return false
	}
	return len(str) >= len(substr) && (str == substr || len(str) > len(substr) && containsHelper(str, substr))
}

func containsHelper(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func main() {
	klog.SetLevel(klog.LevelDebug)

	svr := kitex.NewServer(
		new(UserServiceImpl),
		server.WithServiceAddr(&net.TCPAddr{Port: 8888}),
	)

	err := svr.Run()
	if err != nil {
		log.Fatalf("Kitex server failed to run: %v", err)
	}
}