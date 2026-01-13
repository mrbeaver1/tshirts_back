package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	usecase "github.com/mrbeaver1/tshirts_back/internal/domain/service"
	v1 "github.com/mrbeaver1/tshirts_back/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

// GatewayServer handles REST to gRPC translation
type GatewayServer struct {
	gwmux    *runtime.ServeMux
	httpMux  *http.ServeMux
	endpoint string
}

// NewGatewayServer creates a new GatewayServer
func NewGatewayServer(endpoint string, opts ...runtime.ServeMuxOption) *GatewayServer {
	gwmux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(incomingHeaderMatcher),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true, // Use snake_case for JSON fields
				EmitUnpopulated: false,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	return &GatewayServer{
		gwmux:    gwmux,
		httpMux:  http.NewServeMux(),
		endpoint: endpoint,
	}
}

// RegisterProductHandlerFromEndpoint registers the product service handler from an endpoint
func (gw *GatewayServer) RegisterProductHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return v1.RegisterProductServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}

// RegisterProductHandler registers the product service handler with the gateway
func (gw *GatewayServer) RegisterProductHandler(ctx context.Context, conn *grpc.ClientConn) error {
	return v1.RegisterProductServiceHandler(ctx, gw.gwmux, conn)
}

// SetupRoutes sets up the routes for the gateway server
func (gw *GatewayServer) SetupRoutes() {
	// Register the gRPC-Gateway mux with the main HTTP mux
	gw.httpMux.Handle("/", gw.gwmux)

	// Add health check endpoint
	gw.httpMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	// Add API version prefix route
	gw.httpMux.Handle("/v1/", gw.gwmux)
}

// GetHTTPHandler returns the HTTP handler for the gateway server
func (gw *GatewayServer) GetHTTPHandler() http.Handler {
	return gw.httpMux
}

// incomingHeaderMatcher customizes which headers are passed from HTTP to gRPC
func incomingHeaderMatcher(key string) (string, bool) {
	// Convert header keys to lowercase for comparison
	lowerKey := strings.ToLower(key)

	switch lowerKey {
	case "authorization":
		return key, true
	case "content-type":
		return key, true
	case "x-requested-with":
		return key, true
	default:
		// Pass through headers that start with X-
		if strings.HasPrefix(lowerKey, "x-") {
			return key, true
		}
		return "", false
	}
}

// Run starts the gateway server
func (gw *GatewayServer) Run(ctx context.Context, port string, productUseCase usecase.ProductUseCaseInterface) error {
	// Set up gRPC connection options
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// Connect to gRPC server
	conn, err := grpc.DialContext(ctx, gw.endpoint, opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server: %w", err)
	}
	defer conn.Close()

	// Register handlers
	if err := gw.RegisterProductHandler(ctx, conn); err != nil {
		return fmt.Errorf("failed to register product handler: %w", err)
	}

	// Set up routes
	gw.SetupRoutes()

	// Start HTTP server
	addr := ":" + port
	fmt.Printf("Starting gRPC-Gateway server on %s\n", addr)

	server := &http.Server{
		Addr:    addr,
		Handler: gw.GetHTTPHandler(),
	}

	// Start the server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("gRPC-Gateway server error: %v\n", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Shutdown the server gracefully
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return server.Shutdown(shutdownCtx)
}
