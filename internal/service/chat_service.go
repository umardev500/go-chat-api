package service

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/umardev500/common/constants"
	"github.com/umardev500/common/model"
	"github.com/umardev500/common/utils"
	"github.com/umardev500/gochat/api/proto"
	"github.com/umardev500/gochat/internal/domain"
	"github.com/umardev500/gochat/internal/repository"
	localUtils "github.com/umardev500/gochat/internal/utils"
)

type ChatService interface {
	FindChatList(ctx context.Context, jid, csid string) *model.Response
	UpdateUnread(ctx context.Context, jid, csid string, value int64) *model.Response
	PushMessage(ctx context.Context, pushChat *domain.PushMessage) *model.Response
	StreamingReceiver(req *proto.StreamingRequest)
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
	if socketId == "" {
		log.Error().Msgf("Socket ID is empty: %s", socketId)
		return
	}

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

	chats, err := s.chatRepo.FindChats(ctx, jid, csid, nil)
	if err != nil {
		fmt.Println(err)
		return utils.CrateResponse(fiber.StatusInternalServerError, "Failed to find chat list", nil)
	}

	return utils.CrateResponse(fiber.StatusOK, "Find chat list", chats)
}

func (s *chatService) getCsId(ctx context.Context, jid string) string {
	// Retrieve the CSID from and active chat session
	// The session status must be marked as `active`
	// Once found, extract the CSID of the active session
	chats, err := s.chatRepo.FindChats(ctx, jid, "", utils.ToPtr(domain.ChatStatusActive))
	if err != nil {
		log.Error().Err(err).Msg("Failed to find active chat")
		return ""
	}

	if len(chats) > 0 {
		return chats[0].Csid
	}

	return ""
}

func (s *chatService) PushMessage(ctx context.Context, pushMessage *domain.PushMessage) *model.Response {
	var broadcastMessage = domain.MessageBroadcastResponse{
		Data: pushMessage,
	}

	var jid = pushMessage.Metadata.Jid
	var csid = "xyz" // dummy data only
	var userIdContext = ctx.Value(constants.UserIdContextKey)
	if userIdContext != nil {
		// Retrieve the CSID from the context.
		// This is required only for sending data from the internal app.
		// On the other hand, if the message comes from a WhatsApp client,
		// it will not contain the CSID.

		csid = userIdContext.(string)
	} else {
		// Retrieve the CSID from and active chat session
		// The session status must be marked as `active`
		// Once found, extract the CSID of the active session
		chats, err := s.chatRepo.FindChats(ctx, jid, "", utils.ToPtr(domain.ChatStatusActive))
		if err != nil {
			log.Error().Err(err).Msg("Failed to find active chat")
			return utils.CrateResponse(fiber.StatusInternalServerError, "Failed to find chat list", nil)
		}

		if len(chats) > 0 {
			csid = chats[0].Csid
		}
	}

	initialChatData := domain.CreateChat{
		Jid:    jid,
		Csid:   csid,
		Status: string(domain.ChatStatusActive),
		Messages: []interface{}{
			map[string]interface{}{
				"message":  pushMessage.Message,
				"metadata": &pushMessage.Metadata,
				"unread":   true,
			},
		},
	}

	exist, err := s.chatRepo.CreateChat(ctx, jid, csid, initialChatData)
	if err != nil {
		return utils.CrateResponse(fiber.StatusInternalServerError, "Failed to create chat", nil)
	}

	if exist {
		pushMessage.Unread = true
		err = s.chatRepo.PushMessage(ctx, jid, csid, pushMessage)
		if err != nil {
			return utils.CrateResponse(fiber.StatusInternalServerError, "Failed to push chat", nil)
		}
	} else {
		// The initial data
		broadcastMessage.IsInitial = true
	}

	s.broadcasetWs(csid, &domain.WebsocketBroadcast{
		Type: string(domain.BroadcastMessage),
		Data: &broadcastMessage,
	})

	return utils.CrateResponse(fiber.StatusOK, "Push chat", nil)
}

func (s *chatService) broadcasetWs(socketId string, data interface{}) {
	if socketId == "" {
		log.Error().Msgf("Socket ID is empty: %s", socketId)
		return
	}

	conn := localUtils.WsGetClient(socketId)
	if conn == nil {
		log.Error().Msgf("Websocket client not found: %s", socketId)
		return
	}

	conn.WriteJSON(data)
}

func (s *chatService) StreamingReceiver(req *proto.StreamingRequest) {
	switch msg := req.Message.(type) {
	case *proto.StreamingRequest_StreamingPicture:
		localUtils.GetStreamingClient().PicReqChan <- req
	case *proto.StreamingRequest_StreamTyping:
		s.broadcastTypingStatus(msg)
	case *proto.StreamingRequest_StreamingOnline:
		s.broadcastOnlineStatus(msg)
	default:
		fmt.Println("Unknown streaming message type")
	}
}

func (s *chatService) broadcastOnlineStatus(req *proto.StreamingRequest_StreamingOnline) {
	csid := s.getCsId(context.Background(), req.StreamingOnline.Jid)
	s.broadcasetWs(csid, &domain.WebsocketBroadcast{
		Type: string(domain.BroadcastOnline),
		Data: &proto.StreamingOnlineRequest{
			Jid:    req.StreamingOnline.Jid,
			Online: req.StreamingOnline.Online,
		},
	})
}

func (s *chatService) broadcastTypingStatus(req *proto.StreamingRequest_StreamTyping) {
	csid := s.getCsId(context.Background(), req.StreamTyping.Jid)
	s.broadcasetWs(csid, &domain.WebsocketBroadcast{
		Type: string(domain.BroadcastTyping),
		Data: &proto.StreamTypingRequest{
			Jid:    req.StreamTyping.Jid,
			Typing: req.StreamTyping.Typing,
		},
	})
}

func (s *chatService) UpdateUnread(ctx context.Context, jid, csid string, value int64) *model.Response {
	err := s.chatRepo.UpdateUnread(ctx, jid, csid, value)
	if err != nil {
		return utils.CrateResponse(fiber.StatusInternalServerError, "Failed to update unread", nil)
	}

	return utils.CrateResponse(fiber.StatusOK, "Update unread", nil)
}
