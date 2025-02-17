package http

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/common/constants"
	commonUtils "github.com/umardev500/common/utils"
	"github.com/umardev500/gochat/api/proto"
	"github.com/umardev500/gochat/internal/domain"
	"github.com/umardev500/gochat/internal/service"
	"github.com/umardev500/gochat/internal/utils"
)

type ChatHandler interface {
	FetchChatList(c *fiber.Ctx) error
	PushMessage(c *fiber.Ctx) error
	UpdateUnread(c *fiber.Ctx) error
	WaPicture(c *fiber.Ctx) error
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
	var request domain.PushMessage
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	validate := validator.New()
	if err := validate.Struct(&request); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.Locals(constants.UserIdContextKey)
	ctx = context.WithValue(ctx, constants.UserIdContextKey, userId)

	resp := h.chatService.PushMessage(ctx, &request)

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

func (h *chatHandler) WaPicture(c *fiber.Ctx) error {
	client := utils.GetStreamingClient()
	if client == nil {
		return c.JSON("No client")
	}

	client.ResChan <- &proto.StreamingResponse{
		Message: &proto.StreamingResponse_StreamingPicture{
			StreamingPicture: &proto.StreamingPictureResponse{
				Jid: c.Params("jid"),
			},
		},
	}

	data := <-client.ReqChan
	resp := commonUtils.CrateResponse(fiber.StatusOK, "Chat picture", data.GetStreamingPicture().Url)

	return c.Status(resp.StatusCode).JSON(resp)
}

func (h *chatHandler) WsHandler(c *websocket.Conn) {
	id := c.Locals("id").(string)
	utils.WsAddClient(id, c)

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
