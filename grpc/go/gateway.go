package main

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/piggywave/forever_rpc/grpc/go/user"
)

func main() {
	grpcEndpoint := "localhost:50051"
	httpPort := ":8080"

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	log.Printf("gRPC Gateway listening on %s (forwarding to gRPC %s)", httpPort, grpcEndpoint)
	if err := http.ListenAndServe(httpPort, mux); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
