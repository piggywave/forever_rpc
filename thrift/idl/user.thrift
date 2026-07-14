namespace py user
namespace go user

struct User {
    1: i64 id
    2: string name
    3: string email
    4: i32 age
    5: string address
    6: list<string> tags
}

struct GetUserRequest {
    1: i64 user_id
}

struct GetUserResponse {
    1: User user
    2: i32 code
    3: string message
}

struct CreateUserRequest {
    1: string name
    2: string email
    3: i32 age
    4: string address
    5: list<string> tags
}

struct CreateUserResponse {
    1: i64 user_id
    2: i32 code
    3: string message
}

struct UpdateUserRequest {
    1: i64 id
    2: optional string name
    3: optional string email
    4: optional i32 age
    5: optional string address
}

struct UpdateUserResponse {
    1: bool success
    2: i32 code
    3: string message
}

struct DeleteUserRequest {
    1: i64 user_id
}

struct DeleteUserResponse {
    1: bool success
    2: i32 code
    3: string message
}

struct ListUsersRequest {
    1: i32 page
    2: i32 page_size
    3: optional string keyword
}

struct ListUsersResponse {
    1: list<User> users
    2: i32 total
    3: i32 code
    4: string message
}

service UserService {
    GetUserResponse GetUser(1: GetUserRequest request)
    CreateUserResponse CreateUser(1: CreateUserRequest request)
    UpdateUserResponse UpdateUser(1: UpdateUserRequest request)
    DeleteUserResponse DeleteUser(1: DeleteUserRequest request)
    ListUsersResponse ListUsers(1: ListUsersRequest request)
}