//go:build wireinject

package injector

import (
	"net"

	"github.com/google/wire"
	grpcCmd "github.com/umardev500/gochat/cmd/grpc"
	grpcHandler "github.com/umardev500/gochat/internal/handler/grpc"
)

var ProvideChatGrpcHandler = wire.NewSet(grpcHandler.NewChatGrpHandler, ChatSet)

func InitializeGRPCServer(lis net.Listener) *grpcCmd.GRPCServer {
	wire.Build(ProvideChatGrpcHandler, grpcCmd.NewGrpcServer)
	return nil
}
