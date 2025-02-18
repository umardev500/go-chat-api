package utils

import (
	"sync"

	"github.com/umardev500/gochat/api/proto"
)

type streaming struct {
	Stream     proto.WaService_StreamingServer
	PicReqChan chan *proto.StreamingRequest
	PicResChan chan *proto.StreamingResponse
}

var (
	streamingClient *streaming
	streamingMu     sync.Mutex
)

func DeleteStreamingClient() {
	streamingMu.Lock()
	defer streamingMu.Unlock()

	streamingClient = nil
}

func GetStreamingClient() *streaming {
	streamingMu.Lock()
	defer streamingMu.Unlock()

	return streamingClient
}

func SetStreamingClient(stream proto.WaService_StreamingServer) {
	streamingMu.Lock()
	defer streamingMu.Unlock()

	streamingClient = &streaming{
		Stream:     stream,
		PicReqChan: make(chan *proto.StreamingRequest, 1),
		PicResChan: make(chan *proto.StreamingResponse, 1),
	}
}
