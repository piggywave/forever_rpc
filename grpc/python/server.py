import grpc
from concurrent import futures
import time
import user_pb2
import user_pb2_grpc

users = {}
next_id = 1

class UserServiceServicer(user_pb2_grpc.UserServiceServicer):
    def GetUser(self, request, context):
        user = users.get(request.user_id)
        if user:
            return user_pb2.GetUserResponse(
                user=user,
                code=200,
                message="success"
            )
        return user_pb2.GetUserResponse(
            user=user_pb2.User(),
            code=404,
            message="user not found"
        )

    def CreateUser(self, request, context):
        global next_id
        user_id = next_id
        next_id += 1
        user = user_pb2.User(
            id=user_id,
            name=request.name,
            email=request.email,
            age=request.age,
            address=request.address,
            tags=request.tags
        )
        users[user_id] = user
        return user_pb2.CreateUserResponse(
            user_id=user_id,
            code=201,
            message="user created"
        )

    def UpdateUser(self, request, context):
        user = users.get(request.id)
        if not user:
            return user_pb2.UpdateUserResponse(
                success=False,
                code=404,
                message="user not found"
            )
        if request.name:
            user.name = request.name
        if request.email:
            user.email = request.email
        if request.age:
            user.age = request.age
        if request.address:
            user.address = request.address
        return user_pb2.UpdateUserResponse(
            success=True,
            code=200,
            message="user updated"
        )

    def DeleteUser(self, request, context):
        if request.user_id in users:
            del users[request.user_id]
            return user_pb2.DeleteUserResponse(
                success=True,
                code=200,
                message="user deleted"
            )
        return user_pb2.DeleteUserResponse(
            success=False,
            code=404,
            message="user not found"
        )

    def ListUsers(self, request, context):
        all_users = list(users.values())
        if request.keyword:
            keyword = request.keyword.lower()
            all_users = [u for u in all_users 
                        if keyword in u.name.lower() or keyword in u.email.lower()]
        start = (request.page - 1) * request.page_size
        end = start + request.page_size
        paginated_users = all_users[start:end]
        return user_pb2.ListUsersResponse(
            users=paginated_users,
            total=len(all_users),
            code=200,
            message="success"
        )

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    user_pb2_grpc.add_UserServiceServicer_to_server(UserServiceServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("gRPC server started on port 50051")
    try:
        while True:
            time.sleep(86400)
    except KeyboardInterrupt:
        server.stop(0)

if __name__ == '__main__':
    serve()