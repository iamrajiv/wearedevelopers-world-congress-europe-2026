package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	hellopb "github.com/iamrajiv/wearedevelopers-world-congress-europe-2026/examples/hello/gen/proto/hello/v1"
)

// server implements the GreeterServiceServer interface
type server struct {
	hellopb.UnimplementedGreeterServiceServer
}

// SayHello implements the SayHello RPC method
func (s *server) SayHello(ctx context.Context, req *hellopb.SayHelloRequest) (*hellopb.SayHelloResponse, error) {
	return &hellopb.SayHelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.Name),
	}, nil
}

func runGRPCServer() error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	hellopb.RegisterGreeterServiceServer(grpcServer, &server{})

	// Register reflection service for grpcurl
	reflection.Register(grpcServer)

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve gRPC: %v", err)
	}
	return nil
}

func runRESTGateway() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create gRPC-Gateway mux
	mux := runtime.NewServeMux()

	// Register the service handler
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := hellopb.RegisterGreeterServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		return fmt.Errorf("failed to register gateway: %v", err)
	}

	// Start HTTP server
	log.Println("gRPC-Gateway (REST) server listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		return fmt.Errorf("failed to serve REST: %v", err)
	}
	return nil
}

func main() {
	// Run gRPC server in a goroutine
	go func() {
		if err := runGRPCServer(); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	// Run REST gateway (blocks)
	if err := runRESTGateway(); err != nil {
		log.Fatalf("REST gateway error: %v", err)
	}
}
