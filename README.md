# forever_rpc

Protobuf + gRPC and Thrift learning project. Learn RPC framework development with Python and Go.

## 📚 Overview

This project provides a comprehensive learning path for understanding RPC (Remote Procedure Call) frameworks, covering:

- **Protobuf** + **gRPC** - Google's modern RPC framework based on HTTP/2
- **Thrift** - Apache's cross-language RPC framework

Both implementations feature the same User CRUD service, allowing you to compare and understand the differences between these two popular RPC technologies.

## 📁 Project Structure

```
forever_rpc/
├── grpc/                    # gRPC implementation
│   ├── proto/
│   │   └── user.proto       # Protobuf service definition
│   ├── python/
│   │   ├── server.py        # Python gRPC server
│   │   └── client.py        # Python gRPC client
│   └── go/
│       ├── server.go        # Go gRPC server
│       └── client.go        # Go gRPC client
├── thrift/                  # Thrift implementation
│   ├── idl/
│   │   └── user.thrift      # Thrift service definition
│   ├── python/
│   │   ├── server.py        # Python Thrift server
│   │   └── client.py        # Python Thrift client
│   └── go/
│       ├── server.go        # Go Thrift server
│       └── client.go        # Go Thrift client
├── docs/
│   └── LEARNING_GUIDE.md    # Detailed learning guide
├── Makefile                 # Build automation
└── README.md                # This file
```

## 🚀 Quick Start

### Prerequisites

```bash
# Python dependencies
pip install grpcio grpcio-tools thrift

# Protobuf compiler
brew install protobuf        # macOS
sudo apt install protobuf-compiler  # Ubuntu

# Thrift compiler
brew install thrift          # macOS
sudo apt install thrift-compiler    # Ubuntu

# Go dependencies (for Go examples)
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### Build & Run

```bash
# Build all generated code
make all

# Build individual components
make grpc-python
make thrift-python
make grpc-go
make thrift-go
```

#### gRPC Python

```bash
# Terminal 1 - Start server
cd grpc/python && python server.py

# Terminal 2 - Run client
cd grpc/python && python client.py
```

#### Thrift Python

```bash
# Terminal 1 - Start server
cd thrift/python && python server.py

# Terminal 2 - Run client
cd thrift/python && python client.py
```

#### gRPC Go

```bash
# Terminal 1 - Start server
cd grpc/go && go mod init grpc-go && go mod tidy && go run server.go

# Terminal 2 - Run client
cd grpc/go && go run client.go
```

#### Thrift Go

```bash
# Terminal 1 - Start server
cd thrift/go && go mod init thrift-go && go mod tidy && go run server.go

# Terminal 2 - Run client
cd thrift/go && go run client.go
```

## 📖 Learn More

Check out the [Learning Guide](docs/LEARNING_GUIDE.md) for detailed documentation on:

1. **Protobuf Syntax** - Message definitions, services, and fields
2. **gRPC Concepts** - HTTP/2, streaming, and four RPC modes
3. **Thrift Syntax** - Structs, services, and protocols
4. **Comparison** - gRPC vs Thrift feature comparison
5. **Advanced Topics** - Streaming RPC, authentication, interceptors

## 🔧 API Reference

The User service implements the following RPC methods:

| Method | Description |
|--------|-------------|
| `CreateUser` | Create a new user |
| `GetUser` | Get user by ID |
| `UpdateUser` | Update user information |
| `DeleteUser` | Delete user by ID |
| `ListUsers` | List users with pagination and search |

## 📝 License

MIT License

## 🙋 Contributing

Feel free to submit issues and pull requests!

## 📧 Contact

For questions or feedback, please open an issue on GitHub.