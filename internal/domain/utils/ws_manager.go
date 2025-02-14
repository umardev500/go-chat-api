package utils

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
)

var (
	clients = make(map[string]*websocket.Conn) // Map of connection IDs to connections
	mu      sync.Mutex
)

func WsAddClient(id string, conn *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()

	clients[id] = conn
}

func WsGetClient(id string) *websocket.Conn {
	mu.Lock()
	defer mu.Unlock()

	return clients[id]
}

func WsGetClients() map[string]*websocket.Conn {
	mu.Lock()
	defer mu.Unlock()

	return clients
}

func WsRemoveClient(id string) {
	mu.Lock()
	defer mu.Unlock()

	delete(clients, id)
}
