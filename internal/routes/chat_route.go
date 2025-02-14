package routes

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/umardev500/gochat/internal/handler/http"
	"github.com/umardev500/gochat/internal/middleware"
)

type ChatRouteImpl struct {
	chatHandler http.ChatHandler
}

// 🔹 Implement the interface
func (r *ChatRouteImpl) Api(router fiber.Router) {
	chat := router.Group("/chat")
	chat.Get("/", r.chatHandler.FetchChatList)
	chat.Patch("/update-unread", r.chatHandler.UpdateUnread)
	chat.Post("/push-message", r.chatHandler.PushMessage)
	chat.Get("/ws", middleware.WsMiddlewareCheckAuth, middleware.WsMiddlewareUpgrade, websocket.New(r.chatHandler.WsHandler))
}

func (r *ChatRouteImpl) Web(router fiber.Router) {}

// 🔹 Constructor function
func NewChatRoute(ch http.ChatHandler) *ChatRouteImpl {
	return &ChatRouteImpl{
		chatHandler: ch,
	}
}
