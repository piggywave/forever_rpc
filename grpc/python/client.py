import grpc
import user_pb2
import user_pb2_grpc

def run():
    channel = grpc.insecure_channel('localhost:50051')
    stub = user_pb2_grpc.UserServiceStub(channel)

    print("=== Create User ===")
    create_response = stub.CreateUser(user_pb2.CreateUserRequest(
        name="Alice",
        email="alice@example.com",
        age=25,
        address="123 Main St",
        tags=["developer", "engineer"]
    ))
    print(f"Created user with ID: {create_response.user_id}")

    create_response2 = stub.CreateUser(user_pb2.CreateUserRequest(
        name="Bob",
        email="bob@example.com",
        age=30,
        address="456 Oak Ave",
        tags=["manager"]
    ))
    print(f"Created user with ID: {create_response2.user_id}")

    print("\n=== Get User ===")
    get_response = stub.GetUser(user_pb2.GetUserRequest(user_id=1))
    print(f"User ID: {get_response.user.id}")
    print(f"Name: {get_response.user.name}")
    print(f"Email: {get_response.user.email}")
    print(f"Age: {get_response.user.age}")
    print(f"Address: {get_response.user.address}")
    print(f"Tags: {get_response.user.tags}")

    print("\n=== Update User ===")
    update_response = stub.UpdateUser(user_pb2.UpdateUserRequest(
        id=1,
        name="Alice Updated",
        age=26
    ))
    print(f"Update success: {update_response.success}")
    print(f"Message: {update_response.message}")

    print("\n=== List Users ===")
    list_response = stub.ListUsers(user_pb2.ListUsersRequest(
        page=1,
        page_size=10
    ))
    print(f"Total users: {list_response.total}")
    for u in list_response.users:
        print(f"  - {u.id}: {u.name} ({u.email})")

    print("\n=== List Users with Keyword ===")
    list_response2 = stub.ListUsers(user_pb2.ListUsersRequest(
        page=1,
        page_size=10,
        keyword="alice"
    ))
    print(f"Total users matching 'alice': {list_response2.total}")
    for u in list_response2.users:
        print(f"  - {u.id}: {u.name} ({u.email})")

    print("\n=== Delete User ===")
    delete_response = stub.DeleteUser(user_pb2.DeleteUserRequest(user_id=2))
    print(f"Delete success: {delete_response.success}")
    print(f"Message: {delete_response.message}")

    print("\n=== List Users After Delete ===")
    list_response3 = stub.ListUsers(user_pb2.ListUsersRequest(
        page=1,
        page_size=10
    ))
    print(f"Total users: {list_response3.total}")
    for u in list_response3.users:
        print(f"  - {u.id}: {u.name} ({u.email})")

    channel.close()

if __name__ == '__main__':
    run()