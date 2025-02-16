package grpc

import (
	"context"
	"fmt"

	"github.com/umardev500/gochat/api/proto"
	"github.com/umardev500/gochat/internal/domain"
	"github.com/umardev500/gochat/internal/service"
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
		c.chatService.PushMessage(ctx, &pushChat)
	case *proto.SendMessageRequest_ImageMessage:
		fmt.Println(msg)
	}
	// Implement the logic to send a message
	return nil, nil
}
