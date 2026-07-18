PROTOC := protoc
THRIFT := thrift
GRPC_PY_PLUGIN := grpc_python_plugin
GOOGLEAPIS := third_party/googleapis

GOPATH := $(shell go env GOPATH)
GRPC_GATEWAY_PLUGIN := $(GOPATH)/bin/protoc-gen-grpc-gateway
OPENAPI_PLUGIN := $(GOPATH)/bin/protoc-gen-openapiv2

PROTO_DIR := grpc/proto
PROTOS := $(wildcard $(PROTO_DIR)/*.proto)

.PHONY: all proto-python proto-go grpc-python grpc-go \
        gateway-go openapi \
        thrift-python-thriftonly thrift-go-thriftonly thrift-python thrift-go \
        kitex-gen kitex-build \
        clean list-protos

all: grpc-python grpc-go gateway-go openapi thrift-python thrift-go

list-protos:
	@echo "Found proto files:"
	@$(foreach proto,$(PROTOS),echo "  $(proto)";)

# ---- Protobuf: 不带 RPC（只生成 message 结构体） ----

proto-python:
	mkdir -p grpc/python
	$(PROTOC) --proto_path=$(PROTO_DIR) \
		--proto_path=$(GOOGLEAPIS) \
		--python_out=grpc/python \
		$(PROTOS)
	@echo ">>> 已生成 *_pb2.py（仅 message 定义，无 gRPC service）"

proto-go:
	$(PROTOC) --proto_path=$(PROTO_DIR) \
		--proto_path=$(GOOGLEAPIS) \
		--go_out=. \
		--go_opt=module=github.com/piggywave/forever_rpc \
		$(PROTOS)
	@echo ">>> 已生成 grpc/go/*/*.pb.go（仅 message 定义，无 gRPC service）"

# ---- Protobuf + gRPC: 带 RPC 框架 ----

grpc-python: proto-python
	$(PROTOC) --proto_path=$(PROTO_DIR) \
		--proto_path=$(GOOGLEAPIS) \
		--plugin=protoc-gen-grpc_python=$(shell which grpc_python_plugin) \
		--grpc_python_out=grpc/python \
		$(PROTOS)
	@echo ">>> 已生成 *_pb2.py + *_pb2_grpc.py（message + gRPC service）"

grpc-go: proto-go
	$(PROTOC) --proto_path=$(PROTO_DIR) \
		--proto_path=$(GOOGLEAPIS) \
		--go-grpc_out=. \
		--go-grpc_opt=module=github.com/piggywave/forever_rpc \
		$(PROTOS)
	@echo ">>> 已生成 grpc/go/*/*.pb.go + *_grpc.pb.go（message + gRPC service）"

# ---- gRPC Gateway: HTTP 网关（REST → gRPC 转码） ----

gateway-go: grpc-go
	$(PROTOC) --proto_path=$(PROTO_DIR) \
		--proto_path=$(GOOGLEAPIS) \
		--plugin=protoc-gen-grpc-gateway=$(GRPC_GATEWAY_PLUGIN) \
		--grpc-gateway_out=. \
		--grpc-gateway_opt=module=github.com/piggywave/forever_rpc \
		--grpc-gateway_opt=logtostderr=true \
		$(PROTOS)
	@echo ">>> 已生成 grpc/go/*/*.pb.gw.go（HTTP → gRPC 网关）"

# ---- OpenAPI / Swagger 文档 ----

openapi:
	mkdir -p docs
	$(PROTOC) --proto_path=$(PROTO_DIR) \
		--proto_path=$(GOOGLEAPIS) \
		--plugin=protoc-gen-openapiv2=$(OPENAPI_PLUGIN) \
		--openapiv2_out=docs \
		--openapiv2_opt=logtostderr=true \
		--openapiv2_opt=json_names_for_fields=false \
		$(PROTOS)
	@echo ">>> 已生成 docs/*.swagger.json（Swagger/OpenAPI 文档）"

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
	rm -f grpc/python/*_pb2*.py
	rm -rf grpc/go/*/
	rm -rf thrift/python/user thrift/go/user

# ---- Kitex: 字节跳动高性能 RPC 框架 ----

kitex-gen:
	cd kitex && KITEX=$(GOPATH)/bin/kitex $$KITEX -module github.com/piggywave/forever_rpc/kitex idl/user.thrift
	@echo ">>> 已生成 kitex_gen/（Kitex 代码生成）"

kitex-build:
	-mv go.work go.work.tmp 2>/dev/null
	cd kitex && go build -o /tmp/kitex_server server.go && go build -o /tmp/kitex_client client.go
	-mv go.work.tmp go.work 2>/dev/null
	@echo ">>> Kitex 服务端和客户端编译成功"