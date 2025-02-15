//go:build wireinject

package injector

import (
	"net"

	"github.com/google/wire"
	"github.com/umardev500/gochat/cmd/rest"
)

func InitializeHttpServer(lis net.Listener) *rest.HttpServer {
	wire.Build(InitializeRoutes, rest.NewHttpServer)
	return nil
}
