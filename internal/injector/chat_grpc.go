//go:build wireinject

package injector

import (
	"net"

	"github.com/google/wire"
	grpcCmd "github.com/umardev500/gochat/cmd/grpc"
	grpcHandler "github.com/umardev500/gochat/internal/handler/grpc"
)

func ProvideChatGrpHandler() *grpcHandler.ChatGrpHandlerImpl {
	wire.Build(grpcHandler.NewChatGrpHandler)
	return nil
}

func InitializeGRPCServer(lis net.Listener) *grpcCmd.GRPCServer {
	wire.Build(ProvideChatGrpHandler, grpcCmd.NewGrpcServer)
	return nil
}
