package middleware

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func WsMiddlewareCheckAuth(c *fiber.Ctx) error {
	tokenString := c.Query("token")
	if tokenString == "" {
		log.Error().Msg("token string is empty")
		return fiber.ErrUnauthorized
	}

	// TODO: validate the token
	// for now just set the token string to the context
	c.Locals("id", tokenString)

	return c.Next()
}

func WsMiddlewareUpgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}

	return fiber.ErrUpgradeRequired
}
