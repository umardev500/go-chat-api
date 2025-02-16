package domain

import "github.com/umardev500/gochat/api/proto"

type ChatStatus string

const (
	ChatStatusQueued ChatStatus = "queued"
)

// MessageType defines different message categories
const (
	MessageTypeText  string = "TEXT"
	MessageTypeImage string = "IMAGE"
	MessageTypeVideo string = "VIDEO"
	MessageTypeAudio string = "AUDIO"
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

type PushChat struct {
	Message  interface{}     `json:"message" validate:"required"`
	Metadata *proto.Metadata `json:"metadata" validate:"required"`
}

type MessageBroadcastResponse struct {
	InitialChat *Chat       `json:"initialChat,omitempty"`
	Message     interface{} `json:"message,omitempty"`
}
