//go:build wireinject

package injector

import (
	"github.com/google/wire"
	"github.com/umardev500/common/database"
	"github.com/umardev500/gochat/internal/handler/http"
	"github.com/umardev500/gochat/internal/repository"
	"github.com/umardev500/gochat/internal/routes"
	"github.com/umardev500/gochat/internal/service"
)

// 🔹 Shared Wire Set (Reused across different initializers)
var ChatSet = wire.NewSet(
	database.NewMongo,            // Database
	repository.NewChatRepository, // Repository
	service.NewChatService,       // Service
	http.NewChatHandler,          // Handler
)

// ✅ Provide Chat Handler (used in multiple places)
func ProvideChatHandler() http.ChatHandler {
	wire.Build(ChatSet)
	return nil
}

var ChatRouteSet = wire.NewSet(
	ProvideChatHandler,
	routes.NewChatRoute, // Chat route
)

// ✅ Provide Chat Route
func ProvideChatRoute() *routes.ChatRouteImpl {
	wire.Build(ChatRouteSet)

	return nil
}
