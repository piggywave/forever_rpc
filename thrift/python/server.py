from thrift.server import TServer
from thrift.transport import TSocket
from thrift.transport import TTransport
from thrift.protocol import TBinaryProtocol

from user import UserService
from user.ttypes import *

users = {}
next_id = 1

class UserServiceHandler:
    def GetUser(self, request):
        user = users.get(request.user_id)
        if user:
            return GetUserResponse(user=user, code=200, message="success")
        return GetUserResponse(user=User(), code=404, message="user not found")

    def CreateUser(self, request):
        global next_id
        user_id = next_id
        next_id += 1
        user = User(
            id=user_id,
            name=request.name,
            email=request.email,
            age=request.age,
            address=request.address,
            tags=request.tags
        )
        users[user_id] = user
        return CreateUserResponse(user_id=user_id, code=201, message="user created")

    def UpdateUser(self, request):
        user = users.get(request.id)
        if not user:
            return UpdateUserResponse(success=False, code=404, message="user not found")
        if request.name is not None:
            user.name = request.name
        if request.email is not None:
            user.email = request.email
        if request.age is not None:
            user.age = request.age
        if request.address is not None:
            user.address = request.address
        return UpdateUserResponse(success=True, code=200, message="user updated")

    def DeleteUser(self, request):
        if request.user_id in users:
            del users[request.user_id]
            return DeleteUserResponse(success=True, code=200, message="user deleted")
        return DeleteUserResponse(success=False, code=404, message="user not found")

    def ListUsers(self, request):
        all_users = list(users.values())
        if request.keyword:
            keyword = request.keyword.lower()
            all_users = [u for u in all_users
                        if keyword in u.name.lower() or keyword in u.email.lower()]
        start = (request.page - 1) * request.page_size
        end = start + request.page_size
        paginated_users = all_users[start:end]
        return ListUsersResponse(users=paginated_users, total=len(all_users), code=200, message="success")

def serve():
    handler = UserServiceHandler()
    processor = UserService.Processor(handler)
    transport = TSocket.TServerSocket(host='localhost', port=9090)
    tfactory = TTransport.TBufferedTransportFactory()
    pfactory = TBinaryProtocol.TBinaryProtocolFactory()

    server = TServer.TSimpleServer(processor, transport, tfactory, pfactory)
    print("Thrift server started on port 9090")
    server.serve()

if __name__ == '__main__':
    serve()