package http

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/umardev500/gochat/internal/service"
)

type ChatHandler interface {
	FetchChatList(c *fiber.Ctx) error
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
