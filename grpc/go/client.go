package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/piggywave/forever_rpc/grpc/go/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserServiceClient(conn)

	ctx := context.Background()

	fmt.Println("=== Create User ===")
	createRes, err := c.CreateUser(ctx, &pb.CreateUserRequest{
		Name:    "Alice",
		Email:   "alice@example.com",
		Age:     25,
		Address: "123 Main St",
		Tags:    []string{"developer", "engineer"},
	})
	if err != nil {
		log.Fatalf("CreateUser failed: %v", err)
	}
	fmt.Printf("Created user with ID: %d\n", createRes.UserId)

	createRes2, err := c.CreateUser(ctx, &pb.CreateUserRequest{
		Name:    "Bob",
		Email:   "bob@example.com",
		Age:     30,
		Address: "456 Oak Ave",
		Tags:    []string{"manager"},
	})
	if err != nil {
		log.Fatalf("CreateUser failed: %v", err)
	}
	fmt.Printf("Created user with ID: %d\n", createRes2.UserId)

	fmt.Println("\n=== Get User ===")
	getRes, err := c.GetUser(ctx, &pb.GetUserRequest{UserId: 1})
	if err != nil {
		log.Fatalf("GetUser failed: %v", err)
	}
	fmt.Printf("User ID: %d\n", getRes.User.Id)
	fmt.Printf("Name: %s\n", getRes.User.Name)
	fmt.Printf("Email: %s\n", getRes.User.Email)
	fmt.Printf("Age: %d\n", getRes.User.Age)
	fmt.Printf("Address: %s\n", getRes.User.Address)
	fmt.Printf("Tags: %v\n", getRes.User.Tags)

	fmt.Println("\n=== Update User ===")
	name := "Alice Updated"
	age := int32(26)
	updateRes, err := c.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id:   1,
		Name: &name,
		Age:  &age,
	})
	if err != nil {
		log.Fatalf("UpdateUser failed: %v", err)
	}
	fmt.Printf("Update success: %v\n", updateRes.Success)
	fmt.Printf("Message: %s\n", updateRes.Message)

	fmt.Println("\n=== List Users ===")
	listRes, err := c.ListUsers(ctx, &pb.ListUsersRequest{Page: 1, PageSize: 10})
	if err != nil {
		log.Fatalf("ListUsers failed: %v", err)
	}
	fmt.Printf("Total users: %d\n", listRes.Total)
	for _, u := range listRes.Users {
		fmt.Printf("  - %d: %s (%s)\n", u.Id, u.Name, u.Email)
	}

	fmt.Println("\n=== List Users with Keyword ===")
	keyword := "alice"
	listRes2, err := c.ListUsers(ctx, &pb.ListUsersRequest{Page: 1, PageSize: 10, Keyword: &keyword})
	if err != nil {
		log.Fatalf("ListUsers failed: %v", err)
	}
	fmt.Printf("Total users matching 'alice': %d\n", listRes2.Total)
	for _, u := range listRes2.Users {
		fmt.Printf("  - %d: %s (%s)\n", u.Id, u.Name, u.Email)
	}

	fmt.Println("\n=== Delete User ===")
	deleteRes, err := c.DeleteUser(ctx, &pb.DeleteUserRequest{UserId: 2})
	if err != nil {
		log.Fatalf("DeleteUser failed: %v", err)
	}
	fmt.Printf("Delete success: %v\n", deleteRes.Success)
	fmt.Printf("Message: %s\n", deleteRes.Message)

	fmt.Println("\n=== List Users After Delete ===")
	listRes3, err := c.ListUsers(ctx, &pb.ListUsersRequest{Page: 1, PageSize: 10})
	if err != nil {
		log.Fatalf("ListUsers failed: %v", err)
	}
	fmt.Printf("Total users: %d\n", listRes3.Total)
	for _, u := range listRes3.Users {
		fmt.Printf("  - %d: %s (%s)\n", u.Id, u.Name, u.Email)
	}
}