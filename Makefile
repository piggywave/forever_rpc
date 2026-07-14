PROTOC := protoc
THRIFT := thrift
GRPC_PY_PLUGIN := grpc_python_plugin

.PHONY: all proto-python proto-go grpc-python grpc-go \
        thrift-python-thriftonly thrift-go-thriftonly thrift-python thrift-go \
        clean

all: grpc-python grpc-go thrift-python thrift-go

# ---- Protobuf: 不带 RPC（只生成 message 结构体） ----

proto-python:
	mkdir -p grpc/python
	$(PROTOC) --proto_path=grpc/proto \
		--python_out=grpc/python \
		grpc/proto/user.proto
	@echo ">>> 只生成了 user_pb2.py（仅 message 定义，无 gRPC service）"

proto-go:
	mkdir -p grpc/go/user
	$(PROTOC) --proto_path=grpc/proto \
		--go_out=grpc/go \
		--go_opt=paths=source_relative \
		grpc/proto/user.proto
	@echo ">>> 只生成了 user.pb.go（仅 message 定义，无 gRPC service）"

# ---- Protobuf + gRPC: 带 RPC 框架 ----

grpc-python: proto-python
	$(PROTOC) --proto_path=grpc/proto \
		--plugin=protoc-gen-grpc_python=$(shell which grpc_python_plugin) \
		--grpc_python_out=grpc/python \
		grpc/proto/user.proto
	@echo ">>> 已生成 user_pb2.py + user_pb2_grpc.py（message + gRPC service）"

grpc-go: proto-go
	$(PROTOC) --proto_path=grpc/proto \
		--go-grpc_out=grpc/go \
		--go-grpc_opt=paths=source_relative \
		grpc/proto/user.proto
	@echo ">>> 已生成 user.pb.go + user_grpc.pb.go（message + gRPC service）"

# ---- Thrift: 不带 RPC（只生成 struct 数据结构，用 no_service 选项） ----

thrift-python-thriftonly:
	mkdir -p thrift/python
	$(THRIFT) --gen py:no_service -out thrift/python thrift/idl/user.thrift
	@echo ">>> 只生成了 ttypes.py 等（仅 struct 定义，无 service）"

thrift-go-thriftonly:
	mkdir -p thrift/go/user
	$(THRIFT) --gen go:no_service -out thrift/go thrift/idl/user.thrift
	@echo ">>> 只生成了 ttypes.go 等（仅 struct 定义，无 service）"

# ---- Thrift: 带 RPC 框架（完整 service） ----

thrift-python:
	mkdir -p thrift/python
	$(THRIFT) --gen py -out thrift/python thrift/idl/user.thrift
	@echo ">>> 已生成 ttypes.py + UserService.py（struct + service）"

thrift-go:
	mkdir -p thrift/go/user
	$(THRIFT) --gen go -out thrift/go thrift/idl/user.thrift
	@echo ">>> 已生成 ttypes.go + user-service.go（struct + service）"

clean:
	rm -rf grpc/python/user_pb2*.py grpc/go/user
	rm -rf thrift/python/user thrift/go/user