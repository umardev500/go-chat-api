package domain

type ChatStatus string

const (
	ChatStatusQueued ChatStatus = "queued"
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
	Mt   string       `json:"mt"` // Message type
	Data PushChatData `json:"data"`
}

type PushChatData struct {
	IsInitial   bool        `json:"-"`
	InitialChat CreateChat  `json:"-"`
	Message     interface{} `json:"message"`
}
