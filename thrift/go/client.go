package main

import (
	"fmt"

	"github.com/apache/thrift/lib/go/thrift"
	"./user"
)

func main() {
	transport, err := thrift.NewTSocket("localhost:9090")
	if err != nil {
		fmt.Printf("Error creating client socket: %v\n", err)
		return
	}

	transportFactory := thrift.NewTBufferedTransportFactory(8192)
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	clientTransport, err := transportFactory.GetTransport(transport)
	if err != nil {
		fmt.Printf("Error getting transport: %v\n", err)
		return
	}

	protocol := protocolFactory.GetProtocol(clientTransport)
	client := user.NewUserServiceClient(protocol)

	if err := clientTransport.Open(); err != nil {
		fmt.Printf("Error opening transport: %v\n", err)
		return
	}
	defer clientTransport.Close()

	fmt.Println("=== Create User ===")
	createRes, err := client.CreateUser(&user.CreateUserRequest{
		Name:    "Alice",
		Email:   "alice@example.com",
		Age:     25,
		Address: "123 Main St",
		Tags:    []string{"developer", "engineer"},
	})
	if err != nil {
		fmt.Printf("CreateUser failed: %v\n", err)
		return
	}
	fmt.Printf("Created user with ID: %d\n", createRes.UserId)

	createRes2, err := client.CreateUser(&user.CreateUserRequest{
		Name:    "Bob",
		Email:   "bob@example.com",
		Age:     30,
		Address: "456 Oak Ave",
		Tags:    []string{"manager"},
	})
	if err != nil {
		fmt.Printf("CreateUser failed: %v\n", err)
		return
	}
	fmt.Printf("Created user with ID: %d\n", createRes2.UserId)

	fmt.Println("\n=== Get User ===")
	getRes, err := client.GetUser(&user.GetUserRequest{UserId: 1})
	if err != nil {
		fmt.Printf("GetUser failed: %v\n", err)
		return
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
	updateRes, err := client.UpdateUser(&user.UpdateUserRequest{
		Id:   1,
		Name: &name,
		Age:  &age,
	})
	if err != nil {
		fmt.Printf("UpdateUser failed: %v\n", err)
		return
	}
	fmt.Printf("Update success: %v\n", updateRes.Success)
	fmt.Printf("Message: %s\n", updateRes.Message)

	fmt.Println("\n=== List Users ===")
	listRes, err := client.ListUsers(&user.ListUsersRequest{Page: 1, PageSize: 10})
	if err != nil {
		fmt.Printf("ListUsers failed: %v\n", err)
		return
	}
	fmt.Printf("Total users: %d\n", listRes.Total)
	for _, u := range listRes.Users {
		fmt.Printf("  - %d: %s (%s)\n", u.Id, u.Name, u.Email)
	}

	fmt.Println("\n=== List Users with Keyword ===")
	keyword := "alice"
	listRes2, err := client.ListUsers(&user.ListUsersRequest{Page: 1, PageSize: 10, Keyword: &keyword})
	if err != nil {
		fmt.Printf("ListUsers failed: %v\n", err)
		return
	}
	fmt.Printf("Total users matching 'alice': %d\n", listRes2.Total)
	for _, u := range listRes2.Users {
		fmt.Printf("  - %d: %s (%s)\n", u.Id, u.Name, u.Email)
	}

	fmt.Println("\n=== Delete User ===")
	deleteRes, err := client.DeleteUser(&user.DeleteUserRequest{UserId: 2})
	if err != nil {
		fmt.Printf("DeleteUser failed: %v\n", err)
		return
	}
	fmt.Printf("Delete success: %v\n", deleteRes.Success)
	fmt.Printf("Message: %s\n", deleteRes.Message)

	fmt.Println("\n=== List Users After Delete ===")
	listRes3, err := client.ListUsers(&user.ListUsersRequest{Page: 1, PageSize: 10})
	if err != nil {
		fmt.Printf("ListUsers failed: %v\n", err)
		return
	}
	fmt.Printf("Total users: %d\n", listRes3.Total)
	for _, u := range listRes3.Users {
		fmt.Printf("  - %d: %s (%s)\n", u.Id, u.Name, u.Email)
	}
}