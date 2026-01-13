package server

import (
	"fmt"
	"net"

	usecase "github.com/mrbeaver1/tshirts_back/internal/domain/service"
	v1 "github.com/mrbeaver1/tshirts_back/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GRPCServer wraps the gRPC server functionality
type GRPCServer struct {
	server        *grpc.Server
	productServer *ProductServer
	port          string
}

// NewGRPCServer creates a new GRPCServer
func NewGRPCServer(productUseCase usecase.ProductUseCaseInterface, port string) *GRPCServer {
	grpcServer := grpc.NewServer()
	productServer := NewProductServer(productUseCase)

	// Register the product service
	v1.RegisterProductServiceServer(grpcServer, productServer)

	// Enable reflection for gRPC services (useful for testing with tools like Postman)
	reflection.Register(grpcServer)

	return &GRPCServer{
		server:        grpcServer,
		productServer: productServer,
		port:          port,
	}
}

// Run starts the gRPC server
func (s *GRPCServer) Run() error {
	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", s.port, err)
	}

	fmt.Printf("Starting gRPC server on port %s\n", s.port)
	return s.server.Serve(lis)
}

// Stop stops the gRPC server gracefully
func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
}

// GetServer returns the underlying gRPC server instance
func (s *GRPCServer) GetServer() *grpc.Server {
	return s.server
}
