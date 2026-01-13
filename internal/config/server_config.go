package config

// ServerConfig holds configuration for gRPC and gRPC-Gateway servers
type ServerConfig struct {
	GRPCPort     string
	GatewayPort  string
	GRPCHost     string
	ReadTimeout  int // in seconds
	WriteTimeout int // in seconds
	IdleTimeout  int // in seconds
	EnableCORS   bool
	CORSOrigins  []string
	Debug        bool
}

// DefaultServerConfig returns a default server configuration
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		GRPCPort:     "50051",
		GatewayPort:  "8080",
		GRPCHost:     "localhost",
		ReadTimeout:  30,
		WriteTimeout: 30,
		IdleTimeout:  120,
		EnableCORS:   true,
		CORSOrigins:  []string{"*"},
		Debug:        true,
	}
}

// GetGRPCAddress returns the full gRPC address
func (c *ServerConfig) GetGRPCAddress() string {
	return c.GRPCHost + ":" + c.GRPCPort
}
