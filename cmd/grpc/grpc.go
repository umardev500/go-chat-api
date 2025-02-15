package grpc

import (
	"net"

	"github.com/rs/zerolog/log"
	"github.com/umardev500/gochat/api/proto"
	grpcHandler "github.com/umardev500/gochat/internal/handler/grpc"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	server   *grpc.Server
	listener net.Listener
}

func NewGrpcServer(lis net.Listener, handler *grpcHandler.ChatGrpHandlerImpl) *GRPCServer {
	server := grpc.NewServer()
	proto.RegisterWaServiceServer(server, handler)

	return &GRPCServer{
		listener: lis,
		server:   server,
	}
}

func (s *GRPCServer) Start() {
	log.Info().Msgf("🚀 gRPC Server is running on %s", s.listener.Addr())
	if err := s.server.Serve(s.listener); err != nil {
		log.Fatal().Err(err).Msg("❌ failed to serve gRPC server")
	}
}
