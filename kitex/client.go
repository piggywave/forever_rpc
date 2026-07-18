package main

import (
	"context"
	"fmt"

	"github.com/cloudwego/kitex/client"
	kitex "github.com/piggywave/forever_rpc/kitex/kitex_gen/user/userservice"
	"github.com/piggywave/forever_rpc/kitex/kitex_gen/user"
)

func main() {
	cli, err := kitex.NewClient(
		"user",
		client.WithHostPorts("localhost:8888"),
	)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	createResp, err := cli.CreateUser(ctx, &user.CreateUserRequest{
		Name:    "Alice",
		Email:   "alice@example.com",
		Age:     25,
		Address: "Beijing",
		Tags:    []string{"admin", "vip"},
	})
	if err != nil {
		fmt.Printf("CreateUser failed: %v\n", err)
		return
	}
	fmt.Printf("CreateUser: user_id=%d, code=%d\n", createResp.UserId, createResp.Code)

	getResp, err := cli.GetUser(ctx, &user.GetUserRequest{UserId: createResp.UserId})
	if err != nil {
		fmt.Printf("GetUser failed: %v\n", err)
		return
	}
	fmt.Printf("GetUser: id=%d, name=%s, email=%s\n", getResp.User.Id, getResp.User.Name, getResp.User.Email)

	updateResp, err := cli.UpdateUser(ctx, &user.UpdateUserRequest{
		Id:    createResp.UserId,
		Email: strPtr("alice.new@example.com"),
	})
	if err != nil {
		fmt.Printf("UpdateUser failed: %v\n", err)
		return
	}
	fmt.Printf("UpdateUser: success=%v, code=%d\n", updateResp.Success, updateResp.Code)

	listResp, err := cli.ListUsers(ctx, &user.ListUsersRequest{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		fmt.Printf("ListUsers failed: %v\n", err)
		return
	}
	fmt.Printf("ListUsers: total=%d, users=%d\n", listResp.Total, len(listResp.Users))

	deleteResp, err := cli.DeleteUser(ctx, &user.DeleteUserRequest{UserId: createResp.UserId})
	if err != nil {
		fmt.Printf("DeleteUser failed: %v\n", err)
		return
	}
	fmt.Printf("DeleteUser: success=%v, code=%d\n", deleteResp.Success, deleteResp.Code)

	fmt.Println("All operations completed!")
}

func strPtr(s string) *string {
	return &s
}