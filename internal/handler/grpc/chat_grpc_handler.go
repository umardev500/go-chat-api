package grpc

import (
	"context"
	"fmt"

	"github.com/umardev500/gochat/api/proto"
	"github.com/umardev500/gochat/internal/domain"
)

type ChatGrpHandlerImpl struct {
	proto.UnimplementedWaServiceServer
}

func NewChatGrpHandler() *ChatGrpHandlerImpl {
	return &ChatGrpHandlerImpl{}
}

func (c *ChatGrpHandlerImpl) SendMessage(ctx context.Context, req *proto.SendMessageRequest) (*proto.SendMessageResponse, error) {
	var pushChat = domain.PushChat{}

	switch msg := req.Message.(type) {
	case *proto.SendMessageRequest_TextMessage:
		pushChat.Mt = domain.MessageTypeText
		pushChat.Data.Message = msg
		fmt.Println(pushChat)
	case *proto.SendMessageRequest_ImageMessage:
		fmt.Println(msg)
	}
	// Implement the logic to send a message
	return nil, nil
}
