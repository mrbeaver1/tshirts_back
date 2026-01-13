package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	container "github.com/mrbeaver1/tshirts_back/internal/application/container"
	"github.com/mrbeaver1/tshirts_back/internal/config"
	"github.com/mrbeaver1/tshirts_back/internal/infrastructure/database"
	"github.com/mrbeaver1/tshirts_back/internal/infrastructure/server"
)

func main() {
	// Load server configuration
	serverCfg := config.DefaultServerConfig()

	// Run migrations
	err := database.RunMigrations()
	if err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	// Create dependency container
	cont := container.NewContainer()

	// Get product use case
	productUseCase := cont.GetProductUseCase()

	// Create gRPC server
	grpcServer := server.NewGRPCServer(productUseCase, serverCfg.GRPCPort)

	// Create gRPC-Gateway server
	gatewayServer := server.NewGatewayServer(serverCfg.GetGRPCAddress())

	// Set up channel to listen for interrupt signal to terminate server
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup

	// Channel to signal when servers are ready
	grpcReady := make(chan bool)
	gatewayReady := make(chan bool)

	// Start gRPC server in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("Starting gRPC server on %s", serverCfg.GetGRPCAddress())

		lis, err := net.Listen("tcp", ":"+serverCfg.GRPCPort)
		if err != nil {
			log.Fatalf("Failed to listen on port %s: %v", serverCfg.GRPCPort, err)
		}

		// Signal that gRPC server is ready to accept connections
		close(grpcReady)

		if err := grpcServer.GetServer().Serve(lis); err != nil {
			log.Printf("gRPC server failed to serve: %v", err)
		}
	}()

	// Wait for gRPC server to be ready before starting gateway
	<-grpcReady

	// Start gRPC-Gateway server in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Create context with cancellation for gateway server
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		log.Printf("Starting gRPC-Gateway server on port %s", serverCfg.GatewayPort)

		if err := gatewayServer.Run(ctx, serverCfg.GatewayPort, productUseCase); err != nil {
			if err != http.ErrServerClosed {
				log.Printf("gRPC-Gateway server failed to run: %v", err)
			}
		}

		// Signal that gateway server is ready
		close(gatewayReady)
	}()

	log.Printf("Servers started - gRPC on %s, gRPC-Gateway on port %s",
		serverCfg.GetGRPCAddress(), serverCfg.GatewayPort)
	log.Println("Press Ctrl+C to stop the servers")

	// Wait for interrupt signal
	<-interrupt
	log.Println("Shutting down servers...")

	// Stop gRPC server gracefully
	grpcServer.Stop()

	// Wait for all goroutines to finish
	wg.Wait()
	log.Println("Servers stopped")
}
