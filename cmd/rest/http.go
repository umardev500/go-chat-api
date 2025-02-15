package rest

import (
	"context"
	"net"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/common/model"
	"github.com/umardev500/common/router"
)

type HttpServer struct {
	server   *fiber.App
	listener net.Listener
	routes   []model.Route
}

func NewHttpServer(lis net.Listener, routes []model.Route) *HttpServer {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	return &HttpServer{
		server:   app,
		listener: lis,
		routes:   routes,
	}
}

func (s *HttpServer) Start(ctx context.Context) {
	router.NewFiberRouter(s.server, s.routes).Setup()

	ch := make(chan error, 1)
	go func() {
		log.Info().Msg("🚀 Rest server is running")
		ch <- s.server.Listener(s.listener)
	}()

	select {
	case <-ctx.Done():
		log.Info().Msg("🛑 Rest server is shutting down")
		s.server.Shutdown()
		return
	case err := <-ch:
		log.Fatal().Err(err).Msg("❌ failed to start rest server")
	}
}
