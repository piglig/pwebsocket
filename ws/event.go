package ws

import (
	"github.com/labstack/gommon/log"
	"sync"
)

type EventType string

const (
	SingleChatEvent EventType = "single_chat"
	GroupChatEvent  EventType = "group_chat"
	BroadChatEvent  EventType = "broadcast_chat"
)

type Event struct {
	Type EventType
}

type EventHandler func(eventType EventType, client *Client)

var (
	handlers   = make(map[EventType]EventHandler)
	handlerMux sync.Mutex
)

func init() {
	RegisterEvent(SingleChatEvent, nil)
	RegisterEvent(GroupChatEvent, nil)
	RegisterEvent(BroadChatEvent, nil)
}

func RegisterEvent(eventType EventType, handler EventHandler) {
	if !isValidEventType(eventType) {
		log.Fatalf("invalid event type %s", eventType)
	}

	handlerMux.Lock()
	defer handlerMux.Unlock()
	_, ok := handlers[eventType]
	if ok {
		log.Fatalf("event type %s has been exist", eventType)
	}

	handlers[eventType] = handler
}

func isValidEventType(eventType EventType) bool {
	switch eventType {
	case SingleChatEvent:
	case GroupChatEvent:
	case BroadChatEvent:
	default:
		return false
	}
	return true
}
