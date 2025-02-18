package domain

import "github.com/umardev500/gochat/api/proto"

type PresenceStatus string

var (
	PresenceStatusTyping PresenceStatus = "typing"
	PresenceStatusOnline PresenceStatus = "online"
)

type ChatStatus string

const (
	ChatStatusActive ChatStatus = "active"
	ChatStatusQueued ChatStatus = "queued"
)

type MessageType string

// MessageType defines different message categories
const (
	MessageTypeText  MessageType = "TEXT"
	MessageTypeImage MessageType = "IMAGE"
	MessageTypeVideo MessageType = "VIDEO"
	MessageTypeAudio MessageType = "AUDIO"
)

type Chat struct {
	Jid     string      `json:"jid"`
	Csid    string      `json:"csid"`
	Status  string      `json:"status"`
	Unread  int         `json:"unread"`
	Message interface{} `json:"message"`
}

type CreateChat struct {
	Jid      string        `json:"jid"`
	Csid     string        `json:"csid"`
	Status   string        `json:"status"`
	Unread   int           `json:"unread"`
	Messages []interface{} `json:"messages"`
}

type PushMessage struct {
	Message  interface{}     `json:"message" validate:"required"`
	Metadata *proto.Metadata `json:"metadata" validate:"required"`
	Unread   bool            `json:"unread"`
}

type MessageBroadcastResponse struct {
	IsInitial bool        `json:"isInitial"`
	Data      interface{} `json:"data,omitempty"`
}

type WebsocketBroadcast struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}
