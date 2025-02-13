package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/umardev500/gochat/internal/handler/http"
)

type ChatRouteImpl struct {
	chatHandler http.ChatHandler
}

// 🔹 Implement the interface
func (r *ChatRouteImpl) Api(router fiber.Router) {
	chat := router.Group("/chat")
	chat.Get("/", r.chatHandler.FetchChatList)
}

func (r *ChatRouteImpl) Web(router fiber.Router) {}

// 🔹 Constructor function
func NewChatRoute(ch http.ChatHandler) *ChatRouteImpl {
	return &ChatRouteImpl{
		chatHandler: ch,
	}
}
