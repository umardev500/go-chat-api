package grpc

import (
	"context"

	"github.com/umardev500/gochat/api/proto"
)

type ChatGrpHandlerImpl struct {
	proto.UnimplementedWaServiceServer
}

func NewChatGrpHandler() *ChatGrpHandlerImpl {
	return &ChatGrpHandlerImpl{}
}

func (c *ChatGrpHandlerImpl) SendMessage(ctx context.Context, req *proto.SendMessageRequest) (*proto.SendMessageResponse, error) {
	// Implement the logic to send a message
	return nil, nil
}
