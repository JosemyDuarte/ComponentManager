package internal

import (
	context "context"
	"errors"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/JosemyDuarte/ComponentManager/proto/ping"
)

type Server struct {
	port    string
	server  *grpc.Server
	service ping.UnsafePingServiceServer
}

// NewServer creates a new gRPC server with the given options and service
func NewServer(service ping.UnsafePingServiceServer, port string, opt ...grpc.ServerOption) *Server {
	return &Server{
		service: service,
		port:    port,
		server:  grpc.NewServer(opt...),
	}
}

func (s *Server) Name() string {
	return "grpc_server"
}

func (s *Server) Start(initReady chan struct{}) error {
	lis, err := net.Listen("tcp", s.port)
	if err != nil {
		// wrap the error and return it
		return fmt.Errorf("problems starting gRPC server %w", err)
	}

	// close the initReady channel to signal that the server is ready
	close(initReady)

	if servErr := s.server.Serve(lis); !errors.Is(servErr, grpc.ErrServerStopped) {
		return fmt.Errorf("problems starting gRPC server %w", servErr)
	}

	return nil
}

func (s *Server) Shutdown(_ context.Context) error {
	s.server.GracefulStop()

	return nil
}

func (s *Server) StartTimeout() time.Duration {
	return 5 * time.Second
}
