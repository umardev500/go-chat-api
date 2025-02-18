package grpc

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/umardev500/gochat/api/proto"
	"github.com/umardev500/gochat/internal/domain"
	"github.com/umardev500/gochat/internal/service"
	"github.com/umardev500/gochat/internal/utils"
)

type ChatGrpHandlerImpl struct {
	proto.UnimplementedWaServiceServer
	chatService service.ChatService
}

func NewChatGrpHandler(chatService service.ChatService) *ChatGrpHandlerImpl {
	return &ChatGrpHandlerImpl{
		chatService: chatService,
	}
}

func (c *ChatGrpHandlerImpl) SendMessage(ctx context.Context, req *proto.SendMessageRequest) (*proto.SendMessageResponse, error) {
	var pushChat = domain.PushMessage{}

	switch msg := req.Message.(type) {
	case *proto.SendMessageRequest_TextMessage:
		pushChat.Message = msg.TextMessage
		pushChat.Metadata = req.Metadata
		pushChat.Unread = true
		c.chatService.PushMessage(ctx, &pushChat)
	case *proto.SendMessageRequest_ImageMessage:
		fmt.Println(msg)
	}
	// Implement the logic to send a message
	return nil, nil
}

func (h *ChatGrpHandlerImpl) Streaming(stream proto.WaService_StreamingServer) error {
	utils.SetStreamingClient(stream)
	log.Info().Msg("✅ Streaming client is connected")

	go func() {
		c := utils.GetStreamingClient()
		for msg := range c.ResChan {
			if err := c.Stream.Send(msg); err != nil {
				log.Err(err).Msgf("failed to send streaming data to the client")
			}
		}
	}()

	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Error().Msg("failed to receive streaming message")
			return err
		}

		h.chatService.Streaming(msg)
	}
}
