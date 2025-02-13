package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/common/router"
	"github.com/umardev500/gochat/internal/injector"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("failed to load env")
	}
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	routes := injector.InitializeRoutes()
	router.NewFiberRouter(app, routes).Setup()

	ch := make(chan error, 1)
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
			log.Warn().Msgf("PORT not set, defaulting to %s", port)
		}

		log.Info().Msgf("starting server on port %s", port)

		addr := ":" + port
		ch <- app.Listen(addr)
	}()

	select {
	case err := <-ch:
		log.Fatal().Err(err).Msg("failed to start server")
	case <-ctx.Done():
		log.Info().Msg("shutting down server")
		app.Shutdown()
	}
}
