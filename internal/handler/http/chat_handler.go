package http

import (
	"context"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/gochat/internal/domain"
	"github.com/umardev500/gochat/internal/domain/utils"
	"github.com/umardev500/gochat/internal/service"
)

type ChatHandler interface {
	FetchChatList(c *fiber.Ctx) error
	PushMessage(c *fiber.Ctx) error
	UpdateUnread(c *fiber.Ctx) error
	WsHandler(c *websocket.Conn)
}

type chatHandler struct {
	chatService service.ChatService
}

func NewChatHandler(cs service.ChatService) ChatHandler {
	return &chatHandler{
		chatService: cs,
	}
}

func (h *chatHandler) FetchChatList(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jid := c.Query("jid")
	csid := c.Query("csid")

	resp := h.chatService.FindChatList(ctx, jid, csid)

	return c.Status(resp.StatusCode).JSON(resp)
}

func (h *chatHandler) PushMessage(c *fiber.Ctx) error {
	var request domain.PushChat
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jid := c.Query("jid")
	csid := c.Query("csid")

	if jid == "" || csid == "" {
		return fiber.ErrBadRequest
	}

	resp := h.chatService.PushMessage(ctx, jid, csid, &request)

	return c.Status(resp.StatusCode).JSON(resp)
}

func (h *chatHandler) UpdateUnread(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jid := c.Query("jid")
	csid := c.Query("csid")
	value := c.QueryInt("value", 1)

	resp := h.chatService.UpdateUnread(ctx, jid, csid, int64(value))

	return c.Status(resp.StatusCode).JSON(resp)
}

func (h *chatHandler) WsHandler(c *websocket.Conn) {
	id := c.Locals("id").(string)

	var (
		err error
	)

	log.Info().Msgf("Client connected: %s", id)

	for {
		_, _, err = c.ReadMessage()
		if err != nil {
			log.Error().Msgf("Failed to read message: %v", err)
			utils.WsRemoveClient(id)
			break
		}
	}
}
