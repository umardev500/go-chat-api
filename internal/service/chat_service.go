package service

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/common/model"
	"github.com/umardev500/common/utils"
	"github.com/umardev500/gochat/internal/domain"
	"github.com/umardev500/gochat/internal/repository"
	localUtils "github.com/umardev500/gochat/internal/utils"
)

type ChatService interface {
	FindChatList(ctx context.Context, jid, csid string) *model.Response
	UpdateUnread(ctx context.Context, jid, csid string, value int64) *model.Response
	PushMessage(ctx context.Context, jid, csid string, pushChat *domain.PushChat) *model.Response
}

type chatService struct {
	chatRepo repository.ChatRepository
}

func NewChatService(repo repository.ChatRepository) ChatService {

	return &chatService{
		chatRepo: repo,
	}
}

func (s *chatService) broadcast(socketId string, chat *domain.PushChat) {
	conn := localUtils.WsGetClient(socketId)
	if conn == nil {
		log.Error().Msgf("Websocket client not found: %s", socketId)
		return
	}

	conn.WriteJSON(chat)
}

func (s *chatService) FindChatList(ctx context.Context, jid, csid string) *model.Response {
	if csid == "" {
		return utils.CrateResponse(fiber.StatusBadRequest, "Csid is required", nil)
	}

	chats, err := s.chatRepo.FindChats(ctx, jid, csid)
	if err != nil {
		fmt.Println(err)
		return utils.CrateResponse(fiber.StatusInternalServerError, "Failed to find chat list", nil)
	}

	return utils.CrateResponse(fiber.StatusOK, "Find chat list", chats)
}

func (s *chatService) PushMessage(ctx context.Context, jid, csid string, pushChat *domain.PushChat) *model.Response {
	chatData := domain.CreateChat{
		Jid:      jid,
		Csid:     csid,
		Status:   string(domain.ChatStatusQueued),
		Unread:   1, // Unread is 1 if the first message isn't from customer service
		Messages: []interface{}{pushChat.Data.Message},
	}

	// Check if this is the first message in the chat
	exist, err := s.chatRepo.CreateChat(ctx, jid, csid, chatData)
	if err != nil {
		return utils.CrateResponse(fiber.StatusInternalServerError, "Failed to create chat", nil)
	}

	if exist {
		// Push the message only if the chat exists
		err = s.chatRepo.PushMessage(ctx, jid, csid, pushChat.Data.Message)
		if err != nil {
			return utils.CrateResponse(fiber.StatusInternalServerError, "Failed to push chat", nil)
		}
	} else {
		pushChat.Data.IsInitial = exist
		pushChat.Data.InitialChat = chatData
	}

	// TODO: broadcast to the gprc client
	s.broadcast(csid, pushChat)

	return utils.CrateResponse(fiber.StatusOK, "Push chat", nil)
}

func (s *chatService) UpdateUnread(ctx context.Context, jid, csid string, value int64) *model.Response {
	err := s.chatRepo.UpdateUnread(ctx, jid, csid, value)
	if err != nil {
		return utils.CrateResponse(fiber.StatusInternalServerError, "Failed to update unread", nil)
	}

	return utils.CrateResponse(fiber.StatusOK, "Update unread", nil)
}
