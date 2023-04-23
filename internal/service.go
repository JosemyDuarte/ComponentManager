package internal

import (
	context "context"

	"github.com/JosemyDuarte/ComponentManager/proto/ping"
)

type Service struct {
	ping.UnimplementedPingServiceServer
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Ping(_ context.Context, _ *ping.PingRequest) (*ping.PingResponse, error) {
	return &ping.PingResponse{Message: "pong"}, nil
}
