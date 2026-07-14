PROTOC := protoc
THRIFT := thrift
GRPC_PY_PLUGIN := grpc_python_plugin

.PHONY: all grpc-python thrift-python grpc-go thrift-go clean

all: grpc-python thrift-python grpc-go thrift-go

grpc-python:
	mkdir -p grpc/python
	$(PROTOC) --proto_path=grpc/proto \
		--python_out=grpc/python \
		--grpc_python_out=grpc/python \
		grpc/proto/user.proto

thrift-python:
	mkdir -p thrift/python
	$(THRIFT) --gen py -out thrift/python thrift/idl/user.thrift

grpc-go:
	mkdir -p grpc/go/user
	$(PROTOC) --proto_path=grpc/proto \
		--go_out=grpc/go \
		--go_opt=paths=source_relative \
		--go-grpc_out=grpc/go \
		--go-grpc_opt=paths=source_relative \
		grpc/proto/user.proto

thrift-go:
	mkdir -p thrift/go/user
	$(THRIFT) --gen go -out thrift/go thrift/idl/user.thrift

clean:
	rm -rf grpc/python/user_pb2*.py grpc/go/user
	rm -rf thrift/python/user thrift/go/user