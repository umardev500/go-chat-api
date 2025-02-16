package service

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/common/constants"
	"github.com/umardev500/common/model"
	"github.com/umardev500/common/utils"
	"github.com/umardev500/gochat/internal/domain"
	"github.com/umardev500/gochat/internal/repository"
	localUtils "github.com/umardev500/gochat/internal/utils"
)

type ChatService interface {
	FindChatList(ctx context.Context, jid, csid string) *model.Response
	UpdateUnread(ctx context.Context, jid, csid string, value int64) *model.Response
	PushMessage(ctx context.Context, pushChat *domain.PushChat) *model.Response
}

type chatService struct {
	chatRepo repository.ChatRepository
}

func NewChatService(repo repository.ChatRepository) ChatService {

	return &chatService{
		chatRepo: repo,
	}
}

func (s *chatService) broadcast(socketId string, chat *domain.MessageBroadcastResponse) {
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

func (s *chatService) PushMessage(ctx context.Context, pushChat *domain.PushChat) *model.Response {
	var broadcastMessage = domain.MessageBroadcastResponse{}

	var jid = pushChat.Metadata.Jid
	var csid = ""
	var userIdContext = ctx.Value(constants.UserIdContextKey)
	if userIdContext != nil {
		// Get the csid from the context
		// this needed for only send data from the internal app
		// in the other hand if message coming from whatsapp client
		// that will not containing the csid
		csid = userIdContext.(string)
	} else {
		// Retrieve the CSID from and active chat session
		// The session status must be marked as `active`
		// One found, extract the CSID of the active session

	}

	initialChatData := domain.CreateChat{
		Jid:    jid,
		Csid:   csid,
		Status: string(domain.ChatStatusQueued),
		Unread: 1, // Unread is 1 if the first message isn't from customer service
		Messages: []interface{}{
			map[string]interface{}{
				"message":  pushChat.Message,
				"metadata": &pushChat.Metadata,
			},
		},
	}

	exist, err := s.chatRepo.CreateChat(ctx, jid, csid, initialChatData)
	if err != nil {
		return utils.CrateResponse(fiber.StatusInternalServerError, "Failed to create chat", nil)
	}

	if exist {
		err = s.chatRepo.PushMessage(ctx, jid, csid, pushChat)
		if err != nil {
			return utils.CrateResponse(fiber.StatusInternalServerError, "Failed to push chat", nil)
		}
	} else {
		// The data for broadcasting to the client
		broadcastMessage.InitialChat = &domain.Chat{
			Jid:     jid,
			Csid:    csid,
			Status:  string(domain.ChatStatusQueued),
			Unread:  1,
			Message: pushChat.Message,
		}
		broadcastMessage.Message = nil
	}

	s.broadcast(csid, &broadcastMessage)

	return utils.CrateResponse(fiber.StatusOK, "Push chat", nil)
}

func (s *chatService) UpdateUnread(ctx context.Context, jid, csid string, value int64) *model.Response {
	err := s.chatRepo.UpdateUnread(ctx, jid, csid, value)
	if err != nil {
		return utils.CrateResponse(fiber.StatusInternalServerError, "Failed to update unread", nil)
	}

	return utils.CrateResponse(fiber.StatusOK, "Update unread", nil)
}
