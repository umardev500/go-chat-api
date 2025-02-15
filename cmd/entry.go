package main

import (
	"context"
	"net"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/soheilhy/cmux"
	"github.com/umardev500/gochat/internal/injector"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("failed to load env")
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Warn().Msgf("PORT not set, defaulting to %s", port)
	}

	log.Info().Msgf("starting server on port %s", port)

	addr := ":" + port
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}
	m := cmux.New(lis)

	// Match gRPC and HTTP connections
	grpcL := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpL := m.Match(cmux.HTTP1Fast())

	go func() {
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
		defer cancel()
		injector.InitializeHttpServer(httpL).Start(ctx)
	}()

	go func() {
		injector.InitializeGRPCServer(grpcL).Start()
	}()

	log.Info().Msgf("🚀 Server is running")
	if err := m.Serve(); err != nil {
		log.Fatal().Err(err).Msg("❌ failed to serve")
	}
}
