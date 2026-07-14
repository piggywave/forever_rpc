package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "./user"

	"google.golang.org/grpc"
)

var users = make(map[int64]*pb.User)
var nextID int64 = 1

type userService struct {
	pb.UnimplementedUserServiceServer
}

func (s *userService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user := users[req.UserId]
	if user != nil {
		return &pb.GetUserResponse{
			User:    user,
			Code:    200,
			Message: "success",
		}, nil
	}
	return &pb.GetUserResponse{
		User:    &pb.User{},
		Code:    404,
		Message: "user not found",
	}, nil
}

func (s *userService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	userID := nextID
	nextID++
	user := &pb.User{
		Id:      userID,
		Name:    req.Name,
		Email:   req.Email,
		Age:     req.Age,
		Address: req.Address,
		Tags:    req.Tags,
	}
	users[userID] = user
	return &pb.CreateUserResponse{
		UserId:  userID,
		Code:    201,
		Message: "user created",
	}, nil
}

func (s *userService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	user := users[req.Id]
	if user == nil {
		return &pb.UpdateUserResponse{
			Success: false,
			Code:    404,
			Message: "user not found",
		}, nil
	}
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Age != nil {
		user.Age = *req.Age
	}
	if req.Address != nil {
		user.Address = *req.Address
	}
	return &pb.UpdateUserResponse{
		Success: true,
		Code:    200,
		Message: "user updated",
	}, nil
}

func (s *userService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if _, exists := users[req.UserId]; exists {
		delete(users, req.UserId)
		return &pb.DeleteUserResponse{
			Success: true,
			Code:    200,
			Message: "user deleted",
		}, nil
	}
	return &pb.DeleteUserResponse{
		Success: false,
		Code:    404,
		Message: "user not found",
	}, nil
}

func (s *userService) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	var allUsers []*pb.User
	for _, u := range users {
		allUsers = append(allUsers, u)
	}

	if req.Keyword != nil {
		keyword := *req.Keyword
		var filtered []*pb.User
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

	return &pb.ListUsersResponse{
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
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &userService{})
	fmt.Println("gRPC server started on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}